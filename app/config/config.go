package config

import (
	"app/log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)
import _ "app/log"

type Config struct {
	AppConfig
	MysqlConfig
	RedisConfig
}

type AppConfig struct {
	ListenAddress string
}

type MysqlConfig struct {
	Dsn           string
	SlowThreshold time.Duration //seconds
}

type RedisConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

var cfg *Config

func init() {
	l := log.L().With(zap.String("mod", "config"))
	defer func(l *zap.Logger) {
		err := l.Sync()
		if err != nil {
			// TODO
			panic(err)
		}
	}(l)
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		l.Error("read config failed", zap.Error(err))
		l.Warn("use default config")
	}

	v.SetDefault("app.listenAddress", "localhost:8080")
	v.SetDefault("mysql.dsn", "root:@@XXIIAA1122@tcp(127.0.0.1:3306)/zhihu?charset=utf8&parseTime=True&loc=Local")
	v.SetDefault("mysql.slowThreshold", 200*time.Millisecond)
	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.user", "root")

	if err := v.Unmarshal(cfg); err != nil {
		// TODO
		l.Panic("failed to unmarshal config", zap.Error(err))
	}

	l.Info("config inited")
}

func C() *Config {
	if cfg == nil {
		// TODO
		panic("config is nil")
	}
	return cfg
}
