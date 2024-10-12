package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go setData(ch)
	go getData(ch)

	time.Sleep(1e9)
}

func setData(ch chan string) {
	//
	ch <- "Hi~"
	ch <- "this"
	ch <- "is"
	ch <- "Beijing"
	ch <- "Where are you?"
}

// ch <- int1 表示：用通道 ch 发送变量 int1（双目运算符，中缀 = 发送）
func getData(ch chan string) {
	var input string
	// time.Sleep(2e9)
	for {
		input = <-ch // input 从信道中解析内容
		fmt.Printf("%s ", input)
	}
}
