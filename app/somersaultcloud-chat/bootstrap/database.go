package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/mongo"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/mysql"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/redis"
	"context"
	"fmt"
	"log"
	"time"
)

func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m := env.Mongo
	dbHost := m.Host
	dbPort := m.Port
	dbUser := m.User
	dbPass := m.Pass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NewRedisDatabase(env *Env) redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rconf := env.Redis
	dbAddr := rconf.Addr
	dbPassword := rconf.Password

	client, err := redis.NewRedisClient(redis.NewRedisApplication(dbAddr, dbPassword))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NewMysqlDatabase(env *Env) mysql.Client {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m := env.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.DB)

	client, err := mysql.NewMysqlClient(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

func NewDatabases(env *Env) *Databases {
	return &Databases{
		Redis: NewRedisDatabase(env),
		Mysql: NewMysqlDatabase(env),
	}
}
