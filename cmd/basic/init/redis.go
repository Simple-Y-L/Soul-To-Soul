package inits

import (
	"fmt"
	"kratosItems/cmd/basic/config"
	"log"

	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	addr := fmt.Sprintf("%s:%d", config.Configs.Redis.Host, config.Configs.Redis.Port)
	config.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,                          // use default Addr
		Password: config.Configs.Redis.Password, // no password set
		DB:       config.Configs.Redis.DB,       // use default DB
	})

	_, err := config.RDB.Ping(config.Ctx).Result()
	if err != nil {
		panic("redis inits failed:" + err.Error())
	}
	log.Println("redis inits success")
}
