package infrastructure

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitRedis() {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	RedisClient = redis.NewClient(&redis.Options{Addr: host + ":" + port})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis: ", err)
	}
}
