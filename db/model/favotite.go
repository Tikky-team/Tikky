package model

import "gorm.io/gorm"

// Favorite 点赞表 /*
type Favorite struct {
	gorm.Model
	UserId  uint32 `gorm:"not null;uniqueIndex:user_video"`
	VideoId uint32 `gorm:"not null;uniqueIndex:user_video;index:video"`
}
