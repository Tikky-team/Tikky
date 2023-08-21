package main

import (
	"Tikky/handler"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	douyin := r.Group("/douyin")

	// feed&publish service
	douyin.GET("/feed", handler.FeedAction)
	publishGroup := douyin.Group("/publish")
	publishGroup.POST("/action/", handler.PublishAction)
	publishGroup.GET("/list", handler.PublishList)

	// user service
	userGroup := douyin.Group("/user")
	userGroup.GET("/", handler.GetUserInfo)

	// auth service
	userGroup.POST("/register/", handler.Register)
	userGroup.POST("/login/", handler.Login)

	// favotite&comment service
	favoriteGroup := douyin.Group("/favotite")
	favoriteGroup.POST("/action/", handler.FavoriteAction)
	favoriteGroup.GET("/list/", handler.FavoriteList)
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action/", handler.CommentAction)
	commentGroup.GET("/list/", handler.CommentList)
}
