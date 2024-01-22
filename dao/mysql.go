package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"mini-gpt/setting"
)

var (
	DB *gorm.DB
)

var logger logrus.Logger

func InitMySQL(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		logger.Error(err)
		return
	}
	return DB.DB().Ping()
}

func Close() {
	err := DB.Close()
	if err != nil {
		logger.Error(err)
		return
	}
}
