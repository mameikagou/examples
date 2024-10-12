package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go pump(ch)
	go suck(ch)
	fmt.Println(<-ch)
	time.Sleep(1 * time.Second)
}

func pump(ch chan int) {
	for i := 0; ; i++ {
		ch <- i
	}
}
func suck(ch chan int) {
	for {
		fmt.Println(<-ch)
	}
}
