package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() {
	l, err := zap.NewProduction() //json structured logger
	if err != nil {
		panic(err)
	}

	Log = l
}

func Sync() {
	_ = Log.Sync()
}
