package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// 将日志数据写入Elasticsearch

type ESClient struct {
	client      *elastic.Client
	index       string
	logDataChan chan interface{}
}

var (
	esClient *ESClient
)

func InitES(addr, index string, goroutineNum, maxSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + addr))
	if err != nil {
		// Handle error
		panic(err)
	}
	//fmt.Printf("%#v\n", client)

	// 初始化esClient 避免空指针异常
	esClient = &ESClient{
		client:      client,
		index:       index,
		logDataChan: make(chan interface{}, maxSize),
	}
	fmt.Println("connect to es success")
	// 从通道中取出数据,写入到es中去
	for i := 0; i < goroutineNum; i++ {
		go sendToES()
	}
	return
}

func sendToES() {
	for m1 := range esClient.logDataChan {
		//b, err := json.Marshal(m1)
		//if err != nil {
		//	fmt.Printf("marshal m1 failed, err:%v\n", err)
		//	continue
		//}
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// PutLogData 通过一个首字母大写的函数将通道暴露出去 从包外接收msg,发送到chan中
func PutLogData(msg interface{}) {
	esClient.logDataChan <- msg
}
