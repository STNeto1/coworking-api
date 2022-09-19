package user

import (
	"net/http"

	"ginn/pkg/common/exceptions"
	"ginn/pkg/common/models"

	"github.com/gin-gonic/gin"
)

func (h handler) GetUsers(c *gin.Context) {
	var users []models.User

	if result := h.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Error fetching users"))
		return
	}

	c.JSON(http.StatusOK, &users)
}
