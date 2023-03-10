// @Author huzejun 2023/1/4 21:13:00
package kafka

import (
	"code.dream.com/studygolang/day16/log_transfer/es"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

//1、初始化kafka连接

//2、从kafka里面取出日志数据

func Init(addr []string, topic string) (err error) {
	//创建新的消费者
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { //遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err)
			return
		}
		//defer pc.AsyncClose()
		//异步从每个分区消费信息
		fmt.Println("start to consume...")
		go func(sarama.PartitionConsumer) {
			fmt.Println("in sarama.PartitionConsumer")
			for msg := range pc.Messages() {
				//logDataChan<-msg //为了将同步流程异步化，所以将取出的日志数据先放到channel中
				fmt.Println(msg.Topic, string(msg.Value))
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value, &m1)
				if err != nil {
					fmt.Printf("unmarshal msg failed,err:%v\n", err)
					continue
				}
				es.PutLogData(m1)
			}
		}(pc)
	}
	return
}
