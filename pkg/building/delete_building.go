package building

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) DeleteBuilding(c *gin.Context) {
	user := authorization.ExtractUser(c)
	id := c.Param("id")

	var building models.Building
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&building).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Building not found"))
		return
	}

	if result := h.DB.Delete(&building); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
