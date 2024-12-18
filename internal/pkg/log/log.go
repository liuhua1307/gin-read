package log

import (
	logger "github.com/liuhua1307/gin-read/pkg/log"
)

var log logger.Logger

func RegisterLog(logger logger.Logger) {
	log = logger
}

func Log() logger.Logger {
	if log == nil {
		panic("implement not found for interface Logger, please register")
	}
	return log
}
