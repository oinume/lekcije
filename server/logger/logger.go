package logger

import (
	"os"

	"github.com/uber-go/zap"
)

var (
	AccessLogger = zap.New(zap.NewJSONEncoder(), zap.Output(os.Stdout))
	AppLogger    = zap.New(zap.NewJSONEncoder(), zap.Output(os.Stderr))
)
