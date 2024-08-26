package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger
var LoggerSlave *zap.SugaredLogger
var LoggerAccess *zap.Logger

var path = "/Users/markov/Documents/code/go_code/vpeel/log/"

func InitLogger() {
	writeSyncerError := getLogWriter(path + "info.log")
	writeSyncerMonitor := getLogWriter(path + "error.log")
	writeSyncerSlave := getLogWriter(path + "slave.log")
	writeSyncerAccess := getLogWriter(path + "access.log")
	encoder := getEncoder()

	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel && lv <= zapcore.WarnLevel
	})

	coreError := zapcore.NewCore(encoder, writeSyncerError, infoLevel)
	coreMonitor := zapcore.NewCore(encoder, writeSyncerMonitor, zapcore.ErrorLevel)
	coreSlave := zapcore.NewCore(encoder, writeSyncerSlave, zapcore.InfoLevel)
	coreAccess := zapcore.NewCore(encoder, writeSyncerAccess, zapcore.InfoLevel)

	core := zapcore.NewTee(coreError, coreMonitor)
	Logger = zap.New(core, zap.AddCaller()).Sugar()
	LoggerSlave = zap.New(coreSlave, zap.AddCaller()).Sugar()
	LoggerAccess = zap.New(coreAccess, zap.AddCaller())
}

func SyncLogger() {
	if Logger != nil {
		Logger.Sync()
	}
	if LoggerAccess != nil {
		LoggerAccess.Sync()
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1,
		MaxBackups: 500,
		MaxAge:     7,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
