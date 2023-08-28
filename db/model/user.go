package model

import "gorm.io/gorm"

// User 用户表 /*
type User struct {
	gorm.Model
	WorkCount int64 `gorm:"default:0" json:"work_count"` // 作品数
}
