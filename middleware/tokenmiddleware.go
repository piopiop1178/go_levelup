package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	handler "github.com/piopiop1178/go_levelup/handler/auth"
	"github.com/piopiop1178/go_levelup/models"
)

type TokenMiddleware struct {
	TokenHdlr handler.TokenHandler
	TokenDb   models.TokenDb
}

func (tm *TokenMiddleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStirng := tm.TokenHdlr.ExtractTokenString(c.Request)
		if tokenStirng == "" {
			c.JSON(http.StatusUnauthorized, "Unauthorized access")
			c.Abort()
			return
		}

		accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
		token, err := tm.TokenHdlr.GetTokenFromTokenString(tokenStirng, accessTokenKey)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			c.Abort()
			return
		}

		if err := tm.TokenHdlr.CheckTokenValidation(token); err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			c.Abort()
			return
		}

		_, uuid, err := tm.TokenHdlr.ExtractUserIdandUuid(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		_, dbErr := tm.TokenDb.CheckAccessTokenValidation(uuid)
		if dbErr != nil {
			c.JSON(http.StatusUnauthorized, "로그아웃 된 토큰")
			c.Abort()
			return
		}
		c.Next() //Next?? 가 뭔지??
	}
}
