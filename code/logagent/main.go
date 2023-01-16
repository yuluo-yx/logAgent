package main

import (
	"code/logagent/conf"
	"code/logagent/core"
	"code/logagent/etcd"
	"fmt"
	"github.com/sirupsen/logrus"
)

/* 日志收集客户端，类似的开源项目 fileBeat */

func run() {
	select {}
}

func main() {

	//配置文件初始化
	config := conf.Init()

	// 连接kafka
	err := core.InitKafka([]string{config.KafkaHost + config.KafkaPort}, config.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	// 初始化etcd连接
	//从etcd中拉取要收集的日志文件配置
	err = etcd.Init([]string{config.EtcdHost + config.EtcdPort})
	if err != nil {
		logrus.Errorf("init etcd failed, err:%v\n", err)
		return
	}
	// 从etcd中拉取最新的日志收集配置
	allConf, err := etcd.GetConf(config.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v\n", err)
		return
	}
	// 输出从etcd中获取的配置信息
	fmt.Println(allConf)
	// 派小弟 监控etcd中对应key的变换情况 不能阻塞 用go关键字修饰
	go etcd.WatchConf(config.CollectKey)

	// 连接tailf  根据从etcd中获取的配置项初始化
	err = core.InitTail(allConf)
	if err != nil {
		logrus.Errorf("init tail failed, err:%v", err)
		return
	}
	logrus.Info("init tail success!")

	run()
}
