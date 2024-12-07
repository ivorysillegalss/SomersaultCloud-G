package bootstrap

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type ExporterEnv struct {
	AppEnv string `mapstructure:"appenv"`

	Grpc struct {
		Monitor struct {
			Port   int    `mapstructure:"port"`
			Server string `mapstructure:"server"`
		} `mapstructure:"monitor"`
	} `mapstructure:"grpc"`

	DiscoveryConfig struct {
		Endpoints   []string      `mapstructure:"endpoints"`
		Timeout     time.Duration `mapstructure:"timeout"`
		Username    string        `mapstructure:"username"`
		Password    string        `mapstructure:"password"`
		ServicePath string        `mapstructure:"service_path"`
	} `mapstructure:"discovery"`

	BusinessConfig struct {
		Address []struct {
			Name      string `mapstructure:"name"`
			Endpoints string `mapstructure:"endpoints"`
			IP        string `mapstructure:"ip"`
			Port      int32  `mapstructure:"port"`
		} `mapstructure:"address"`
	} `mapstructure:"business_config"`
}

func NewEnv() *ExporterEnv {
	env := ExporterEnv{}
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
