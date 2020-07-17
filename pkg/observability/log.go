package observability

import (
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewProduction creates a new zap logger and
// put this into the global logger. It also return a http
// handler that can be used to change the log level dynamically.
func NewProduction(appname, hostname string) http.Handler {
	sampled := zap.NewProductionEncoderConfig()
	sampled.EncodeTime = zapcore.ISO8601TimeEncoder
	atom := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(sampled), os.Stdout, atom)
	core = zapcore.NewSamplerWithOptions(core, time.Second, 100, 50)
	log := zap.New(core, zap.AddCaller())
	log = log.With(zap.String("application", appname), zap.String("hostname", hostname))
	zap.ReplaceGlobals(log)
	return atom
}

// NewDevelopment returns a new Dev Logger
func NewDevelopment(appname, hostname string) {
	dev, _ := zap.NewDevelopment()
	log := dev.With(zap.String("application", appname), zap.String("hostname", hostname))
	zap.ReplaceGlobals(log)
}

// A WriteSyncer is an io.Writer that can also flush any buffered data.
// Used only for testing config
type WriteSyncer interface {
	io.Writer
	Sync() error
}

// no op
type noOpWriteSyncer struct {
}

func (d noOpWriteSyncer) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (d noOpWriteSyncer) Sync() error {
	return nil
}

// NewTesting returns a new Test Logger with disabled stack trace
// and can take WriteSyncer for output
func NewTesting(ws WriteSyncer) *zap.Logger {
	if ws == nil {
		ws = noOpWriteSyncer{}
	}
	enc := zap.NewDevelopmentEncoderConfig()
	enc.TimeKey = "ts"
	enc.MessageKey = "message"
	enc.LevelKey = "level"
	en := zapcore.NewJSONEncoder(enc)
	core := zapcore.NewCore(en, ws, zap.DebugLevel)
	l := zap.New(core)
	return l
}
