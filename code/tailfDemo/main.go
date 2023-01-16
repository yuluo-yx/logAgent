package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	fileName := `E:\Go\project\日志收集项目\logs\test.log`
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	/*打开文件读取数据*/
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Printf("tail %s failed, err: %v\n", fileName, err)
		return
	}

	/*开始读取*/
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines // chan tail.Line
		if !ok {
			fmt.Printf("tail file close reopen, fileName:%s\n", tails.Filename)
			time.Sleep(time.Second) // 读取出错等一秒
			continue
		}
		fmt.Println("msg: ", msg.Text)
	}
}
