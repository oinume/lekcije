package log

import (
	"context"
	"errors"
	"io"
	"io/ioutil"

	//"github.com/bugsnag/bugsnag-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger        *zap.Logger
	notifyEnabled bool
}

func New(w io.Writer, notifyEnabled bool) *Logger {
	l := newZapLogger(nil, []io.Writer{w}, zapcore.InfoLevel) // TODO: Specify level from outside
	return &Logger{
		logger:        l,
		notifyEnabled: notifyEnabled,
	}
}

func (l *Logger) Debug(_ context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(_ context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(_ context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
	if !l.notifyEnabled {
		return
	}

	var e error
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			if err, ok := f.Interface.(error); ok {
				e = err
			}
		}
	}
	if e == nil {
		e = errors.New("<no error>") // Set error explicitly because Bugsnag raises error if err is nil
	}
	//l.notifyError(ctx, e, bugsnag.MetaData{"msg": {"msg": msg}})
}

/*
func (l *Logger) notifyError(ctx context.Context, err error, meta bugsnag.MetaData) {
	if err := bugsnag.Notify(err, ctx, meta); err != nil {
		l.logger.Warn("bugsnag.Notify failed", zap.Error(err))
	}
}
*/

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		logger:        l.logger.With(fields...),
		notifyEnabled: l.notifyEnabled,
	}
}

func newZapLogger(
	encoderConfig *zapcore.EncoderConfig,
	writers []io.Writer,
	logLevel zapcore.Level,
	options ...zap.Option,
) *zap.Logger {
	if encoderConfig == nil {
		c := zap.NewProductionEncoderConfig()
		c.LevelKey = "severity"
		c.TimeKey = "timestamp"
		c.EncodeTime = zapcore.RFC3339NanoTimeEncoder
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
