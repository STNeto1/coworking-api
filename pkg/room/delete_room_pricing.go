package room

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteRoomPricing(c *gin.Context) {
	id := c.Param("id")
	user := authorization.ExtractUser(c)

	var pricing models.RoomPricing
	if err := h.DB.Where("id = ?", id).Preload("Room.Building").First(&pricing).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	if user.ID != pricing.Room.Building.UserID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Unauthorized("Unauthorized"))
		return
	}

	if result := h.DB.Delete(&pricing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
