package repository

import (
	"fmt"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/log"
	"my_zhihu_backend/app/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDBConn(cfg config.ReadConfigFunc) *gorm.DB {
	c := cfg().Mysql
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName)), &gorm.Config{
		Logger:         log.NewGormZapLogger(log.L().With(zap.String("module", "mysql")), 200*time.Millisecond),
		TranslateError: true,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(new(model.User), new(model.UserFollowers), new(model.Answer), new(model.Question), new(model.Comment)); err != nil {
		panic(err)
	}
}
