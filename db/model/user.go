package model

import "gorm.io/gorm"

// User 用户表 /*
type User struct {
	gorm.Model
	Username        string  `gorm:"not null;unique;size: 32;index"`
	Password        *string `gorm:"not null;size: 32"`
	Avatar          string  `json:"avatar"`                      // 用户头像
	BackgroundImage string  `json:"background_image"`            // 用户个人页顶部大图
	FavoriteCount   int64   `json:"favorite_count"`              // 喜欢数
	FollowCount     int64   `json:"follow_count"`                // 关注总数
	FollowerCount   int64   `json:"follower_count"`              // 粉丝总数
	IsFollow        bool    `json:"is_follow"`                   // true-已关注，false-未关注
	Signature       string  `json:"signature"`                   // 个人简介
	TotalFavorited  int64   `json:"total_favorited"`             // 获赞数量
	WorkCount       int64   `gorm:"default:0" json:"work_count"` // 作品数
	Token           string  `json:"token"`
}
