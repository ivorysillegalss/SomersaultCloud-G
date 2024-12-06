package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"app_env" yaml:"app_env"`
	ServerAddress  string `mapstructure:"server_address" yaml:"server_address"`
	ContextTimeout int    `mapstructure:"context_timeout" yaml:"context_timeout"`

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

	JwtSecretToken     string `mapstructure:"jwt_secret_token" yaml:"jwt_secret_token"`
	ApiOpenaiSecretKey string `mapstructure:"api_openai_secret_key" yaml:"api_openai_secret_key"`

	Serializer string `mapstructure:"serializer" yaml:"serializer"`
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
