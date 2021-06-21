package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/piopiop1178/go_levelup/models"
)

type LoginWorker struct {
	TokenHandler TokenHandler
	TokenDb      models.TokenDb
	Db           models.DB
	//method에 user 사용하는데 들어가야하는지?????
}

func (w *LoginWorker) Login(c *gin.Context) {
	var u models.User //worker 요소로 담아야하는지??
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid data provided")
		return
	}

	userId, err := w.Db.CheckLoginDetails(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Login details incorrect")
		return
	}
	u.ID = userId

	ti, err := w.TokenHandler.CreateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := w.TokenDb.SaveTokenToDb(uint64(u.ID), ti)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ti.AccessToken,
		"refresh_token": ti.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func (w *LoginWorker) Logout(c *gin.Context) {
	accessUuid, err := w.TokenHandler.ExtractAccessUuid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}
	deleted, delErr := w.TokenDb.DeleteToken(accessUuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}
