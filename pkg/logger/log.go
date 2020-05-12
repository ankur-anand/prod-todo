package logger

import (
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

// Duration constructs a string field with the given key and time value
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}

// Info logs info level message
func (log *Logger) Info(msg string, fields ...Field) {
	log.sampled.Info(msg, fields...)
}

// Error logs Error level message
func (log *Logger) Error(msg string, field ...Field) {
	log.nonSampled.Error(msg, field...)
}

// Panic logs panic level message
func (log *Logger) Panic(msg string, field ...Field) {
	log.nonSampled.Panic(msg, field...)
}

// Fatal logs fatal level message
func (log *Logger) Fatal(msg string, field ...Field) {
	log.nonSampled.Fatal(msg, field...)
}

// Skip constructs a no-op field, which is often useful when handling invalid
// inputs in other Field constructors.
func Skip() Field {
	return Field{Type: zapcore.SkipType}
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: zapcore.ErrorType, Interface: err}
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) Field {
	return NamedError("error", err)
}

// Sync should be called before application exit.
func (log *Logger) Sync() {
	if log.isDev {
		_ = log.sampled.Sync()
		return
	}
	_ = log.sampled.Sync()
	_ = log.nonSampled.Sync()
}

// NewProduction returns a new Production ready Logger
func NewProduction() (*Logger, error) {
	sampled := zap.NewProductionConfig()
	nonSampled := zap.NewProductionConfig()
	// disable error sampling
	nonSampled.Sampling = nil
	logSamp, err := sampled.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	logNonSamp, err := nonSampled.Build(zap.AddCallerSkip(1))
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
	logDev, err := dev.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{
		sampled:    logDev,
		nonSampled: logDev,
		isDev:      true,
	}, nil
}
