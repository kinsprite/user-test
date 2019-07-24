package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/json-iterator/go"
	"go.elastic.co/apm/module/apmgin"
)

const prefixV1 = "/api/user/v1"
const prefixV2 = "/api/user/v2"

func init() {
	log.Println("INFO   user.go init")
}

func main() {
	initDB()

	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	engine.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"token": "10020320090",
		})
	})

	v1 := engine.Group(prefixV1)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1.GET("/userInfoBySession", func(c *gin.Context) {
		userID := 10001
		userInfo := getUserInfo(userID)

		if userInfo == nil {
			c.JSON(http.StatusNotFound, userInfo)
			return
		}

		c.JSON(http.StatusOK, userInfo)
	})

	v1.POST("/new", func(c *gin.Context) {
		var input userInfoInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userInfo := &UserInfo{
			Name:  input.Name,
			Email: input.Email,
		}

		createUserInfoToDB(userInfo)

		if userInfo.ID != 0 {
			setUserInfoToCache(userInfo)
		}

		c.JSON(http.StatusOK, userInfo)
	})

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

func getUserInfo(userID int) *UserInfo {
	userInfo := getUserInfoFromCache(userID)

	if userInfo == nil {
		userInfo = getUserInfoFromDB(userID)

		if userInfo != nil {
			setUserInfoToCache(userInfo)
		}
	}

	return userInfo
}
