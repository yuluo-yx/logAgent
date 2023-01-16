# Go 日志收集项目

## 环境搭建

kafka  zoopeeker

> 修改zookeeper的data配置
>
> 修改kafka的zookeeper配置和server中listeners=PLAINTEXT://:9092 改为 \#listeners=PLAINTEXT://192.0.0.1:9092

`bin\windows\zookeeper-server-start.bat config\zookeeper.properties`

`bin\windows\kafka-server-start.bat config\server.properties`

安装`sarama`

```go
go get github.com/Shopify/sarama 
```

**注意：这里必须安装v1.19版本的sarama，v1.20之后加入了zstd压缩算法，需要使用cgo，在windowsi下编译会报如下错误：`exec: "gcc":executable fie not in %PATH%`**

在go.mod中指定版本，使用`go mod download`下载依赖

启动消费者获取数据

`bin\windows\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic shopping --from-beginning`

代码：

```go
package main

import (
   "fmt"
   _ "fmt"
   "github.com/Shopify/sarama"
)

func main() {
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
```

tail包的使用：读取日志文件

`go get github.com/hpcloud/tail`

使用的第三方库

```go
go get github.com/go-ini/ini   读取配置文件
go get github.com/sirupsen/logrus  日志依赖
```



## 日志收集（logAgent）项目

### 配置文件版本

#### ini配置文件读取

##### 下载

##### 代码

#### kafka初始化

#### tail初始化

#### 小结

#### run

#### 总结

> 1. 尽量不要暴露太多的变量
> 2. 使用tailf读日志文件，发送到kafka
> 3. 使用了chan和goroutine
> 4. 目前存在的问题：
>    - 只能读取一个日志文件，不支持多个日志文件读取
>    - topic配置单一，实际业务可能有多个topic
>    - 有多少个日志路径，应该就有多少个tailObj
> 5. 解决存在的问题

## logAgent引入etcd

#### etcd介绍

https://docs.qq.com/doc/DTndrQXdXYUxUU09O?opendocxfrom=admin

类似于zookeeper(java中用得到)，etcd\consul（go中用得到）***注册、配置中心***

#### etcd使用

##### 安装

https://github.com/etcd-io/etcd/releases

##### 启动

命令行输入 `SET ETCECTL_API=3`使用v3版本。v2版本没有put命令

双击etcd.exe

如果出现`no help topic for 'put'`，就是没有设置环境变量

命令行：

put 设置值 `etcdctl.exe --endpoints=http://127.0.0.1:2379 put yuluo "test"`

get 获取值 `etcdctl.exe --endpoints=http://127.0.0.1:2379 get yuluo`

##### 使用go操作etcd

`go get go.etcd.io/etcd/clientv3`

如果安装失败，在terminal中运行

```bash
go mod edit -require=google.golang.org/grpc@v1.26.0

go get -u -x google.golang.org/grpc@v1.26.0

之后重新 go get
```

demo

```go
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
   _, err = cli.Put(ctx, "s4", "真好")
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
```

从cmd中获取

```bash
etcdctl.exe --endpoints=http://127.0.0.1:2379 get s4
```

#### watch 操作

> 监控etcd中一个key的变化，更改，创建，删除……

```go
package main

import (
   "context"
   "fmt"
   "go.etcd.io/etcd/clientv3"
   "time"
)

/* 演示 etcd的watch 操作 */

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

   //watch
   wCh := cli.Watch(context.Background(), "s4")
   for wresp := range wCh {
      for _, evt := range wresp.Events {
         fmt.Printf("type:%s key:%s, value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
      }
   }
}
```

#### kafka消费实例

```java
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
```

#### 从etcd加载日志配置项

```go
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
```

#### watch日志配置项变更

监听日志收集路径的配置变化

#### 管理日志收集任务

在程序启动之后，拉取了最新的配置之后，就应该派一个小弟去监控etcd中`collect_log_conf`这个key的变化情况

### 遗留问题

如果logAgent项目停止，需要记录读取的日志位置，下次启动之后，接着上次继续读取，参考filebeat实现

 

## 下一步工作

每台服务器上的logagent的收集项可能都不一致，我们需要让logagent去etcd中获取自己的配置

### go获取本机ip

```go
// GetOutboundIP 获取本机ip
func GetOutboundIP() (ip string, err error) {
   conn, err := net.Dial("udp", "8.8.8.8:80")
   if err != nil {
      return
   }
   defer func(conn net.Conn) {
      err := conn.Close()
      if err != nil {
         logrus.Error("net conn err:%s", err)
      }
   }(conn)

   localAddr := conn.LocalAddr().(*net.UDPAddr)
   //fmt.Println("本机ip地址：", localAddr.String())
   ip = strings.Split(localAddr.IP.String(), ":")[0]
   return
}
```

### go获取系统信息 gopsutil

> psutil 是一个python库。gopsutil是go的实现
>
> 采集一些系统和监控信息等
>
> `go get github.com/shirou/gopsutil`

### influxDB时序数据库

监控数据库，按时间做索引。go语言编写的。实现分布式和水平伸缩扩展。支持类sql语言查询

#### 名词介绍

