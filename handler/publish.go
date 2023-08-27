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
	"net/http"
	"os"
	"path"
	"strconv"
)

func PublishAction(c *gin.Context) {
	form, _ := c.MultipartForm()
	file := form.File["data"]
	title := form.Value["title"][0]
	var userinfo model.User
	userid, _ := c.Get("username")
	db.Db.Table("users").Where("id = ?", userid).First(&userinfo)
	var data = make([]byte, file[0].Size)
	contenType := http.DetectContentType(data)
	if contenType != "video/mp4" {
		message := "invalid content type"
		c.JSON(http.StatusResetContent, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	openFile, err := file[0].Open()
	if err != nil {
		message := "Open file failed"
		c.JSON(http.StatusUnauthorized, &model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	readSize, _ := openFile.Read(data)
	if readSize != int(file[0].Size) {
		message := "Size not match"
		c.JSON(http.StatusUnauthorized, model.Response{
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
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}

	err = os.WriteFile(filePath, data, os.FileMode(0755))
	if err != nil {
		message := "Failed to save file"
		c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	coverurl := GetCover(filePath)
	playurl := VideoConvert(filePath)
	var Video model.Video
	err = db.Db.AutoMigrate(&model.Video{})
	if err != nil {
		message := "Failed to create a table"
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	Video = model.Video{
		PlayURL:  playurl,
		Title:    title,
		UserId:   uint32(userinfo.ID),
		CoverURL: coverurl,
	}
	db.Db.Save(&Video)
	result2 := db.Db.Table("users").Where("ID = ?", userid).Update("work_count", userinfo.WorkCount+1)
	if result2.Error != nil {
		message := "Failed to Update"
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	message := "Success"
	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  message,
	})
}
func GetCover(filePath string) (coverurl string) {
	inputFile := filePath
	imagename := ksuid.New().String()
	outfile := path.Join("cover", imagename+".png")
	dir := path.Dir(outfile)
	err := os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	err = ffmpeg_go.Input(inputFile).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output(outfile, ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "png"}).
		WithOutput(buf, os.Stdout).Run()
	if err != nil {
		fmt.Println("Failed to capture image: ", err)
		return ""
	}
	host, _ := os.Hostname()
	coverpath := string(host) + "/douyin/publish/list/" + outfile
	img, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println("Failed to capture image: ", err)
		return ""
	}
	err = imaging.Save(img, coverpath)
	if err != nil {
		fmt.Println("Failed to capture image: ", err)
		return ""
	}
	return coverpath
}
func VideoConvert(filepath string) (palyurl string) {
	inputFile := filepath
	filename := ksuid.New().String()
	outfile := path.Join("out", filename+".mp4")
	dir := path.Dir(outfile)
	_ = os.MkdirAll(dir, os.FileMode(0755))
	err := ffmpeg_go.Input(inputFile).Output(outfile, ffmpeg_go.KwArgs{
		"profile:v": "main",
		"movflags":  "+faststart",
		"crf":       26,
	}).OverWriteOutput().Run()
	if err != nil {
		fmt.Println("Transcoding failed", err)
		return ""
	}
	host, _ := os.Hostname()
	videourl := string(host) + "/douyin/publish/action/" + outfile
	return videourl
}
func PublishList(c *gin.Context) {
	_, tokenbool := c.GetQuery("token")
	userId, userIdbool := c.GetQuery("user_id")
	userid, _ := strconv.ParseInt(userId, 10, 64)
	if !tokenbool || !userIdbool {
		message := "no passed data"
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	var userinfo = model.User{}
	err := db.Db.Table("users").Where("id = ?", userid).First(&userinfo)
	if err.Error != nil {
		message := "no pass data"
		c.JSON(http.StatusInternalServerError, model.Response{
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
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  message,
		})
		return
	}
	var videos []*model.VideoRsp
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
	Smessage := "Success"
	c.JSON(http.StatusOK, model.ListRsp{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  Smessage,
		},
		VideoList: videos,
	})
}
