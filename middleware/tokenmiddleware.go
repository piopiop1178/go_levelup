package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler "github.com/piopiop1178/go_levelup/handler/auth"
)

type TokenMiddleware struct {
	TokenHdlr handler.TokenHandler
}

func (tm *TokenMiddleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tm.TokenHdlr.CheckTokenValidation(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next() //Next?? 가 뭔지??
	}
}
