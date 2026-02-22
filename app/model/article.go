package model

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Title       string `gorm:"type:varchar(255);not null;index:idx_fulltext,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"`
	Body        string `gorm:"type:text;not null;index:idx_fulltext,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"`
	AuthorId    int64  `gorm:"type:int;index"`
	User        User   `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ViewCount   int    `gorm:"default:0"`
	LikeCount   int    `gorm:"default:0"`
	IsAvailable bool   `gorm:"default:true;index"`
}

type Answer struct {
	gorm.Model
	QuestionId  uint     `gorm:"index;not null"`
	Question    Question `gorm:"foreignKey:QuestionId"`
	AuthorId    int64    `gorm:"type:int;index"`
	User        User     `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content     string   `gorm:"type:text;not null;index:,class:FULLTEXT,option:WITH PARSER ngram VISIBLE"`
	LikeCount   int      `gorm:"default:0"`
	IsAvailable bool     `gorm:"default:true;index"`
}

type Comment struct {
	gorm.Model
	AnswerId uint     `gorm:"index;not null"`
	AuthorId int64    `gorm:"type:int;index"`
	User     User     `gorm:"foreignKey:AuthorId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content  string   `gorm:"type:text;not null"`
	ParentId *uint    `gorm:"index"`
	Parent   *Comment `gorm:"foreignKey:ParentId"`
}
