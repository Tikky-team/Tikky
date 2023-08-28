package handler

import (
	"Tikky/db"
	"Tikky/db/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func FeedAction(c *gin.Context) {
	latestTime := c.Query("latest_time")
	var tempTime int64
	now := time.Now().UnixMilli()
	if latestTime == "0" {
		latestTime = strconv.FormatInt(now, 10)
		tempTime, _ = strconv.ParseInt(latestTime, 10, 0)
	} else {
		tempTime, _ = strconv.ParseInt(latestTime, 10, 0)
	}
	var tempVideo []model.Video
	result := db.Db.Table("videos").Where("index >= ?", time.UnixMilli(tempTime)).Order("index desc").Find(&tempVideo)
	if result.RowsAffected > 1 || tempTime == now {
		var tempVideo []model.Video
		Vrow := db.Db.Table("videos").Find(&tempVideo) //获取video表中行数
		var videos []*VideoRsp
		video := make([]model.Video, Vrow.RowsAffected)
		now := time.Now().UnixMilli()
		tempTime, _ = strconv.ParseInt(strconv.FormatInt(now, 10), 10, 0)
		db.Db.Table("videos").Where("index <= ?", time.UnixMilli(tempTime)).Limit(30).Order("index desc").Find(&video)
		for _, v := range video {
			User := GetUser(v.UserId)
			videos = append(videos, &VideoRsp{
				Author:        User,
				CommentCount:  v.CommentCount,
				CoverURL:      v.CoverURL,
				FavoriteCount: v.FavoriteCount,
				ID:            v.ID,
				IsFavorite:    v.IsFavorite,
				PlayURL:       v.PlayURL,
				Title:         v.Title,
			})
		}
		if video == nil || len(video) == 0 {
			message := "Successfully refreshed"
			c.JSON(http.StatusForbidden, FeedRsp{
				StatusCode: 1,
				StatusMsg:  message,
			})
			return
		} //反复刷新返回响应
		var NextTime int64
		NextTime = video[0].CreatedAt.UnixNano() / int64(time.Millisecond)
		message := "Success"
		c.JSON(http.StatusOK, FeedRsp{
			StatusCode: 0,
			StatusMsg:  message,
			NextTime:   &NextTime,
			VideoList:  videos,
		})
	} else {
		message := "fail"
		c.JSON(http.StatusInternalServerError, FeedRsp{
			StatusCode: 1,
			StatusMsg:  message,
			NextTime:   nil,
			VideoList:  nil,
		})
	}
}
func GetUser(userid uint32) (u *UserRsp) { //通过userid得到user
	var tempUser UserRsp
	result := db.Db.Table("users").Where("id = ?", userid).First(&tempUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil
	}
	return &tempUser
}

type UserRsp struct {
	ID              int64  `json:"id"`                          //用户id
	Name            string `json:"name"`                        // 用户名称
	FollowCount     int64  `json:"follow_count"`                // 关注总数
	FollowerCount   int64  `json:"follower_count"`              // 粉丝总数
	IsFollow        bool   `json:"is_follow"`                   // true-已关注，false-未关注
	Avatar          string `json:"avatar"`                      // 用户头像
	BackgroundImage string `json:"background_image"`            // 用户个人页顶部大图
	Signature       string `json:"signature"`                   // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`             // 获赞数量
	WorkCount       int64  `gorm:"default:0" json:"work_count"` // 作品数
	FavoriteCount   int64  `json:"favorite_count"`              // 喜欢数
}
type VideoRsp struct {
	Author        *UserRsp `json:"author"`         // 视频作者信息
	CommentCount  int64    `json:"comment_count"`  // 视频的评论总数
	CoverURL      string   `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"` // 视频的点赞总数
	ID            uint     `json:"id"`             // 视频唯一标识
	IsFavorite    bool     `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string   `json:"play_url"`       // 视频播放地址
	Title         string   `json:"title"`          // 视频标题
}
type FeedRsp struct {
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`  // 返回状态描述
	NextTime   *int64      `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList  []*VideoRsp `json:"video_list"`  // 视频列表
}
