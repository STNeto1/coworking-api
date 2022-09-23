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

type UpdateRoomPricing struct {
	Price float32 `json:"price" binding:"required,min=1"`
	Type  string  `json:"type" binding:"required,Enum=shift_daily_weekly"`
}

func (h handler) UpdateRoomPricing(c *gin.Context) {
	id := c.Param("id")
	user := authorization.ExtractUser(c)

	body := UpdateRoomPricing{}

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

	var pricing models.RoomPricing
	if err := h.DB.Where("id = ?", id).Preload("Room.Building").First(&pricing).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	if user.ID != pricing.Room.Building.UserID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Unauthorized("Unauthorized"))
		return
	}

	pricing.Price = body.Price
	pricing.Type = body.Type

	if result := h.DB.Save(&pricing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusOK, pricing)
}
