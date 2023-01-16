package utils

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

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
