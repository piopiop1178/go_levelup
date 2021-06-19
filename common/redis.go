package common

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client
var ctx = context.Background() //context가 뭘까???

func RedisInit() {
	os.Setenv("REDIS_DSN", "redis:6379") //docker compose로 연결되어있으면 서비스 이름으로 연결 가능
	dsn := os.Getenv("REDIS_DSN")

	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	//check connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis init!!!")
}
