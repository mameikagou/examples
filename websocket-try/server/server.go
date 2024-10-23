package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 传递地址格式  go run server.go -a :2222
// 9000是默认值
var a = flag.String("a", ":8080", "这是-help的提示")

var upgrader = websocket.Upgrader{}

func echo(ctx *gin.Context, done chan struct{}) {
	writer, request := ctx.Writer, ctx.Request
	conn, err := upgrader.Upgrade(writer, request, nil) // 更新成websocket.Upgrader的方法

	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage() // 读传过来的东西
		if err != nil {
			log.Println("read:", err)
			break
		}
		if err := conn.WriteMessage(messageType, message); err != nil { // 原样将读取的东西发回给客户端
			log.Println("write:", err)
			break
		}
		log.Printf("这是server recv :%s", message)
	}
	// 连接关闭时发送 done 信号
	done <- struct{}{}
}

func main() {
	flag.Parse()
	r := gin.Default()
	done := make(chan struct{})
	r.GET("/echo", func(ctx *gin.Context) {
		echo(ctx, done)
	})
	go func() {
		<-done
		log.Println("WebSocket 连接已关闭")
	}()
	log.Fatal(r.Run(*a))
}

// log.Fatal
// 记录日志消息：将传递给 log.Fatal 的参数格式化为一条日志消息，并将其输出到标准错误输出（stderr）。
// 终止程序：调用 os.Exit(1) 终止程序运行，返回状态码 1，表示程序异常退出。
