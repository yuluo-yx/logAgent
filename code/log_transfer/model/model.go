package model

type Conf struct {
	KafkaHost    string `ini:"kafkaHost"`
	KafkaPort    string `ini:"kafkaPort"`
	Topic        string `ini:"topic"`
	EsHost       string `ini:"esHost"`
	EsPort       string `ini:"esPort"`
	Index        string `ini:"index"`
	MaxChanSize  int    `ini:"maxChanSize"`
	GoroutineNum int    `ini:"goroutineNum"`
}
