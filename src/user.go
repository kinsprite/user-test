package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/json-iterator/go"
	"go.elastic.co/apm/module/apmgin"
)

const prefixV1 = "/api/user/v1"
const prefixV2 = "/api/user/v2"

func main() {
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
		c.JSON(http.StatusOK, gin.H{
			"id":   10001,
			"name": "User 1",
		})
	})

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
