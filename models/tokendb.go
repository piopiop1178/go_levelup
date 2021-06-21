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
	SaveTokenToDb(userId uint64, ti *TokenInfo) error
	CheckAccessTokenValidation(accessUuid string) (userId uint64, err error)
	DeleteToken(Uuid string) (int64, error)
}

type TempTokenDb struct {
}

func (t *TempTokenDb) Init() {
}

func (t *TempTokenDb) SaveTokenToDb(userId uint64, ti *TokenInfo) error {
	//tmp list에 저장?
	return nil
}

func (t *TempTokenDb) CheckAccessTokenValidation(accessUuid string) (userId uint64, err error) {
	return 0, nil
}

func (t *TempTokenDb) DeleteToken(Uuid string) (int64, error) {
	return 1, nil
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

func (r *Redis) SaveTokenToDb(userId uint64, ti *TokenInfo) error {
	ctx := context.Background()

	at := time.Unix(ti.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(ti.RtExpires, 0)
	now := time.Now()

	//Atoi -> 문자열을 숫자로, Itoa -> 숫자를 문자열로
	errAt := r.client.Set(ctx, ti.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAt != nil {
		return errAt
	}

	errRt := r.client.Set(ctx, ti.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errAt != nil {
		return errRt
	}

	return nil
}

func (r *Redis) CheckAccessTokenValidation(accessUuid string) (userId uint64, err error) {
	userIdBeforeParsing, err := r.client.Get(context.Background(), accessUuid).Result()

	if err != nil {
		return 0, nil
	}
	userId, _ = strconv.ParseUint(userIdBeforeParsing, 10, 64)

	return userId, nil
}

func (r *Redis) DeleteToken(Uuid string) (int64, error) {
	deleted, err := r.client.Del(context.Background(), Uuid).Result()

	if err != nil {
		return 0, err
	}

	return deleted, nil
}
