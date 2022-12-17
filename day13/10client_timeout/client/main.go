package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// 客户端
type respData struct {
	resp *http.Response
	err  error
}

func doCall(ctx context.Context) {
	// 造一个客户端
	transport := http.Transport{
		/*DisableKeepAlives: true,*/
	}
	client := http.Client{
		Transport: &transport,
	}

	respChan := make(chan *respData, 1)
	// 造一个请求对象
	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	if err != nil {
		fmt.Printf("new request failed,err:%v\n", err)
		return
	}
	req = req.WithContext(ctx) //敷衍带超时的ctx创建一个新的client request
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		resp, err := client.Do(req) // 执行真正的发请求的操作
		fmt.Printf("client.do resp:%v,err:%v\n", resp, err)
		rd := &respData{
			resp: resp,
			err:  err,
		}
		respChan <- rd
		wg.Done()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("call aip timeout")
	case result := <-respChan:
		fmt.Println("call server api success")
		if result.err != nil {
			fmt.Printf("call server api failed, err:%v\n", result.err)
			return
		}
		defer result.resp.Body.Close()
		data, _ := ioutil.ReadAll(result.resp.Body)
		fmt.Printf("resp:%v\n", string(data))
	}
}

func main() {
	//定义一个100毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel() //调用cancel释放子ogroutine资源
	doCall(ctx)    // 新造一个client,调用client去请求本地的127.0.0.1：8000 把响应的结果放到一个chan中，使用select
}
