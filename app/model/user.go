package model

import (
	"time"

	"gorm.io/gorm"
)

type UserId int64 // 使用雪花算法生成 UserId

type User struct {
	Id             UserId         `gorm:"primaryKey;not null;type:int" json:"id"`
	Username       string         `gorm:"index;not null;type:varchar(50)" json:"username"`
	HPassword      string         `gorm:"not null;type:varchar(128);not null" json:"-"` // 使用bcrypt来生成哈希值 不需要存储盐值 varchar(128)为将来算法升级准备
	Email          string         `gorm:"unique;not null;type:varchar(100);index" json:"email"`
	Followers      []User         `gorm:"many2many:user_followers;foreignKey:Id;joinForeignKey:FollowingID;References:Id;joinReferences:FollowerID" json:"followers"`  // 我的粉丝
	Followings     []User         `gorm:"many2many:user_followers;foreignKey:Id;joinForeignKey:FollowerID;References:Id;joinReferences:FollowingID" json:"followings"` // 我的关注
	FollowerCount  int            `gorm:"not null;type:int;default:0" json:"follower_count"`
	FollowingCount int            `gorm:"not null;type:int;default:0" json:"following_count"`
	Gender         *UserGender    `gorm:"default:0;not null" json:"gender"` // 使用指针类型区分null和默认零值
	Region         string         `gorm:"not null;type:varchar(50);default:'';not null" json:"region"`
	Settings       UserSettings   `gorm:"serializer:json;type:json" json:"settings"`
	Other          UserOtherInfo  `gorm:"serializer:json;type:json" json:"other"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// UserFollowers 联合主键保证关注关系不会重复出现 a->b 不允许重复 但a->b b->a可同时存在
type UserFollowers struct {
	FollowingID UserId    `gorm:"primaryKey;type:int;index" json:"following_id"`
	FollowerID  UserId    `gorm:"primaryKey;type:int;index" json:"follower_id"` // 给粉丝加索引 防止反向查找过慢
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

type UserSettings struct {
	HidePrivacy bool `json:"hide_privacy"`
}