| influxDB名词 | 传统数据库概念 |
| :----------: | :------------: |
|   database   |     数据库     |
| measurement  |     数据表     |
|    point     |     数据行     |

#### point

influxDB中的point相当于传统数据库里的一行数据，由时间戳（time）、数据（field）、标签（tag）组成。

| Point属性 |                传统数据库概念                |
| :-------: | :------------------------------------------: |
|   time    |     每个数据记录时间，是数据库中的主索引     |
|   field   | 各种记录值（没有索引的属性），例如温度、湿度 |
|   tags    |       各种有索引的属性，例如地区、海拔       |

#### Series

`Series`相当于是 InfluxDB 中一些数据的集合，在同一个 database 中，retention policy、measurement、tag sets 完全相同的数据同属于一个 series，同一个 series 的数据在物理上会按照时间顺序排列存储在一起。

#### 创建数据库操作

> 类似mysql
>
> `show databases;`
>
> `create database test`
>
> 查询：
>
> `select * from cpu_percent order by time desc limit 10`
>
> 删除：
>
> ```
>  delete from cpu where time = 1618456817804394700
>  delete from cpu where "tag"="xxx" 
> ```

#### 操作示例 1.x 版本

```go
package main

//influxDB 数据库操作示例

import (
   "fmt"
   client "github.com/influxdata/influxdb1-client/v2"
   "log"
   "time"
)

// influxdb demo

// 数据库链接
func connInflux() client.Client {
   cli, err := client.NewHTTPClient(client.HTTPConfig{
      Addr:     "http://127.0.0.1:8086",
      Username: "admin",
      Password: "",
   })
   if err != nil {
      log.Fatal(err)
   }
   return cli
}

// query 查询
func queryDB(cli client.Client, cmd string) (res []client.Result, err error) {
   q := client.Query{
      Command:  cmd,
      Database: "test",
   }
   if response, err := cli.Query(q); err == nil {
      if response.Error() != nil {
         return res, response.Error()
      }
      res = response.Results
   } else {
      return res, err
   }
   return res, nil
}

// insert 插入
func writesPoints(cli client.Client) {
   bp, err := client.NewBatchPoints(client.BatchPointsConfig{
      Database:  "test",
      Precision: "s", //精度，默认ns
   })
   if err != nil {
      log.Fatal(err)
   }
   tags := map[string]string{"cpu": "ih-cpu"}
   fields := map[string]interface{}{
      "idle":   201.1,
      "system": 43.3,
      "user":   86.6,
   }

   pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
   if err != nil {
      log.Fatal(err)
   }
   bp.AddPoint(pt)
   err = cli.Write(bp)
   if err != nil {
      log.Fatal(err)
   }
   log.Println("insert success")
}

func main() {
   conn := connInflux()
   fmt.Println(conn)

   // insert
   writesPoints(conn)

   // 获取10条数据并展示
   qs := fmt.Sprintf("SELECT * FROM %s LIMIT %d", "cpu_usage", 10)
   res, err := queryDB(conn, qs)
   if err != nil {
      log.Fatal(err)
   }

   for _, row := range res[0].Series[0].Values {
      for j, value := range row {
         log.Printf("j:%d value:%v\n", j, value)
      }
   }
}
```

### grafana

> 展示数据工具，通常用来做监控数据可视化。

#### 配置

在 Windows 上，该`sample.ini`文件位于与`defaults.ini`文件相同的目录中。它包含注释掉的所有设置。复制`sample.ini`并命名`custom.ini`。

Grafana 使用分号（`;`char）来注释`.ini`文件中的行。您必须通过从该行的开头删除来取消注释您正在修改`custom.ini`的文件中的每一行。否则您的更改将被忽略。`grafana.ini``;`

url : `http://127.0.0.1:3000` 默认为 admin/admin

1. 连接时序数据库
2. 创建dashborad
   1. 在面板中创建数据库查询语句，之后显示在面板上
   2. 设置显示图形

### ES介绍

Elasticsearch（ES）是一个基于Lucene构建的开源、分布式、RESTful接口的全文搜索引擎。Elasticsearch还是一个分布式文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索，ES能够横向扩展至数以百计的服务器存储以及处理PB级的数据。可以在极短的时间内存储、搜索和分析大量的数据。通常作为具有复杂搜索场景情况下的核心发动机。

- 电商搜索商品
- github所有
- 收集日志或者交易数据。对数据进行分析，挖掘等

> **需要在配置文件中修改path路径和安全认证开关值改为false，实现免密登录。不然连接报错！**

可以使用postman或浏览器发送请求到`127.0.0.1:9200/_cat/health?v`查看健康状态

```
epoch      timestamp cluster       status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1673774111 09:15:11  elasticsearch green           1         1      2   2    0    0        0             0                  -                100.0%
```

### Kibana

是一个分析和可视化的平台，设计用于和es一起工作。可以用来搜索，查看，并和存储在es中的数据进行交互

### logtransfer 实现

  ![image-20230116120530979](E:\Go\project\日志收集项目\笔记+总结\log收集.assets\image-20230116120530979.png)

## 项目架构图

![image-20230116120604714](E:\Go\project\日志收集项目\笔记+总结\log收集.assets\image-20230116120604714.png)