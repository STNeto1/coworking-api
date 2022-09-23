package room

import (
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ShowRoom(c *gin.Context) {
	id := c.Param("id")

	var room models.Room

	if err := h.DB.Preload("Building.Address").Preload("Schedules", "status = 'open'").First(&room, "id = ?", id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	c.JSON(http.StatusOK, room)
}
