package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/piopiop1178/go_levelup/models"
)

type Todo struct {
	Content string `json:"content"`
	UserId  uint64 `json:"user_id"`
}

type TodoHandler struct {
	TokenHdlr TokenHandler
	Db        models.DB
	TokenDb   models.TokenDb
}

func (th *TodoHandler) CreateTodo(c *gin.Context) {
	var td *Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid data provided")
		return
	}

	tokenString := th.TokenHdlr.ExtractTokenString(c.Request)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, "cannot get Token")
		return
	}

	accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
	accessToken, err := th.TokenHdlr.GetTokenFromTokenString(tokenString, accessTokenKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, AtUuid, err := th.TokenHdlr.ExtractUserIdandUuid(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := th.TokenDb.CheckTokenValidation(AtUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.UserId = userId

	//실제 db에 저장하는 과정 일단 생략
	c.JSON(http.StatusCreated, td)
}
