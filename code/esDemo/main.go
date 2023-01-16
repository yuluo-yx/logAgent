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
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("connect to es success")
	// 创建一条数据
	p1 := Person{Name: "mumu", Age: 21, Married: false}
	put1, err := client.Index().
		// 指定数据库
		Index("user").
		// 把对象转成json
		BodyJson(p1).
		// 插入
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
