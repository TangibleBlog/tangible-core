package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"tangible-core/public/globalstruct"
	"tangible-core/utils"

	"github.com/go-redis/redis/v8"
)

type Struct struct {
	UseRedis bool
	Redis    *redis.Client
	Context  context.Context
}

func Write(cacheConfig Struct, key string, value []byte) {
	if cacheConfig.UseRedis {
		pong, err := cacheConfig.Redis.Ping(cacheConfig.Context).Result()
		if err != nil {
			log.Println("Redis connection failed:", pong, err)
		}
		cacheConfig.Redis.Set(cacheConfig.Context, key, value, 0)
	}

}

func Read(cacheConfig Struct, key string) (bool, []byte) {
	if !cacheConfig.UseRedis {
		return false, nil
	}
	found := false
	pong, err := cacheConfig.Redis.Ping(cacheConfig.Context).Result()
	if err != nil {
		log.Println("Redis connection failed:", pong, err)
	}
	result, err := cacheConfig.Redis.Get(cacheConfig.Context, key).Result()
	found = true
	if err == redis.Nil {
		log.Println("Page not found")
		found = false
	}
	return found, []byte(result)

}

func Flush(cacheConfig Struct) {
	if cacheConfig.UseRedis {
		cacheConfig.Redis.FlushAll(cacheConfig.Context)
		log.Println("Cache refreshed")
	}
}

func LoadCache() Struct {
	var cacheConfig Struct
	cacheConfig.UseRedis = false
	if utils.CheckFileExist("./config/redis.json") {
		log.Println("Use redis as the page cache engine.")
		cacheConfig.UseRedis = true
		var redisConfig globalstruct.RedisConfigStruct
		err := json.Unmarshal(utils.OpenFile("./config/redis.json"), &redisConfig)
		if err != nil {
			log.Fatal(err)
		}
		cacheConfig.Context = context.Background()
		//goland:noinspection GoSnakeCaseUsage
		cacheConfig.Redis = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, strconv.Itoa(redisConfig.Port)),
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		})
		pong, err := cacheConfig.Redis.Ping(cacheConfig.Context).Result()
		if err != nil {
			log.Panic("Redis connection failed:", pong, err)

		}
		log.Println("Redis connection successful:", pong)
	}
	return cacheConfig

}
