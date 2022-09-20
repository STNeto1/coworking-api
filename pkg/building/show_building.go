package building

import (
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ShowBuilding(c *gin.Context) {
	id := c.Param("id")

	var building models.Building

	if err := h.DB.Preload("Address").First(&building, "id = ?", id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Building not found"))
		return
	}

	c.JSON(http.StatusOK, building)
}
