package room

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	user := authorization.ExtractUser(c)

	var room models.Room
	if err := h.DB.Where("id = ?", id).Preload("Building").First(&room).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	if user.ID != room.Building.UserID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Unauthorized("Unauthorized 222"))
		return
	}

	if result := h.DB.Delete(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
