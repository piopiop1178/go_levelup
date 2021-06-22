package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

//extract token string from request header
func (t *TokenHandler) ExtractTokenString(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

//get token from token string
func (t *TokenHandler) GetTokenFromTokenString(tokenString string, key string) (*jwt.Token, error) {
	//token string secret key로 복호화
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//check signing method and return secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		fmt.Println("gettokenfromstring")
		return nil, err
	}

	return token, nil
}

func (t *TokenHandler) CheckTokenValidation(token *jwt.Token) error {
	//type assertion
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("Token is not valid")
	}
	//token의 claims의 type이 jwt.claims인지 확인??
	return nil
}

//token db에서 조회할 토큰의 userid, uuid 추출
func (t *TokenHandler) ExtractUserIdandUuid(token *jwt.Token) (uint64, string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid { //valid 체크 안해도 될 것 같은디 -> jwt.MapCLamis와 jwt.Claims의 vaild가 다른지?

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return 0, "", errors.New("error with find userid")
		}
		uuid, ok := claims["access_uuid"].(string) //refresh uuid 추출하는 상황이 있을지??

		if !ok {
			return 0, "", errors.New("error with find uuid")
		}

		return userId, uuid, nil
	}
	return 0, "", errors.New("not valid token")
}
