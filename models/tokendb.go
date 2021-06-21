package models

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenDb interface {
	Init()
	SaveTokenToDb(userid int64, ti *TokenInfo) error
}

type TempTokenDb struct {
}

func (t *TempTokenDb) Init() {
}

func (t *TempTokenDb) SaveTokenToDb(userid int64, ti *TokenInfo) error {
	//tmp list에 저장?
	return nil
}

type Redis struct {
	client *redis.Client
}

func (r *Redis) Init() {
	var Ctx = context.Background()

	os.Setenv("REDIS_DSN", "redis:6379") //docker compose로 연결되어있으면 서비스 이름으로 연결 가능
	dsn := os.Getenv("REDIS_DSN")

	r.client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	//check connection
	_, err := r.client.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}

func (r *Redis) SaveTokenToDb(userid int64, ti *TokenInfo) error {
	ctx := context.Background()

	at := time.Unix(ti.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(ti.RtExpires, 0)
	now := time.Now()

	//Atoi -> 문자열을 숫자로, Itoa -> 숫자를 문자열로
	errAt := r.client.Set(ctx, ti.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAt != nil {
		return errAt
	}

	errRt := r.client.Set(ctx, ti.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errAt != nil {
		return errRt
	}

	return nil
}
