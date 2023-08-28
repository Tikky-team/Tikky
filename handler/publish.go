package handler

import (
	"Tikky/db"
	"Tikky/db/model"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func PublishAction(c *gin.Context) {
	form, _ := c.MultipartForm()
	file := form.File["data"]
	title := form.Value["title"][0]
	var userinfo model.User
	userid, _ := c.Get("userid")
	db.Db.Table("users").Where("id = ?", userid).First(&userinfo)
	var data = make([]byte, file[0].Size)
	contenType := http.DetectContentType(data)
	if contenType != "video/mp4" {
		message := "invalid content type"
		c.JSON(http.StatusResetContent, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	openFile, err := file[0].Open()
	if err != nil {
		message := "Open file failed"
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	readSize, _ := openFile.Read(data)
	if readSize != int(file[0].Size) {
		message := "Size not match"
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	filename := ksuid.New().String()
	filePath := path.Join("video", filename+".mp4")
	dir := path.Dir(filePath)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		message := "Failed to create directory"
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}

	err = os.WriteFile(filePath, data, os.FileMode(0755))
	if err != nil {
		message := "Failed to save file"
		c.JSON(http.StatusUnauthorized, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	host := c.Request.Host
	coverurl, _ := GetCover(filePath, host)
	playurl, _ := VideoConvert(filePath, host)
	var Video model.Video
	Video = model.Video{
		PlayURL:  playurl,
		Title:    title,
		UserId:   uint32(userinfo.ID),
		CoverURL: coverurl,
	}
	db.Db.Save(&Video)
	message := "Success"
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  message,
	})
}
func GetCover(filePath string, host string) (coverurl string, err error) {
	const (
		imageFormat = "png"
		frameNumber = 1
	)
	inputFile := filePath
	imagename := ksuid.New().String()
	outfile := filepath.Join("/douyin/publish/list/", imagename+"."+imageFormat)
	dir := filepath.Dir(outfile)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		log.Println("Failed to create directory:", err)
		return
	}
	buf := bytes.NewBuffer(nil)
	err = ffmpeg_go.Input(inputFile).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNumber)}).
		Output(outfile, ffmpeg_go.KwArgs{"vframes": frameNumber, "format": "image2", "vcodec": "png"}).
		WithOutput(buf, os.Stdout).Run()
	if err != nil {
		log.Println("Failed to capture image: ", err)
		return
	}
	coverpath := filepath.Join(host, outfile)
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Println("Failed to decode image: ", err)
		return
	}
	err = imaging.Save(img, coverpath)
	if err != nil {
		log.Println("Failed to save image: ", err)
		return
	}
	return coverpath, nil
}
func VideoConvert(filePath string, host string) (palyurl string, err error) {
	const videoFormat = "mp4"
	inputFile := filePath
	filename := ksuid.New().String()
	outfile := filepath.Join("/douyin/publish/action/", filename+"."+videoFormat)
	dir := filepath.Dir(outfile)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		log.Println("Failed to create directory:", err)
		return
	}
	err = ffmpeg_go.Input(inputFile).Output(outfile, ffmpeg_go.KwArgs{
		"profile:v": "main",
		"movflags":  "+faststart",
		"crf":       26,
	}).OverWriteOutput().Run()
	if err != nil {
		log.Println("Transcoding failed", err)
		return
	}
	videourl := filepath.Join(host, outfile)
	return videourl, nil
}
func PublishList(c *gin.Context) {
	userId, userIdbool := c.GetQuery("user_id")
	userid, _ := strconv.ParseInt(userId, 10, 64)
	if !userIdbool {
		message := "no passed data"
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	var userinfo = model.User{}
	err := db.Db.Table("users").Where("id = ?", userid).First(&userinfo)
	if err.Error != nil {
		message := "no pass data"
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	count, _ := strconv.Atoi(strconv.FormatInt(userinfo.WorkCount, 10))
	var video []model.Video
	err = db.Db.Table("videos").Limit(count).Where("user_id = ?", userinfo.ID).Find(&video)
	if err.Error != nil {
		message := "Failed to search Video list"
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	var videos []*VideoRsp
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
	Smessage := "Success"
	c.JSON(http.StatusOK, ListRsp{
		StatusCode: 0,
		StatusMsg:  Smessage,
		VideoList:  videos,
	})
}

type Response struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}
type ListRsp struct {
	StatusCode int64       `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`  // 返回状态描述
	VideoList  []*VideoRsp `json:"video_list"`  // 用户发布的视频列表
}
