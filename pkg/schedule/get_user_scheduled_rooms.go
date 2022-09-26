package schedule

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetUserScheduledRooms(c *gin.Context) {
	user := authorization.ExtractUser(c)

	var schedules []models.Schedule
	if err := h.DB.Where("user_id = ?", user.ID).Preload("Room.Building.Address").Find(&schedules).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.InternalServerError("Error while getting schedules"))
		return
	}

	c.JSON(http.StatusOK, schedules)
}
