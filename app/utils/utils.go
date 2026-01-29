package utils

import (
	"crypto/sha512"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrCodeEncryptionFailed = 102

type Utils struct{}

func (u *Utils) GenerateUUID() string {
	uid := uuid.New() // TODO: 生成uuid可能会panic?
	return uid.String()
}

func (u *Utils) EncryptPassword(password string) ([]byte, error) {
	sum := sha512.Sum512([]byte(password))
	return bcrypt.GenerateFromPassword(sum[:], bcrypt.MinCost)
}

func (u *Utils) ValidatePassword(hPassword, password string) bool {
	sum := sha512.Sum512([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hPassword), sum[:])
	return err == nil
}

func Ptr[T any](v T) *T {
	return &v
}
