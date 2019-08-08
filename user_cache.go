package main

import (
	"fmt"
	"log"
	"os"

	redis "github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var redisServerURL = "redis:6379"
var client *redis.Client

func initCache() {
	url := os.Getenv("REDIS_SERVER_URL")

	if url != "" {
		redisServerURL = url
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisServerURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Println("INFO    Redis client init")
}

func closeCache() {
	client.Close()
}

func getUserCacheKey(userID int) string {
	return fmt.Sprintf("user_%d", userID)
}

func getUserInfoFromCache(userID int) *UserInfo {
	key := getUserCacheKey(userID)
	val, err := client.Get(key).Result()

	if err != nil {
		log.Println("WARN    redis can't get: ", key, ", ", err)
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
	key := getUserCacheKey(userInfo.ID)
	value, err := json.Marshal(userInfo)

	if err != nil {
		log.Println("WARN    redis set failed: ", key, ", ", err)
		return
	}

	err = client.Set(key, string(value), 0).Err()

	if err != nil {
		log.Printf("ERROR    Fail to store user %d info to redis\n", userInfo.ID)
		return
	}

	log.Printf("INFO    Saving user %d info to redis OK\n", userInfo.ID)
}
