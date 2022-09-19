package auth

import (
	"coworking/pkg/common/authorization"

	"github.com/gin-gonic/gin"
)

func (h handler) Profile(c *gin.Context) {
	user := authorization.ExtractUser(c)

	c.JSON(200, gin.H{
		"user": user,
	})
}
