package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Access = NewZapLogger(&zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
	App = NewZapLogger(&zap.Config{
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	})
)

func init() {
	err := zap.RegisterEncoder("debug", func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewConsoleEncoder(encoderConfig), nil
	})
	if err != nil {
		panic(err)
	}
	// TODO: level
	//if !config.IsProductionEnv() {
	//	App.SetLevel(zap.DebugLevel)
	//}
}

func InitializeAccessLogger() {
	// TODO: OutputPaths
	Access = NewZapLogger(&zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	})
}

func InitializeAppLogger() {
	App = NewZapLogger(&zap.Config{
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	})
	// TODO: level
	//if !config.IsProductionEnv() {
	//	App.SetLevel(zap.DebugLevel)
	//}
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

func NewZapLogger(config *zap.Config, opts ...zap.Option) *zap.Logger {
	var c zap.Config
	if _, ok := os.LookupEnv("ZAP_DEBUG"); ok {
		c = zap.NewDevelopmentConfig()
		c.Encoding = "debug"
		c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		c = zap.NewDevelopmentConfig()
	}

	c.OutputPaths = config.OutputPaths
	c.ErrorOutputPaths = config.ErrorOutputPaths
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	c.DisableStacktrace = true
	c.Level.SetLevel(zapcore.DPanicLevel) // not show logs normally
	if v := os.Getenv("ZAP_LEVEL"); v != "" {
		var lv zapcore.Level
		if err := lv.UnmarshalText([]byte(v)); err == nil {
			c.Level.SetLevel(lv)
		} else {
			panic(fmt.Sprintf("Unknown zap log level: %v", v))
		}
	}

	l, err := c.Build(opts...)
	if err != nil {
		panic(l)
	}
	return l
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
