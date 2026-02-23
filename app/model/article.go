package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID          int64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Title       string         `gorm:"type:varchar(255);not null;index:idx_fulltext,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"` // title 和 body 联合索引
	Content     string         `gorm:"type:text;not null;index:idx_fulltext,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"`
	AuthorId    int64          `gorm:"type:int;index"`
	User        User           `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsAvailable bool           `gorm:"default:true;index"`
}

type Answer struct {
	ID          int64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	QuestionId  int64          `gorm:"index;not null"`
	AuthorId    int64          `gorm:"type:int;index"`
	User        User           `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content     string         `gorm:"type:text;not null;index:,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"`
	LikeCount   int            `gorm:"default:0"`
	IsAvailable bool           `gorm:"default:true;index"`
}

type Comment struct {
	ID          int64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	AnswerId    int64          `gorm:"index;not null"`
	AuthorId    int64          `gorm:"type:int;index"`
	User        User           `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content     string         `gorm:"type:text;not null"`
	ParentId    *int64         `gorm:"index"`
	Parent      *Comment       `gorm:"foreignKey:ParentId"`
	LikeCount   int            `gorm:"default:0"`
	IsAvailable bool           `gorm:"default:true;index"`
}
