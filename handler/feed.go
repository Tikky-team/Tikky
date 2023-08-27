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
		var videos []*model.VideoRsp
		video := make([]model.Video, Vrow.RowsAffected)
		now := time.Now().UnixMilli()
		tempTime, _ = strconv.ParseInt(strconv.FormatInt(now, 10), 10, 0)
		db.Db.Table("videos").Where("index <= ?", time.UnixMilli(tempTime)).Limit(30).Order("index desc").Find(&video)
		for _, v := range video {
			User := GetUser(v.UserId)
			videos = append(videos, &model.VideoRsp{
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
			c.JSON(http.StatusForbidden, model.FeedRsp{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  message,
				},
			})
			return
		} //反复刷新返回响应
		var NextTime int64
		NextTime = video[0].CreatedAt.UnixNano() / int64(time.Millisecond)
		message := "Success"
		c.JSON(http.StatusOK, model.FeedRsp{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  message,
			},
			NextTime:  &NextTime,
			VideoList: videos,
		})
	} else {
		message := "fail"
		c.JSON(http.StatusInternalServerError, model.FeedRsp{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  message,
			},
			NextTime:  nil,
			VideoList: nil,
		})
	}
}
func GetUser(userid uint32) (u *model.UserRsp) { //通过userid得到user
	var tempUser model.UserRsp
	result := db.Db.Table("users").Where("id = ?", userid).First(&tempUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil
	}
	return &tempUser
}
