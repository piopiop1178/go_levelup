package main

import (
	"github.com/gin-gonic/gin"
	"github.com/piopiop1178/go_levelup/routers"
)

func main() {
	r := gin.Default()

	routers.Init(r)

	r.Run()
}
