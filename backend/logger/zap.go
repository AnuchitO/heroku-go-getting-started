package logger

import (
	"log"
	"os"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func NewZap() (*zap.Logger, func()) {
	// if err := os.MkdirAll("log", os.ModePerm); err != nil {
	// 	log.Fatalf("Error creating backend/log: %v", err)
	// }
	config := loadConfigLog(getEnv(os.Getenv))
	logger, err := config.Build(
		zap.AddCaller(),
	)
	if err != nil {
		log.Fatal(err)
	}
	undo := zap.ReplaceGlobals(logger)

	return logger, func() {
		undo()
		_ = logger.Sync()
	}
}

type getEnv func(string) string

func (fn getEnv) Getenv(key string) string {
	return fn(key)
}

type env interface {
	Getenv(name string) string
}

func loadConfigLog(e env) zap.Config {
	if e.Getenv("ENV") == "LOCAL" {
		return zap.NewDevelopmentConfig()
	} else {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableStacktrace = true
		config.OutputPaths = []string{
			"stdout",
		}
		return config
	}
}
