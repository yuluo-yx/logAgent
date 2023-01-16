package core

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

// InitKafka 初始化全局kafka连接
func InitKafka(address []string, chanSize int64) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("producer closed, err: ", err)
		return
	}

	// 初始化MsgChan
	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	//使用后台的goroutine从MshChan中读数据
	go sendMsg()

	return
}

// 从MsgChan中读取消息，发送给kafka
func sendMsg() {
	for {
		select {
		case msg := <-MsgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed, err: ", err)
				return
			}
			logrus.Infof("send msg to kafka. pid:%v offset:%v", pid, offset)
		}
	}
}
