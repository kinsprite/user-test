package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"go.elastic.co/apm/module/apmgin"
)

const prefixV1 = "/api/user/v1"
const prefixV2 = "/api/user/v2"

var redisServerURL = "redis-cache:6379"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type redisPoolItem struct {
	client *redis.Client
}

var redisPool = sync.Pool{
	New: newRedisConnection,
}

func init() {
	url := os.Getenv("REDIS_SERVER_URL")

	if url != "" {
		redisServerURL = url
	}
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
		userInfo := getUserInfoFromCache(userID)

		if userInfo == nil {
			userInfo = &UserInfo{
				ID:   userID,
				Name: "User 1",
			}

			setUserInfoToCache(userInfo)
		}

		c.JSON(http.StatusOK, userInfo)
	})

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

func newRedisConnection() interface{} {
	client := redis.NewClient(&redis.Options{
		Addr:     redisServerURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &redisPoolItem{
		client: client,
	}
}

func getUserCacheKey(userID int) string {
	return fmt.Sprintf("user_%d", userID)
}

func getUserInfoFromCache(userID int) *UserInfo {
	poolItem := redisPool.Get().(*redisPoolItem)
	defer redisPool.Put(poolItem)

	client := poolItem.client

	if client == nil {
		return nil
	}

	key := getUserCacheKey(userID)
	val, err := client.Get(key).Result()

	if err != nil {
		return nil
	}

	log.Printf("INFO    redis got: %s, %s\n", key, val)

	var userInfo UserInfo
	err = json.Unmarshal([]byte(val), &userInfo)

	if err != nil {
		return nil
	}

	return &userInfo
}

func setUserInfoToCache(userInfo *UserInfo) {
	poolItem := redisPool.Get().(*redisPoolItem)
	defer redisPool.Put(poolItem)

	client := poolItem.client

	if client == nil {
		return
	}

	key := getUserCacheKey(userInfo.ID)
	value, err := json.Marshal(userInfo)

	if err != nil {
		return
	}

	err = client.Set(key, string(value), 0).Err()

	if err != nil {
		log.Printf("ERROR    Fail to store user %d info to redis\n", userInfo.ID)
		return
	}

	log.Printf("INFO    Saving user %d info to redis OK\n", userInfo.ID)
}
