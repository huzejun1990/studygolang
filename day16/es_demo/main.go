package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// Elasticsearch demo

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.32.110:9200"))
	if err != nil {
		// Handle err
		panic(err)
	}
	fmt.Println("connect to es success")
	// 创建一条数据
	p1 := Person{Name: "guan", Age: 18, Married: true}
	put1, err := client.Index().
		Index("user").
		BodyJson(p1).
		Do(context.Background())
	if err != nil {
		//Handle error
		panic(err)
	}
	fmt.Printf("indexed user %s to index %s,type %s\n", put1.Id, put1.Index, put1.Type)

}
