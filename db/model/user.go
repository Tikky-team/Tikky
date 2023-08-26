package model

import "gorm.io/gorm"

// User 用户表 /*
type User struct {
	gorm.Model
	Username string `gorm:"not null;unique;size: 32;index"`
	Password string `gorm:"not null;size: 32"`
}
