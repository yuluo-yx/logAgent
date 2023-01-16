package config

import (
	"code/log_transfer/model"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

func Init() (conf *model.Conf) {
	// 配置对象
	var cfg = new(model.Conf)
	// 本地读取环境变量  注意：此路径是从main.go开始算起的路径
	file, err := ini.Load("./config/logtransfer.ini")
	if err != nil {
		// 处理读取配置文件异常
		// panic 直译为 运行时恐慌 当panic被抛出异常后，如果我们没有在程序中添加任何保护措施的话，程序就会打印出panic的详细情况之后，终止运行
		logrus.Error("配置文件读取错误，请检查文件路径:", err)
		return
	}

	LoadKafKaConfig(cfg, file)
	LoadEsConfig(cfg, file)
	LoadSysConfig(cfg, file)

	//fmt.Printf("%#v\n", config)
	fmt.Printf("config init success! %v", cfg)
	return cfg
}

// LoadEsConfig 加载kafka相关配置
func LoadEsConfig(conf *model.Conf, file *ini.File) {
	conf.EsHost = file.Section("es").Key("esHost").String()
	conf.EsPort = file.Section("es").Key("esPort").String()
	conf.Index = file.Section("es").Key("index").String()
}

// LoadKafKaConfig 加载kafka相关配置
func LoadKafKaConfig(conf *model.Conf, file *ini.File) {
	conf.KafkaHost = file.Section("kafka").Key("kafkaHost").String()
	conf.KafkaPort = file.Section("kafka").Key("kafkaPort").String()
	conf.Topic = file.Section("kafka").Key("topic").String()
}

// LoadSysConfig 加载sys相关配置
func LoadSysConfig(conf *model.Conf, file *ini.File) {
	conf.MaxChanSize, _ = file.Section("sys").Key("maxChanSize").Int()
	conf.GoroutineNum, _ = file.Section("sys").Key("goroutineNum").Int()
}
