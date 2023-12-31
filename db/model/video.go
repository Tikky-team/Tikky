package model

import (
	"gorm.io/gorm"
	"time"
)

// Video 视频表 /*
type Video struct {
	gorm.Model
	UserId    uint32    `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title"`
	FileName  string    `json:"play_name"`
	CoverName string    `json:"cover_name"`
	CreatedAt time.Time `gorm:"index"`
}
