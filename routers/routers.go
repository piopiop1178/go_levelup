package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler "github.com/piopiop1178/go_levelup/handler/auth"
	"github.com/piopiop1178/go_levelup/models"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(router *gin.Engine) {
	db := models.TempDb{}
	db.Init()

	tokendb := models.Redis{}
	tokendb.Init()

	tokenhdlr := handler.TokenHandler{}

	//이렇게 초기화하는게 맞는지 모르겠음 / 왜 db, tokendb는 포인터?(인터페이스 내 메소드 사용 가능해야함 -> 포인터?)
	w := &handler.LoginWorker{
		Db:           &db,
		TokenDb:      &tokendb,
		TokenHandler: tokenhdlr,
	}

	th := &handler.TodoHandler{
		Db:        &db,
		TokenDb:   &tokendb,
		TokenHdlr: tokenhdlr,
	}

	//같은 db 들어가야되는디?
	//이거 맞는지 모르겠다
	router.LoadHTMLGlob("templates/*")

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"team": "K3"})
	})

	// router.POST("/login", handler.Login)

	router.POST("/login", w.Login)

	router.POST("/logout", w.Logout)

	router.POST("/token_test", th.CreateTodo)
}
