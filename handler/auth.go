package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/piopiop1178/go_levelup/models"
)

//for test
var user = models.User{
	ID:       1,
	Username: "kangsan",
	Password: "1234",
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

	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

func CreateToken(userid uint64) (token string, err error) {

	os.Setenv("ACCESS_TOKEN_KEY", "kkkkkkkkk") //따로 설정 필요

	//at == access token
	//claim = jwt 구성요소 -> header, claim, signature
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() //why unix?
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims) //HS256 -> 대칭키방식
	token, err = at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_KEY")))

	if err != nil {
		return "", err
	}
	return token, nil
}
