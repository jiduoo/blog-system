package config

import (
	"backend/global"
	"log"

	"github.com/go-redis/redis"
)

func initRedis(){

	addr := AppConfig.Redis.Addr
	db := AppConfig.Redis.DB
	password := AppConfig.Redis.Password 

	RedisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		DB: db,
		Password: password,
	})

	_, err := RedisClient.Ping().Result()

	if err !=nil{
		log.Printf("Failed to connect to Redis, got error: %v", err)
		return
	}

	global.RedisDB = RedisClient
	log.Println("Redis initialized successfully")
}