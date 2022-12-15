// @Author huzejun 2022/12/15 17:12:00
package main

import (
	"context"
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup

//使用channel的方式实现

func worker(ctx context.Context) {
	defer wg.Done()
	go worker2(ctx)
LABEL:
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LABEL
		default:

		}
	}
	//如何接收外部命令实现退出
}

func worker2(ctx context.Context) {
LABEL:
	for {
		fmt.Println("worker2...")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LABEL
		default:

		}
	}
	//如何接收外部命令实现退出
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go worker(ctx)
	// 如何优雅的实现子goroutine
	time.Sleep(time.Second * 5)
	cancel()  //调用cancel函数告诉子goroutine退出
	wg.Wait() //等待...
	fmt.Println("cancel...")
	time.Sleep(time.Second * 5)
	fmt.Println("over")
}
