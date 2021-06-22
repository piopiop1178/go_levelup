package handler

import (
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
	var u models.User

	//request body의 내용 user 구조체로 변환
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid data provided")
		return
	}

	//db에서 올바른 id, password인지 확인
	userId, err := w.Db.CheckLoginDetails(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Login details incorrect")
		return
	}
	u.ID = userId

	//access token, refresh token 사용
	ti, err := w.TokenHandler.CreateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//생성된 토큰 token db에 저장(redis)
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
	//request header에서 access token 추출
	tokenString := w.TokenHandler.ExtractTokenString(c.Request)
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, "cannot get Token")
		return
	}

	//token string token으로 변환(복호화)
	accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
	accessToken, err := w.TokenHandler.GetTokenFromTokenString(tokenString, accessTokenKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//tokendb에서 token의 key인 uuid 추출
	_, accessUuid, err := w.TokenHandler.ExtractUserIdandUuid(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	//tokendb에서 access token 삭제
	atDeleted, atDelErr := w.TokenDb.DeleteToken(accessUuid)
	if atDelErr != nil || atDeleted == 0 {
		c.JSON(http.StatusUnauthorized, "Unauthrized")
		return
	}

	// -------------------------- logout 할 때 refresh token도 삭제 ------------------------------
	// ------------ refresh token header? body? 어디서 받는지 확인중 ----------------------

	// refreshTokenKey := os.Getenv("REFRESH_TOKEN_KEY")
	// refreshToken, err := w.TokenHandler.GetTokenFromTokenString(tokenString, refreshTokenKey)
	// if err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, err.Error())
	// 	return
	// }

	// _, refreshUuid, err := w.TokenHandler.ExtractUserIdandUuid(refreshToken)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, "Unauthrized")
	// 	return
	// }

	// rtDeleted, rtDelErr := w.TokenDb.DeleteToken(refreshUuid)
	// if rtDelErr != nil || rtDeleted == 0 {
	// 	c.JSON(http.StatusUnauthorized, "Unauthrized")
	// 	return
	// }
	// -------------------------- logout 할 때 refresh token도 삭제 ------------------------------

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
