package handler

import (
	"Tikky/db"
	"Tikky/db/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	jwt.StandardClaims
}

// 定义过期时间
const TokenExpireDuration = time.Hour * 2

// 定义secret
var MySecret = []byte("token的密钥")

// 生成jwt
func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		"none",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "my-project",
		},
	}
	//创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthHandler(c *gin.Context) {
	//得到handler
	var user MyClaims
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}

	if user.Name == "cyl" && user.Password == "123456" {
		//生成token
		tokenString, _ := GenToken(user.Name)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func AuthMiddleware(c *gin.Context) {
	{
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中的auth为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中的auth格式错误",
			})
			//阻止调用后续的函数
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的token",
			})
			c.Abort()
			return
		}
		//将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Name)
		//后续的处理函数可以通过c.Get("username")
		c.Next()
	}

}
func Register(c *gin.Context) {
	{
		//获取参数
		name := c.PostForm("name")
		password := c.PostForm("password")

		//数据验证
		if len(name) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于6位",
			})
			return
		}

		//判断用户名是否存在
		var user model.User
		db.Db.Where(" name= ?", name).First(&user)
		if user.ID != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户已存在",
			})
			return
		}

		//创建用户
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    500,
				"message": "密码加密错误",
			})
			return
		}
		newUser := model.User{
			Username: name,
			Password: (string(hasedPassword)),
		}
		db.Db.Create(&newUser)

		//返回结果
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注册成功",
		})

		return
	}
	return
}

func Login(c *gin.Context) {
	//获取参数
	id := c.PostForm("id")
	password := c.PostForm("password")

	//数据验证
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}

	//判断id是否存在
	var user model.User
	db.Db.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码错误",
		})
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
	})
}
