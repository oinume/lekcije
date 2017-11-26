package logger

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: consider to put them into context
var (
	Access *zap.Logger
	App    *zap.Logger
)

func init() {
	err := zap.RegisterEncoder("debug", func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewConsoleEncoder(encoderConfig), nil
	})
	if err != nil {
		panic(err)
	}

	InitializeAccessLogger(os.Stdout)
	appLogLevel := zapcore.InfoLevel
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		appLogLevel = NewLevel(level)
	}
	InitializeAppLogger(os.Stderr, appLogLevel)
}

func InitializeAccessLogger(w io.Writer) {
	Access = NewZapLogger(nil, []io.Writer{w}, zapcore.InfoLevel)
}

func InitializeAppLogger(w io.Writer, logLevel zapcore.Level) {
	App = NewZapLogger(nil, []io.Writer{w}, logLevel)
}

func NewLevel(level string) zapcore.Level {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	case "panic":
		l = zap.PanicLevel
	case "fatal":
		l = zap.FatalLevel
	default:
		l = zap.InfoLevel
	}
	return l
}

func NewZapLogger(
	encoderConfig *zapcore.EncoderConfig, writers []io.Writer, logLevel zapcore.Level, options ...zap.Option,
) *zap.Logger {
	if encoderConfig == nil {
		c := zap.NewProductionEncoderConfig()
		c.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig = &c
	}
	if len(writers) == 0 {
		writers = append(writers, ioutil.Discard)
	}
	enabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})
	cores := make([]zapcore.Core, len(writers))
	for i, w := range writers {
		cores[i] = zapcore.NewCore(zapcore.NewJSONEncoder(*encoderConfig), zapcore.AddSync(w), enabler)
	}
	return zap.New(zapcore.NewTee(cores...), options...)
}

type consoleEncoder struct {
	zapcore.Encoder
	consoleEncoder zapcore.Encoder
}

func NewConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	// TODO: import "github.com/fatih/color"
	//color.NoColor = false // Force enabled

	cfg.StacktraceKey = ""
	cfg2 := cfg
	cfg2.NameKey = ""
	cfg2.MessageKey = ""
	cfg2.LevelKey = ""
	cfg2.CallerKey = ""
	cfg2.StacktraceKey = ""
	cfg2.TimeKey = ""
	return consoleEncoder{
		consoleEncoder: zapcore.NewConsoleEncoder(cfg),
		Encoder:        zapcore.NewJSONEncoder(cfg2),
	}
}
