package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piopiop1178/go_levelup/handler"
	"github.com/piopiop1178/go_levelup/models"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(router *gin.Engine) {
	//이렇게 초기화하는게 맞는지 모르겠음 / 왜 db, tokendb는 포인터?(인터페이스 내 메소드 사용 가능해야함 -> 포인터?)
	w := &handler.LoginWorker{
		Db:           &models.TempDb{},
		TokenDb:      &models.TempTokenDb{},
		TokenHandler: handler.TokenHandler{},
	}
	//이거 맞는지 모르겠다

	router.LoadHTMLGlob("templates/*")

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"team": "K3"})
	})

	// router.POST("/login", handler.Login)

	router.POST("/login", w.Login)
}
