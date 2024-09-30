package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(
		context.Background(), // 返回一个空的 context.Context; 是基础的上下文环境;
		syscall.SIGINT,       // 通常由Ctrl + C触发
		syscall.SIGTERM,
	)
	defer stop()

	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	//使用autotls.RunWithContext函数来运行gin框架的路由器r。这个函数会自动处理 TLS 相关的配置。
	//它会尝试为example1.com和example2.com这两个域名（这里应该是示例域名，实际使用时需要替换为真实的域名）获取 TLS 证书，并将gin框架的服务运行在 HTTPS 模式下。
	log.Fatal(autotls.RunWithContext(ctx, r, "example1.com", "example2.com"))
}
