// @Author huzejun 2022/12/15 17:12:00
package main

import (
	"fmt"
	"sync"

	"time"
)

/**
make 和 new 的区别
都是用来初始化内存
new多用来为基本数据类型（bool、string、int...)初始化内存，返回的是指针
make用来初始化slice、map、channel,返回的是对应类型。
*/

var wg sync.WaitGroup

//使用channel的方式实现

func worker(ch <-chan struct{}) {
	defer wg.Done()
LABEL:
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		select {
		case <-ch:
			break LABEL
		default:
		}
	}
	//如何接收外部命令实现退出
}

func main() {
	var exitChan = make(chan struct{})

	wg.Add(1)
	go worker(exitChan)
	// 如何优雅的实现子goroutine
	time.Sleep(time.Second * 5)
	exitChan <- struct{}{}
	wg.Wait() //等待...
	fmt.Println("over")
}
