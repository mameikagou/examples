package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var a = flag.String("a", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	// log.SetFlags(0)
	// 参数 0 表示清除所有默认的标志位，即不输出任何额外的信息，如时间戳、文件名和行号等。
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	interrupt := make(chan os.Signal, 1)
	// make(chan os.Signal, 1)：创建一个缓冲通道，用于接收操作系统信号。缓冲区大小为 1，意味着通道可以存储一个信号，防止信号丢失。
	signal.Notify(interrupt, os.Interrupt)
	// 使用 signal.Notify 注册信号处理器，捕获 os.Interrupt（通常是 Ctrl+C)

	u := url.URL{Scheme: "ws", Host: *a, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// Dial 方法：用于建立 WebSocket 连接。
	// 第一个参数 u.String()：是一个 URL 的字符串表示形式，指定了要连接的 WebSocket 服务器地址。
	// 第二个参数 nil：是一个可选的 HTTP 请求头，如果不需要传递额外的请求头，可以传入 nil。

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{}) // 用于通知任务完成

	go func() {
		defer close(done) // 确保任务完成后关闭
		for {
			_, message, err := conn.ReadMessage() //多个goroutine不断读取消息
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("这是client recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		// 等待done通道, ticker 通道和 interrupt 通道
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String())) // 将时间发送给服务端
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// websocket.CloseNormalClosure 用于指示连接正常关闭，没有错误发生。

			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				//case <-done:：如果 done 通道接收到信号（即辅助 goroutine 关闭了 done 通道），主循环退出。
				// case <-time.After(time.Second):：如果 done 通道没有接收到信号，主 goroutine 会等待 1 秒，然后继续执行。
			}
			return
		}
	}
}
