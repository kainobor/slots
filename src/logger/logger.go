package logger

import "go.uber.org/zap"

var lgr *zap.SugaredLogger

func Init() error {
    loggerProd, err := zap.NewProduction()
    if err != nil {
        return err
    }
    defer loggerProd.Sync()
    lgr = loggerProd.Sugar()

    return nil
}

func Log(msg string, args ...interface{}) {
    lgr.Infow(msg, args...)
}

func Err(msg string, args ...interface{}) {
    lgr.Errorw(msg, args...)
}
