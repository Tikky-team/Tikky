package handler

	import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	)

	type Result struct {
		User
		Follow
		Video
	}
	type User struct {
		UserID          uint   `gorm:"primaryKey"` //用户id
		Name            string `gorm:"not null"`   //用户名
		Avatar          string `gorm:"not null"`   //用户头像
		BackgroundImage string `gorm:"not null"`   //用户个人页顶部大图
		Signature       string `gorm:"not null"`   //个人简介
	}
	type Follow struct {
		UserID        uint `gorm:"foreignKey:UserID;AssociationForeignKey:ID"` //外键
		FollowCount   uint //关注总数
		FollowerCount uint //粉丝总数
		IsFollow      bool `gorm:"not null"` //是否关注
	}

	type Video struct {
		UserID         uint   `gorm:"foreignKey:UserID;AssociationForeignKey:ID"` //外键
		TotalFavorited string //获赞数量
		WorkCount      uint   //作品数
		FavoriteCount  uint   //喜欢数
	}

	func GetUserInfo(c *gin.Context) {
		DataSourceName := "root:123456@tcp(localhost:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(DataSourceName), &gorm.Config{})
		if err != nil {
			panic("database connect error!")
		}
		db.AutoMigrate(&User{}, &Follow{}, &Video{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": -1,
				"status_msg":  err.Error(),
			})
			return
		}

		userid := c.Param("userid") // 从请求参数中获取id
		var user User
		var follow Follow
		var video Video
		result := db.First(&user, userid)
		if user.UserID == 0 { //id为0返回错误信息
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": -1,
				"status_msg":  "User not found",
			})
			return
		} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{ //未查询到id返回错误信息
				"status_code": -1,
				"status_msg":  result.Error.Error(),
			})
			return
		} else {
			db.First(&follow, userid)
			db.First(&video, userid)
			var result Result
			result.User = user
			result.Follow = follow
			result.Video = video
			c.JSON(http.StatusOK, gin.H{
				"status_code": 0,
				"status_msg":  "User info retrieved successfully",
				"user":        result,
			})
		}
	}

