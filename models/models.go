package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/lzyzsd/go-gin-blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int   `gorm:"primary_key" json:"id"`
	CreatedOn  int64 `json:"created_on"`
	ModifiedOn int64 `json:"modified_on"`
}

func Setup() {
	var (
		err error
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,                 // 慢 SQL 阈值
			LogLevel:      setting.AppSetting.LogLevel, // Log level
			Colorful:      false,                       // 禁用彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                                // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: newLogger,
	})

	if err != nil {
		log.Println(err)
	}
}
