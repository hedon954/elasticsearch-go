package elasticsearch_go

import (
	"github.com/hedon954/go-pkg/zap"
	zap1 "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLogger *zap1.Logger

func init() {
	initLogger()
}

func initLogger() {
	ZapLogger = zap.NewLogger(zap.NewStdoutPlugin(zapcore.ErrorLevel))
}
