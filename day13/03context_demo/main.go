// @Author huzejun 2022/12/15 17:12:00
package main

import (
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup

//全局变量的方式

var exit bool

func worker() {
	defer wg.Done()
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		//如何接收外部命令实现退出
		if exit {
			break
		}
	}
	//如何接收外部命令实现退出
	//wg.Done()
}

func main() {
	wg.Add(1)
	go worker()
	// 如何优雅的实现子goroutine
	time.Sleep(time.Second * 5)
	exit = true
	wg.Wait() //等待...
	fmt.Println("over")
}
