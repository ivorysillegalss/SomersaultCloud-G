package bootstrap

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	AppEnv         string    `mapstructure:"app_env" yaml:"app_env"`
	ServerAddress  string    `mapstructure:"server_address" yaml:"server_address"`
	ContextTimeout int       `mapstructure:"context_timeout" yaml:"context_timeout"`
	Port           int       `mapstructure:"port" yaml:"port"`
	Net            NetConfig `mapstructure:"net" yaml:"net"`

	//可观测性
	Statistics struct {
	} `mapstructure:"statistics" yaml:"statistics"`

	Mongo struct {
		Host string `mapstructure:"host" yaml:"host"`
		Port int    `mapstructure:"port" yaml:"port"`
		User string `mapstructure:"user" yaml:"user"`
		Pass string `mapstructure:"pass" yaml:"pass"`
		Name string `mapstructure:"name" yaml:"name"`
	} `mapstructure:"mongo" yaml:"mongo"`

	Redis struct {
		Addr     string `mapstructure:"addr" yaml:"addr"`
		Password string `mapstructure:"password" yaml:"password"`
	} `mapstructure:"redis" yaml:"redis"`

	Mysql struct {
		User     string `mapstructure:"user" yaml:"user"`
		Password string `mapstructure:"password" yaml:"password"`
		Host     string `mapstructure:"host" yaml:"host"`
		Port     int    `mapstructure:"port" yaml:"port"`
		DB       string `mapstructure:"db" yaml:"db"`
	} `mapstructure:"mysql" yaml:"mysql"`

	Rabbitmq struct {
		User     string `mapstructure:"user" yaml:"user"`
		Password string `mapstructure:"password" yaml:"password"`
		Addr     string `mapstructure:"addr" yaml:"addr"`
		Port     int    `mapstructure:"port" yaml:"port"`
	} `mapstructure:"rabbitmq" yaml:"rabbitmq"`

	Grpc struct {
		Monitor struct {
			Port int `mapstructure:"port" yaml:"port"`
		} `mapstructure:"monitor" yaml:"monitor"`
	} `mapstructure:"grpc" yaml:"grpc"`

	Prometheus struct {
		ServerAddress string `mapstructure:"server_address" yaml:"server_address"`
	} `mapstructure:"prometheus" yaml:"prometheus"`

	JwtSecretToken     string `mapstructure:"jwt_secret_token" yaml:"jwt_secret_token"`
	ApiOpenaiSecretKey string `mapstructure:"api_openai_secret_key" yaml:"api_openai_secret_key"`
	DeepSeekSecretKey  string `mapstructure:"api_deepseek_secret_key" yaml:"api_deepseek_secret_key"`

	Serializer string `mapstructure:"serializer" yaml:"serializer"`
}

// NetConfig 网络配置（限流  跨域等配置）
type NetConfig struct {
	RateLimit RateLimit `mapstructure:"ratelimit" yaml:"ratelimit"`
	Cors      struct {
		Default             bool `mapstructure:"default" yaml:"default"`
		AllowAllOrigin      bool `mapstructure:"allow_all_origin" yaml:"allow_all_origin"`
		AllowAllCredentials bool `mapstructure:"allow_all_credentials" yaml:"allow_all_credentials"`
		//这里可能会有bug
		Headers []string `mapstructure:"headers" yaml:"headers"`
		Methods []string `mapstructure:"methods" yaml:"methods"`
	} `mapstructure:"cors" yaml:"cors"`
}

// RateLimit 限流
// TODO 设置一个Buck和其他限流类型的父类
type RateLimit struct {
	Buck struct {
		Prefix    string `mapstructure:"prefix" yaml:"prefix"`
		Capacity  string `mapstructure:"capacity" yaml:"capacity"`
		Rate      string `mapstructure:"rate" yaml:"rate"`
		Requested string `mapstructure:"requested" yaml:"requested"`
	} `mapstructure:"buck" yaml:"buck"`

	//TODO
	Granularity []string `mapstructure:"granularity" yaml:"granularity"` // 颗粒度 （头 前缀）
}

func NewEnv() *Env {
	env := Env{}
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
