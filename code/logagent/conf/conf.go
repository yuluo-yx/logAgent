package conf

import (
	"code/logagent/utils"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type Conf struct {
	KafkaHost   string `ini:"kafkaHost"`
	KafkaPort   string `ini:"kafkaPort"`
	Topic       string `ini:"topic"`
	LogFilePath string `ini:"logFilePath"`
	ChanSize    int64  `ini:"chanSize"`
	EtcdHost    string `ini:"etcdHost"`
	EtcdPort    string `ini:"etcdPort"`
	CollectKey  string `ini:"collectKey"`
}

func Init() (conf *Conf) {
	// 配置对象
	var config = new(Conf)
	// 本地读取环境变量  注意：此路径是从main.go开始算起的路径
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		// 处理读取配置文件异常
		// panic 直译为 运行时恐慌 当panic被抛出异常后，如果我们没有在程序中添加任何保护措施的话，程序就会打印出panic的详细情况之后，终止运行
		logrus.Error("配置文件读取错误，请检查文件路径:", err)
		return
	}

	LoadKafKaConfig(config, file)
	LoadLogFileConfig(config, file)
	LoadEtcdConfig(config, file)

	//fmt.Printf("%#v\n", config)
	return config
}

// LoadEtcdConfig 加载etcd的配置
func LoadEtcdConfig(conf *Conf, file *ini.File) {
	conf.EtcdHost = file.Section("etcd").Key("etcdHost").String()
	conf.EtcdPort = file.Section("etcd").Key("etcdPort").String()
	conf.CollectKey = file.Section("etcd").Key("collectKey").String()

	// 替换本机ip地址
	ip, err := utils.GetOutboundIP()
	if err != nil {
		logrus.Errorf("load config failed, err:%v", err)
		return
	}
	conf.CollectKey = fmt.Sprintf(conf.CollectKey, ip)
	//fmt.Printf("配置中读取到IP地址：%s", conf.CollectKey)

}

// LoadLogFileConfig 加载日志相关配置
func LoadLogFileConfig(conf *Conf, file *ini.File) {
	conf.LogFilePath = file.Section("collect").Key("logfilePath").String()
}

// LoadKafKaConfig 加载kafka相关配置
func LoadKafKaConfig(conf *Conf, file *ini.File) {
	conf.KafkaHost = file.Section("kafka").Key("kafkaHost").String()
	conf.KafkaPort = file.Section("kafka").Key("kafkaPort").String()
	conf.Topic = file.Section("kafka").Key("topic").String()
	conf.ChanSize, _ = file.Section("kafka").Key("chanSize").Int64()
}
