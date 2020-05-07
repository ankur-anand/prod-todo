package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field for structured logging
type Field = zap.Field

// Logger supports various sampling for different log level
// over the same io.writer
// Sync should be called before application exit.
type Logger struct {
	sampled    *zap.Logger
	nonSampled *zap.Logger
	isDev      bool
}

// String constructs a string field with the given key and value
func String(key string, value string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: value}
}

// Int constructs a string field with the given key and value
func Int(key string, val int) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: int64(val)}
}

func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}

func (log *Logger) Info(msg string, fields ...Field) {
	log.sampled.Info(msg, fields...)
}

func (log *Logger) Error(msg string, field ...Field) {
	log.nonSampled.Error(msg, field...)
}

func (log *Logger) Panic(msg string, field ...Field) {
	log.nonSampled.Panic(msg, field...)
}

func (log *Logger) Fatal(msg string, field ...Field) {
	log.nonSampled.Fatal(msg, field...)
}

// HTTPRequest is an helper method to log an HTTPRequest directly
func (log *Logger) HTTPRequest(message string, statusCode int, duration time.Duration, r *http.Request, err error) {
	if err != nil {
		log.nonSampled.Error(message, String("method", r.Method), String("url", r.URL.String()), Int("status", statusCode), Duration("duration", duration), zap.Error(err))
		return
	}
	log.sampled.Error(message, String("method", r.Method), String("url", r.URL.String()), Int("status", statusCode), Duration("duration", duration))
}

func (log *Logger) Sync() {
	if log.isDev {
		log.sampled.Sync()
		return
	}
	log.sampled.Sync()
	log.nonSampled.Sync()
}

// NewProduction returns a new Production ready Logger
func NewProduction() (*Logger, error) {
	sampled := zap.NewProductionConfig()
	nonSampled := zap.NewProductionConfig()
	// disable error sampling
	nonSampled.Sampling = nil
	logSamp, err := sampled.Build()
	if err != nil {
		return nil, err
	}
	logNonSamp, err := nonSampled.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		sampled:    logSamp,
		nonSampled: logNonSamp,
	}, nil
}

// NewDevelopment returns a new Dev Logger
func NewDevelopment() (*Logger, error) {
	dev := zap.NewDevelopmentConfig()
	logDev, err := dev.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		sampled:    logDev,
		nonSampled: logDev,
		isDev:      true,
	}, nil
}
