package logger

import (
	"testing"

	"bytes"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	//r := require.New(t)
	var b bytes.Buffer
	e := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(e, zapcore.AddSync(&b), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	}))
	l := zap.New(core)
	//r.Nil(err)
	l.Info("TestLogger !")
	_ = l.Sync()
	//l.C
	fmt.Printf("b = %s\n", b.String())
}
