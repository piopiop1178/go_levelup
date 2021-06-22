package handler

import (
	"fmt"
	"net/http"
	"os"

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
	tokenString := w.TokenHandler.ExtractTokenString(c.Request)
	if tokenString == "" {
		fmt.Println("string")
		c.JSON(http.StatusBadRequest, "cannot get Token")
		return
	}

	accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
	accessToken, err := w.TokenHandler.GetTokenFromTokenString(tokenString, accessTokenKey)
	if err != nil {
		fmt.Println("token")
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, accessUuid, err := w.TokenHandler.ExtractUserIdandUuid(accessToken)
	if err != nil {
		fmt.Println("getid")
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	deleted, delErr := w.TokenDb.DeleteToken(accessUuid)
	if delErr != nil || deleted == 0 {
		fmt.Println("delete")
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}

func (w *LoginWorker) TokenRefresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBind(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//get refreshtoken
	refreshTokenKey := os.Getenv("REFRESH_TOKEN_KEY")
	token, err := w.TokenHandler.GetTokenFromTokenString(refreshToken, refreshTokenKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Refresh Token was expired")
		return
	}

	//check refreshtoken validation
	if err := w.TokenHandler.CheckTokenValidation(token); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//get userid, uuid from refreshtoken
	userId, uuid, err := w.TokenHandler.ExtractUserIdandUuid(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Error occurred")
		return
	}

	//기존의 리프레시 토큰 삭제
	deleted, delErr := w.TokenDb.DeleteToken(uuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	ti, err := w.TokenHandler.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	saveErr := w.TokenDb.SaveTokenToDb(userId, ti)
	if saveErr != nil {
		c.JSON(http.StatusForbidden, saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ti.AccessToken,
		"refresh_token": ti.RefreshToken,
	}
	c.JSON(http.StatusCreated, tokens)
}
