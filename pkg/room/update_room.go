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

type UpdateRoomRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h handler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	user := authorization.ExtractUser(c)

	body := UpdateRoomRequestBody{}

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

	var room models.Room
	if err := h.DB.Where("id = ?", id).Preload("Building").First(&room).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	if user.ID != room.Building.UserID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Unauthorized("Unauthorized 222"))
		return
	}

	room.Name = body.Name
	room.Description = body.Description

	if result := h.DB.Save(&room); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusOK, room)
}
