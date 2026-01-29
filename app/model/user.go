package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	UUID           string         `gorm:"primary_key;type:char(36)" json:"uuid"`
	Username       string         `gorm:"unique;not null;type:varchar(50);not null" json:"username"`
	HPassword      string         `gorm:"not null;type:varchar(128);not null" json:"-"` // 使用bcrypt来生成哈希值 不需要存储盐值 varchar(128)为将来算法升级准备
	Email          string         `gorm:"unique;not null;type:varchar(100);not null;index" json:"email"`
	Followers      []UserModel    `gorm:"many2many:user_followers;foreignKey:UUID;joinForeignKey:FollowingID;References:UUID;joinReferences:FollowerID" json:"followers"`  // 我的粉丝
	Followings     []UserModel    `gorm:"many2many:user_followers;foreignKey:UUID;joinForeignKey:FollowerID;References:UUID;joinReferences:FollowingID" json:"followings"` // 我的关注
	FollowerCount  int            `gorm:"default:0;not null" json:"follower_count"`
	FollowingCount int            `gorm:"default:0;not null" json:"following_count"`
	Gender         UserGender     `gorm:"default:0;not null" json:"gender"`
	Region         string         `gorm:"not null;type:varchar(50);default:'';not null" json:"region"`
	Other          UserOtherInfo  `gorm:"serializer:json;type:json" json:"other"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// UserFollowers 联合主键保证关注关系不会重复出现 a->b 不允许重复 但a->b b->a可同时存在
type UserFollowers struct {
	FollowingID string    `gorm:"primaryKey;type:char(36)" json:"following_id"`
	FollowerID  string    `gorm:"primaryKey;type:char(36);index" json:"follower_id"` // 给粉丝加索引 防止反向查找过慢
	CreatedAt   time.Time `json:"created_at"`
}

type UserGender int

const (
	UserGenderSecret UserGender = iota
	UserGenderMale
	UserGenderFemale
)

// UserOtherInfo 个性签名 简介 头像等不参与计算比较的杂项
type UserOtherInfo struct {
	Introduction string `json:"introduction"`
	Icon         string `json:"icon"` // 指向图像的URL地址
}
