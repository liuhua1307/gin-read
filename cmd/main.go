package main

import (
	"context"
	"github.com/liuhua1307/gin-read/internal/configs"
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	logger "github.com/liuhua1307/gin-read/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configs.LogInit("/log")
	server := configs.ServerInit()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	CloseRouter(server)
}

func CloseRouter(srv *http.Server) {
	//创建一个信道捕获系统的信息，buffered channel
	quit := make(chan os.Signal, 1)
	//通过该信道监听这两种信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	//没收到信号这里就会堵塞，不会接着往下执行
	<-quit
	log.Log().Info("Shutting down server...")
	//设置关于此ctx的协程都会在5秒超时时间后自动强制调用关闭函数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//这里用于在5秒内就执行完毕下面的语句后调用结束协程
	defer cancel()
	//这里新创建一个协程，用于关闭服务器,且绑定于上面这个子ctx上
	if err := srv.Shutdown(ctx); err != nil {
		log.Log().Fatal("Server forced to shut down", logger.Field{Key: "error", Value: err})
	}
	log.Log().Info("Server exited")
}
