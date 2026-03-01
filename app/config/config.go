package config

import (
	"errors"
	"my_zhihu_backend/app/log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var l = log.L().With(zap.String("module", "config"))

type ReadConfigFunc func() Config

type Config struct {
	App     AppConfig         `mapstructure:"APP" yaml:"app"`
	Mysql   MysqlConfig       `mapstructure:"MYSQL" yaml:"mysql"`
	Redis   RedisConfig       `mapstructure:"REDIS" yaml:"redis"`
	Prefix  RedisPrefixConfig `mapstructure:"PREFIX" yaml:"prefix"`
	Service ServiceConfig     `mapstructure:"SERVICE" yaml:"service"`
}

type AppConfig struct {
	ListenAddr string `mapstructure:"LISTEN_ADDR" yaml:"listenAddr"`
}

type RedisConfig struct {
	Addr string `mapstructure:"ADDR" yaml:"addr"`
}
type MysqlConfig struct {
	Host     string `mapstructure:"HOST" yaml:"host"`
	Port     int    `mapstructure:"PORT" yaml:"port"`
	DBName   string `mapstructure:"DB_NAME" yaml:"dbName"`
	User     string `mapstructure:"USER" yaml:"user"`
	Password string `mapstructure:"PASSWORD" yaml:"password"`
}

type ServiceConfig struct {
	RefreshTokenExp time.Duration `mapstructure:"REFRESH_TOKEN_EXP" yaml:"refreshTokenExp"`
	AccessTokenExp  time.Duration `mapstructure:"ACCESS_TOKEN_EXP" yaml:"accessTokenExp"`
	Timeout         time.Duration `mapstructure:"TIMEOUT" yaml:"timeout"`
}

type RedisPrefixConfig struct {
	RefreshToken string `mapstructure:"REFRESH_TOKEN" yaml:"refreshToken"`

	UserInfoPrefix   string `mapstructure:"USERINFO_PREFIX" yaml:"userInfoPrefix"`
	UserSearchPrefix string `mapstructure:"USER_SEARCH_PREFIX" yaml:"userSearchPrefix"`
}

var cfg Config
var rmu sync.RWMutex

func InitConfig() {
	// 设置默认值
	viper.SetDefault("app.LISTEN_ADDR", ":8080")
	viper.SetDefault("mysql.HOST", "127.0.0.1")
	viper.SetDefault("mysql.PORT", 3306)
	viper.SetDefault("mysql.DB_NAME", "zhihu")
	viper.SetDefault("mysql.USER", "root")
	viper.SetDefault("mysql.PASSWORD", "@@XXIIAA@@")
	viper.SetDefault("redis.ADDR", "127.0.0.1:6379")
	viper.SetDefault("prefix.REFRESH_TOKEN", "refreshToken::")
	viper.SetDefault("prefix.USERINFO_PREFIX", "userInfo::")
	viper.SetDefault("prefix.USER_SEARCH_PREFIX", "userSearch::")
	viper.SetDefault("service.REFRESH_TOKEN_EXP", 7*24*time.Hour)
	viper.SetDefault("service.ACCESS_TOKEN_EXP", 5*time.Minute)
	viper.SetDefault("service.TIMEOUT", 5*time.Second)

	// 设置配置文件查找路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../")

	// 启用环境变量自动映射
	viper.AutomaticEnv()

	// 设置环境变量前缀分隔符映射
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		l.Error("failed to read config", zap.Error(err))
		l.Warn("default value will be used")
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		l.Panic("failed to unmarshal config", zap.Error(err))
	}

	l.Info("config loaded", zap.Any("config", cfg))
}

func C() Config {
	rmu.RLock()
	defer rmu.RUnlock()
	return cfg
}

func ReviseConfig(revise func(oldCfg Config) Config) error {
	rmu.Lock()
	defer rmu.Unlock()
	nCfg := revise(cfg)
	viper.Set("app", nCfg.App)
	viper.Set("mysql", nCfg.Mysql)
	viper.Set("redis", nCfg.Redis)
	viper.Set("prefix", nCfg.Prefix)
	viper.Set("service", nCfg.Service)
	cfg = nCfg

	if err := viper.WriteConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				l.Error("failed to create and write config", zap.Error(err))
				return err
			}
		} else {
			return err
		}
	}

	l.Info("config updated", zap.Any("config", cfg))
	return nil
}
