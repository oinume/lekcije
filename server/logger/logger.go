package logger

import (
	"os"

	"github.com/uber-go/zap"
)

var (
	AccessLogger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stdout))
	AppLogger    = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stderr))
)
