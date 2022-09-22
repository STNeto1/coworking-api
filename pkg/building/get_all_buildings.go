package building

import (
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) GetAllBuildings(c *gin.Context) {
	name := c.Query("name")
	description := c.Query("description")

	query := h.DB.Preload("Address")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+strings.ToLower(name)+"%")
	}

	if description != "" {
		query = query.Where("description LIKE ?", "%"+strings.ToLower(description)+"%")
	}

	var buildings []models.Building

	if err := query.Find(&buildings).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.InternalServerError("Error while getting buildings"))
		return
	}

	c.JSON(http.StatusOK, buildings)
}
