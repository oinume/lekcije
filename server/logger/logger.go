package logger

import (
	"io"
	"os"

	"github.com/uber-go/zap"
)

var (
	AccessLogger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stdout))
	AppLogger    = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stderr))
)

func InitializeAccessLogger(writer io.Writer) {
	AccessLogger = zap.New(
		zap.NewJSONEncoder(zap.RFC3339Formatter("ts")),
		zap.Output(zap.AddSync(writer)),
	)
}

func InitializeAppLogger(writer io.Writer) {
	AppLogger = zap.New(
		zap.NewJSONEncoder(zap.RFC3339Formatter("ts")),
		zap.Output(zap.AddSync(writer)),
	)
}
