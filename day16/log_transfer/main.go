package main

import (
	"code.dream.com/studygolang/day16/log_transfer/es"
	"code.dream.com/studygolang/day16/log_transfer/kafka"
	"code.dream.com/studygolang/day16/log_transfer/model"
	"fmt"
	"github.com/go-ini/ini"
)

// log transfer
// 从kafka消费日志数据，写入es

func main() {
	// 1、加载配置文件
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		panic(err)
	}
	fmt.Printf("%#v\n", *cfg)
	fmt.Println("load config success")
	//2、连接es
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.Index, cfg.ESConf.GoNum, cfg.ESConf.MaxSize)
	if err != nil {
		fmt.Printf("Init es failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init ES success")

	// 3、连接kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic)
	if err != nil {
		fmt.Printf("connect to kafka failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init kafka success")
	//3、连接es
	/*	err = es.Init(cfg.ESConf.Address,cfg.ESConf.Index,cfg.ESConf.GoNum,cfg.ESConf.MaxSize)
		if err != nil {
			fmt.Printf("Init es failed,err:%v\n",err)
			panic(err)
		}
		fmt.Println("Init ES success")*/
	//让程序在这儿停顿
	select {}

}
