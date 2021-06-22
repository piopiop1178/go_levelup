package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler "github.com/piopiop1178/go_levelup/handler/auth"
	"github.com/piopiop1178/go_levelup/middleware"
	"github.com/piopiop1178/go_levelup/models"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(router *gin.Engine) {
	db := models.TempDb{}
	db.Init()

	tokendb := models.Redis{}
	// tokendb := models.TempTokenDb{}
	tokendb.Init()

	tokenhdlr := handler.TokenHandler{}

	//이렇게 초기화하는게 맞는지 모르겠음 / tokenhandler는 method만 쓸건데 넣어줘야하는지????
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

	tm := &middleware.TokenMiddleware{
		TokenDb:   &tokendb,
		TokenHdlr: tokenhdlr,
	}

	//이거 맞는지 모르겠다
	router.LoadHTMLGlob("templates/*")

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"team": "K3"})
	})

	// router.POST("/login", handler.Login)

	router.POST("/login", w.Login)

	router.POST("/logout", tm.TokenAuthMiddleware(), w.Logout)

	router.POST("/token_test", tm.TokenAuthMiddleware(), th.CreateTodo)

	router.POST("/token/refresh", w.TokenRefresh)
}
