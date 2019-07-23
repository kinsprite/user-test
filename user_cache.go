package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	redis "github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

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
