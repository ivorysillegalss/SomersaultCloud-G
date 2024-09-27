package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`

	MongoHost string `mapstructure:"MONGO_DB_HOST"`
	MongoPort string `mapstructure:"MONGO_DB_PORT"`
	MongoUser string `mapstructure:"MONGO_DB_USER"`
	MongoPass string `mapstructure:"MONGO_DB_PASS"`
	MongoName string `mapstructure:"MONGO_DB_NAME"`

	RedisAddr     string `mapstructure:"REDIS_DB_ADDR"`
	RedisPassword string `mapstructure:"REDIS_DB_PASSWORD"`

	MysqlUser     string `mapstructure:"MYSQL_DB_USER"`
	MysqlPassword string `mapstructure:"MYSQL_DB_PASSWORD"`
	MysqlHost     string `mapstructure:"MYSQL_DB_HOST"`
	MysqlPort     int    `mapstructure:"MYSQL_DB_PORT"`
	MysqlDB       string `mapstructure:"MYSQL_DB_DB"`

	RabbitmqUser     string `mapstructure:"RABBITMQ_USER"`
	RabbitmqPassword string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitmqAddr     string `mapstructure:"RABBITMQ_ADDR"`
	RabbitmqPort     string `mapstructure:"RABBITMQ_PORT"`

	JwtSecretToken     string `mapstructure:"JWT_SECRET_KEY"`
	ApiOpenaiSecretKey string `mapstructure:"API_OPENAI_SECRET_KEY"`

	Serializer string `mapstructure:"SERIALIZER"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

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
