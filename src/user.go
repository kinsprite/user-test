package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func main() {
	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	engine.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"token": "10020320090",
		})
	})

	engine.GET("/userBySession", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":   10001,
			"name": "User 1",
		})
	})

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
