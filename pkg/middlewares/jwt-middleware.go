package middlewares

import (
	"net/http"

	"ginn/pkg/common/authorization"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := authorization.ValidateToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
