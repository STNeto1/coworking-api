package user

import (
	"net/http"

	"coworking/pkg/common/authorization"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteUser(c *gin.Context) {
	user := authorization.ExtractUser(c)

	h.DB.Delete(&user)

	c.JSON(http.StatusNoContent, nil)
}
