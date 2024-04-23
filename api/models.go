package api

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var ModelsMap = new(Config)

type Config struct {
	Models map[string]string
}

func LoadModels() {
	// 读取YAML配置文件
	data, err := ioutil.ReadFile("../conf/models.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 解析YAML配置到Config结构体
	err = yaml.Unmarshal(data, &ModelsMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//// 打印配置信息，验证是否成功读取
	//for serverName, serverConfig := range config.Models {
	//	fmt.Printf("%s runs at %s:%d\n", serverName, serverConfig.Usage, serverConfig.Name)
	//}
}
