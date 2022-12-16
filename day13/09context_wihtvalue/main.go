package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// context.WithValue

type TraceCode string
type UserID string

var wg sync.WaitGroup

func worker(ctx context.Context) {
	key := TraceCode("TRACE_CODE")           //造一个TraceCode类型的结束 n := 100 n:=int32（100）str := "xxx"
	traceCode, ok := ctx.Value(key).(string) // 在子goroutine中获取trace code
	if !ok {
		fmt.Println("invalid trace code")
	}

	useridKey := UserID("USER_ID")
	userid, ok := ctx.Value(useridKey).(int64) // 在子goroutine中获取user ID
	if !ok {
		fmt.Println("invalid user id")
	}
	log.Printf("%s worker func...", traceCode)
	log.Printf("userid:%d", userid)
LOOP:
	for {
		fmt.Printf("worker,trace code:%s\n", traceCode)
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func main() {
	//设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	//在系统的入口中设置 trace code传递给后续启动的goroutine实现日志数据聚合
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "123123132")
	ctx = context.WithValue(ctx, UserID("USER_ID"), int64(142134123153425))
	log.Printf("%s main 函数", "123123132")
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
