package main

import (
	"github.com/gin-gonic/gin"
	"github.com/piopiop1178/go_levelup/common"
	"github.com/piopiop1178/go_levelup/routers"
)

func main() {
	r := gin.Default()

	common.RedisInit()
	routers.Init(r)

	r.Run()
}
