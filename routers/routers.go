package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piopiop1178/go_levelup/handler"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*")

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"team": "K3"})
	})

	router.POST("/login", handler.Login)
}
