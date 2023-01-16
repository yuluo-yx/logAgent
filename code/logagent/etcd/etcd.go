package etcd

/* etcd相关操作 */

import (
	"code/logagent/common"
	"code/logagent/core"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

// GetConf 拉取最新的日志收集配置项
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed, err:%v\n", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s\n", key)
		return
	}

	ret := resp.Kvs[0]
	// ret.Value json格式字符串 反序列化到~
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed, err:%v\n", err)
		return
	}
	return
}

// WatchConf 监控etcd中key的变化情况的函数
func WatchConf(key string) {
	// 循环监听
	for {
		watchCh := client.Watch(context.Background(), key)
		for wresp := range watchCh {
			logrus.Infof("get new conf from etcd!")
			for _, evt := range wresp.Events {
				fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				var newConf []common.CollectEntry
				if evt.Type == clientv3.EventTypeDelete {
					// 如果是删除事件
					logrus.Warning("FBI warning: etcd delete the key!!!")
					core.SendNewConf(newConf)
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf failed, err:%v", err)
					continue
				}
				// 不报错 告诉tailfile这个模块启用新的配置信息了！
				// 使用 channel 没有人接受就是阻塞的
				core.SendNewConf(newConf)
			}
		}
	}
}
