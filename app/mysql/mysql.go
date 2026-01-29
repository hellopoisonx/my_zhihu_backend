package mysql

import (
	"app/config"
	"app/log"
	"app/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const ErrCodeMysql = 101

type Mysql struct {
	db  *gorm.DB
	log *zap.Logger
	cfg *config.Config
}

func NewMysql(l *zap.Logger, cfg *config.Config) (*Mysql, error) {
	db, err := gorm.Open(mysql.Open(cfg.MysqlConfig.Dsn), &gorm.Config{
		Logger:         log.NewGormZapLogger(l, cfg.MysqlConfig.SlowThreshold),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.UserModel{})
	if err != nil {
		return nil, err
	}
	return &Mysql{
		db:  db,
		log: l,
		cfg: cfg,
	}, nil
}
