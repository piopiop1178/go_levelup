package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"

	"github.com/piopiop1178/go_levelup/models"
)

type TokenHandler struct {
	//TokenInfo가 안에 들어가는지??
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

//extract token from request header
func (t *TokenHandler) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func (t *TokenHandler) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := t.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TokenHandler) CheckTokenExpiration(r *http.Request) error {
	token, err := t.VerifyToken(r)

	if err != nil {
		return err
	}
	//token의 claims의 type이 jwt.claims인지 확인??
	fmt.Println(token.Claims.(jwt.Claims))
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	//token의 claims의 type이 jwt.claims인지 확인??
	return nil
}

//token db에서 조회할 토큰 accessuuid 추출
func (t *TokenHandler) ExtractAccessUuid(r *http.Request) (string, error) {
	token, err := t.VerifyToken(r)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		atUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return "", errors.New("error with find access uuid")
		}

		return atUuid, nil
	}
	return "", errors.New("not valid token")
}
