package ioc

import (
	"blog/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitialLogger() logger.LoggerV1 {
	config := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &config)
	if err != nil {
		panic(err)
	}
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
