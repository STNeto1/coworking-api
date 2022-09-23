package room

import (
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) GetAllRooms(c *gin.Context) {
	name := c.Query("name")
	description := c.Query("description")

	query := h.DB

	if name != "" {
		query = query.Where("name LIKE ?", "%"+strings.ToLower(name)+"%")
	}

	if description != "" {
		query = query.Where("description LIKE ?", "%"+strings.ToLower(description)+"%")
	}

	var rooms []models.Room

	if err := query.Preload("Building").Find(&rooms).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.InternalServerError("Error while getting rooms"))
		return
	}

	c.JSON(http.StatusOK, rooms)
}
