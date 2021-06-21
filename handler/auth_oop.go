package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"

	"github.com/piopiop1178/go_levelup/models"
)

type LoginWorker struct {
	TokenHandler TokenHandler
	TokenDb      models.TokenDb
	Db           models.DB
	//method에 user 사용하는데 들어가야하는지?????
}

type TokenHandler struct {
	//TokenInfo가 안에 들어가는지??
}

func (w *LoginWorker) Login(c *gin.Context) {
	var u models.User //worker 요소로 담아야하는지??
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid data provided") //badrequest?
		return
	}

	if w.Db.CheckLoginDetails(u.Username, u.Password) == false {
		c.JSON(http.StatusUnauthorized, "Login details incorrect")
		return
	}

	ti, err := w.TokenHandler.CreateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := w.TokenDb.SaveTokenToDb(int64(u.ID), ti)
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

func (t *TokenHandler) CreateToken(userid uint64) (tokeninfo *models.TokenInfo, err error) {
	ti := &models.TokenInfo{}
	ti.AtExpires = time.Now().Add(time.Minute * 15).Unix() //access token 유효기간 15분 //why unix?
	ti.AccessUuid = uuid.NewV4().String()
	ti.RtExpires = time.Now().Add(time.Hour * 24).Unix() //refresh token 유효기간 하루
	ti.RefreshUuid = uuid.NewV4().String()

	//at == access token
	//claim = jwt 구성요소 -> header, claim, signature
	//create access token
	os.Setenv("ACCESS_TOKEN_KEY", "kkkkkkkkk") //따로 설정 필요

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = ti.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = ti.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims) //HS256 -> 대칭키방식
	ti.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_KEY")))

	if err != nil {
		return nil, err
	}

	//create refresh token
	os.Setenv("REFRESH_TOKEN_KEY", "sssssssss") //따로 설정 필요
	rtClaims := jwt.MapClaims{}
	rtClaims["access_uuid"] = ti.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = ti.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims) //HS256 -> 대칭키방식
	ti.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN_KEY")))

	if err != nil {
		return nil, err
	}

	return ti, nil
}
