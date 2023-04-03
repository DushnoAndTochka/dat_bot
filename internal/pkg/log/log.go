package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var commitID string

var logger *zap.SugaredLogger = nil

func newLogger() {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "/var/log/dushno_and_tochka.log"
	}

	zapConfig := zap.Config{
		// TODO: Настроить уровень логирования и сделать его динамическим в зависимости от среды
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.999999Z07:00"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{logPath},
		ErrorOutputPaths: []string{logPath},
	}

	l, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar().With(zap.String("commit_id", commitID))
}

func GetLogger() *zap.SugaredLogger {
	var once sync.Once
	once.Do(newLogger)
	return logger
}
