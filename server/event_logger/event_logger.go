package event_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New(l *zap.Logger) *Logger {
	return &Logger{logger: l}
}

func (l *Logger) Log(
	userID uint32,
	category,
	action string,
	fields ...zapcore.Field,
) {
	f := make([]zapcore.Field, 0, len(fields)+1)
	f = append(
		f,
		zap.String("category", category),
		zap.String("action", action),
		zap.Uint("userID", uint(userID)),
	)
	f = append(f, fields...)
	l.logger.Info("eventLog", f...)
}
