package room

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateRoomRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Building    string `json:"building" binding:"required"`
}

func (h handler) CreateRoom(c *gin.Context) {
	user := authorization.ExtractUser(c)

	body := CreateRoomRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]exceptions.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = exceptions.ApiError{Param: fe.Field(), Message: exceptions.MsgForTag(fe)}
			}

			c.JSON(http.StatusBadRequest, exceptions.BadValidation(out))
			return
		}
	}

	var building models.Building
	if err := h.DB.Where("id = ? AND user_id = ?", body.Building, user.ID).First(&building).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Building not found"))
		return
	}

	room := models.Room{
		Name:        body.Name,
		Description: body.Description,
		BuildingID:  body.Building,
	}

	if result := h.DB.Save(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusOK, room)
}
