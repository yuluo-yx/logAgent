package core

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context // 可以通过context停止go routine
	cancel context.CancelFunc
}

// tailTask 构造函数 根绝topic和path造一个tailTask对象
func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}

	return tt
}

// Init 使用tail包打开日志文件准备读
func (t *tailTask) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, config)

	return
}

// 核心业务 收集日志发往kafka
func (t *tailTask) run() {
	logrus.Infof("collect for path:%s is running.....", t.path)

	// 循环读数据
	for {
		select {
		case <-t.ctx.Done(): // 只要调用了t.cancel() 就会收到信号，从而停止go routine
			logrus.Infof("path:%s is stopping...", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen, filename: %s\n", t.tObj.Filename)
				// 读取出错等一秒
				time.Sleep(time.Second)
				continue
			}

			// 如果是空行就跳过
			fmt.Printf("%#v\n", line.Text)
			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Info("出现空行，直接跳过了……")
				continue
			}

			// 利用通道将同步的代码改成异步的，
			//将读出的一行日志包装成kafka里面的msg类型
			//fmt.Println("msg: ", line.Text)
			msg := &sarama.ProducerMessage{}
			msg.Topic = "web_log"
			msg.Value = sarama.StringEncoder(line.Text)
			// 放到通道中
			MsgChan <- msg
		}

	}
}
