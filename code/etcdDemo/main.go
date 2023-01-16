package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

/*go 操作etcd demo*/

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}

	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// 字符串配置
	str := "[{\"path\":\"e:/Go/project/日志收集项目/logs/test.log\",\"topic\":\"s4_log\"},{\"path\":\"e:/Go/project/日志收集项目/logs/web.log\",\"topic\":\"web_log\"}]"
	//str := "[{\"path\":\"e:/Go/project/日志收集项目/logs/test.log\",\"topic\":\"s4_log\"},{\"path\":\"e:/Go/project/日志收集项目/logs/web.log\",\"topic\":\"web_log\"},{\"path\":\"e:/Go/project/日志收集项目/logs/mumu.log\",\"topic\":\"mumu_log\"}]"

	_, err = cli.Put(ctx, "collect_log_192.168.2.4_conf", str)
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	gr, err := cli.Get(ctx, "s4")
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v", err)
		return
	}
	for _, ev := range gr.Kvs {
		fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
	}
	cancel()
}
