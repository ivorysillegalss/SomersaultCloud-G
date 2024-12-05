package bootstrap

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type IpConfigEnv struct {
	AppEnv          string `mapstructure:"appenv"`
	DiscoveryConfig struct {
		Endpoints   []string      `mapstructure:"endpoints"`
		Timeout     time.Duration `mapstructure:"timeout"`
		Username    string        `mapstructure:"username"`
		Password    string        `mapstructure:"password"`
		ServicePath string        `mapstructure:"service_path"`
	} `mapstructure:"discovery"`
}

func NewEnv() *IpConfigEnv {
	env := IpConfigEnv{}
	viper.SetConfigFile("somersaultcloud.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
