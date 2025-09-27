package middleware

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := service.AuthService.GetTokenClaims(c)
		if err != nil {
			resp.Fail(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
