package main

import (
	"fmt"
	"go.etcd.io/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 代码连接etcd

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err:%v", err)
		return
	}
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	str := `[{"path":"d:/logs/s5.log","topic":"s5_log"},{"path":"e:/logs/web.log","topic":"web_log"}]`
	//str := `[{"path":"d:/logs/s5.log","topic":"s5_log"},{"path":"e:/logs/web.log","topic":"web_log"},{"path":"f:/logs/naza.log","topic":"naza_log"}]`
	//str := `[{"path":"d:/logs/s5.log","topic":"s5_log"},{"path":"e:/logs/web.log","topic":"web_log"},{"path":"f:/logs/naza.log","topic":"naza_log"},{"path":"f:/logs/naza3.log","topic":"efg_log"}]`
	//_, err = cli.Put(ctx, "collect_log_conf", str)
	_, err = cli.Put(ctx, "collect_log_192.168.32.110_conf", str)
	//_, err = cli.Put(ctx, "collect_log_conf", "很好") collect_log_192.168.32.110_conf
	//_, err = cli.Put(ctx, "s4", "很好")
	if err != nil {
		fmt.Printf("put to etcd failed:err:%v", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//gr, err := cli.Get(ctx, "s4") // collect_log_conf
	//gr, err := cli.Get(ctx, "collect_log_conf") // collect_log_conf
	gr, err := cli.Get(ctx, "collect_log_192.168.32.110_conf") // collect_log_conf
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v", err)
		return
	}
	for _, ev := range gr.Kvs {
		fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
	}
	cancel()
}
