package main

import (
	"code/log_transfer/config"
	"code/log_transfer/es"
	"code/log_transfer/kafka"
	"fmt"
)

// log transfer
// 1. 从kafka消费日志数据 写入es
// 2.

func main() {

	//配置文件初始化
	cfg := config.Init()

	// kafka 初始化
	err := kafka.InitKafka([]string{cfg.KafkaHost + cfg.KafkaPort}, cfg.Topic)
	if err != nil {
		fmt.Printf("connect to kafka failed, err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init kafka success")

	// 写入 es
	err = es.InitES(cfg.EsHost+cfg.EsPort, cfg.Index, cfg.MaxChanSize, cfg.GoroutineNum)
	if err != nil {
		fmt.Printf("Init es failed, err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init es success")

	// 在这儿停顿
	select {}

}
