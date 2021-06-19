package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"

	"github.com/piopiop1178/go_levelup/common"
	"github.com/piopiop1178/go_levelup/models"
)

//for test
var user = models.User{
	ID:       1,
	Username: "kangsan",
	Password: "1234",
}

type TokenInfo struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func Login(c *gin.Context) {
	var u models.User

	//request를 User 구조체에 binding
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid data provided") //badrequest?
		return
	}

	//-------------------------------------
	//db에서 input id, password로 등록된 정보 있나 확인 함수로 교체

	//임시로 위에 설정한 user와 일치하는지 확인
	if u.Username != user.Username || u.Password != user.Password {
		c.JSON(http.StatusUnauthorized, "Login details incorrect")
	}
	//-------------------------------------

	ti, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := SaveTokenInfoToRedis(int64(u.ID), ti)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token":  ti.AccessToken,
		"refresh_token": ti.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func CreateToken(userid uint64) (tokeninfo *TokenInfo, err error) {
	ti := &TokenInfo{}
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

func SaveTokenInfoToRedis(userid int64, ti *TokenInfo) error {
	client := common.GetClient()
	ctx := common.Ctx

	at := time.Unix(ti.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(ti.RtExpires, 0)
	now := time.Now()

	//Atoi -> 문자열을 숫자로, Itoa -> 숫자를 문자열로
	errAt := client.Set(ctx, ti.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAt != nil {
		return errAt
	}

	errRt := client.Set(ctx, ti.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errAt != nil {
		return errRt
	}

	return nil
}
