package util

import (
	"github.com/spf13/viper"
	"log"

	"net"
)

// TODO 修改获取IP的方式，no遍历
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			return ipNet.IP.String()
		}
	}
	return ""
}

// TODO 修改获取port的方式
func GetLocalPort() int32 {
	viper.SetConfigFile("somersaultcloud.yaml")
	viper.SetConfigType("yaml")
	port := viper.Get("port")
	return port.(int32)
}
