package util

import (
	"crypto/sha512"
	"my_zhihu_backend/app/log"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var l = log.L().With(zap.String("module", "util"))

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1) // 目前只有一个机器节点 没有集群
	if err != nil {
		l.Panic("failed to init snowflake node", zap.Error(err))
	}
}

type Util struct{}

func (_ *Util) GenerateUUID() string {
	uid := uuid.New() // TODO: 生成uuid可能会panic?
	return uid.String()
}

func (_ *Util) EncryptPassword(password string) ([]byte, error) {
	sum := sha512.Sum512([]byte(password))
	return bcrypt.GenerateFromPassword(sum[:], bcrypt.MinCost)
}

func (_ *Util) ValidatePassword(hPassword, password string) bool {
	sum := sha512.Sum512([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hPassword), sum[:])
	if err != nil {
		l.Warn("password unmatch", zap.Error(err))
		return false
	}
	return true
}

func (_ *Util) GenerateSnowflakeID() int64 {
	return node.Generate().Int64()
}

func Ptr[T any](v T) *T {
	return &v
}
