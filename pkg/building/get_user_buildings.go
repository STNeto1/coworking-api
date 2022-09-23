package building

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetUserBuildings(c *gin.Context) {
	user := authorization.ExtractUser(c)

	var buildings []models.Building
	if err := h.DB.Where("user_id = ?", user.ID).Preload("Address").Preload("Rooms.Schedules").Find(&buildings).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.InternalServerError("Error while getting buildings"))
		return
	}

	c.JSON(http.StatusOK, buildings)
}
