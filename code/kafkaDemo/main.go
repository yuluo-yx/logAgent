package main

import (
	"fmt"
	_ "fmt"
	"github.com/Shopify/sarama"
	"sync"
)

func main() {

	/*kafka生产者*/
	//kafkaProducer()

	/*kafka消费者*/
	kafkaConsumer()

}

func kafkaProducer() {
	config := sarama.NewConfig()

	/*生产者配置*/
	// ACK
	config.Producer.RequiredAcks = sarama.WaitForAll
	//分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//确认
	config.Producer.Return.Successes = true

	/*链接kafka*/
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err: ", err)
		return
	}
	defer client.Close()

	/*封装消息*/
	msg := &sarama.ProducerMessage{}
	msg.Topic = "shopping"
	msg.Value = sarama.StringEncoder("2022/12/30 20:09")

	/*发送消息*/
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err: ", err)
		return
	}
	fmt.Printf("pid： %v offset: %v\n", pid, offset)
}

func kafkaConsumer() {
	// 创建新的消费者
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	var wg sync.WaitGroup
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s",
					msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}
