// @Author huzejun 2022/12/15 17:12:00
package main

import (
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup

//初始的例子

func worker() {
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
	}
	//如何接收外部命令实现退出
	wg.Done()
}

func main() {
	wg.Add(1)
	go worker()
	// 如何优雅的实现子goroutine
	wg.Wait()
	fmt.Println("over")
}
