package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Client interface {
	Ping() error
	Gorm() *gorm.DB
}

type mysqlClient struct {
	gorm *gorm.DB
}

func (mc *mysqlClient) Gorm() *gorm.DB {
	return mc.gorm
}

func (mc *mysqlClient) Ping() error {
	return mc.gorm.DB().Ping()
}

func NewMysqlClient(dsn string) (Client, error) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil || db == nil {
		log.Fatal(err)
	}
	return &mysqlClient{gorm: db}, nil
}
