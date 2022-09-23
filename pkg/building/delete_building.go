package building

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeleteBuilding(c *gin.Context) {
	user := authorization.ExtractUser(c)
	id := c.Param("id")

	var building models.Building
	if err := h.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&building).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Building not found"))
		return
	}

	var rooms []models.Room
	if err := h.DB.Where("building_id = ?", building.ID).Find(&rooms).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.InternalServerError("Error while getting rooms"))
		return
	}

	h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&building).Error; err != nil {
			return err
		}

		for i := 0; i < len(rooms); i++ {
			if err := tx.Delete(&rooms[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if result := h.DB.Delete(&building); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
