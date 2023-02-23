package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
    "os"
    "time"
)

/**
 * @Author zyq
 * @Date 2023/2/23 11:28 AM
 * @Description
 **/
var (
    logger     *zap.Logger
    Sugar      *zap.SugaredLogger
    debugLevel = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level == zap.DebugLevel
    })
)

const timeFormatDefault = "2006-01-02 15:04:05.000"

func init() {
    var err error
    options := []zap.Option{zap.AddCaller()}
    options = append(options, zap.AddCallerSkip(1))
    logger = zap.New(zapcore.NewTee(zapcore.NewCore(zapcore.NewConsoleEncoder(NewEncoderConfig()),
        zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
        debugLevel,
    )), options...)
    if err != nil {
        log.Fatal(err)
    }
    Sugar = logger.Sugar()
}

func NewEncoderConfig() zapcore.EncoderConfig {
    return zapcore.EncoderConfig{
        // Keys can be anything except the empty string.
        TimeKey:        "T",
        LevelKey:       "L",
        NameKey:        "N",
        CallerKey:      "C",
        MessageKey:     "M",
        StacktraceKey:  "S",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     TimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.Format(timeFormatDefault))
}

func Debug(args ...interface{}) {
    Sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
    Sugar.Debugf(template, args...)
}
