package configs

import (
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	log2 "github.com/liuhua1307/gin-read/pkg/log"
)

func LogInit(path string) {
	// LogInit ...
	log.RegisterLog(log2.NewSlogLogger(path))
}
