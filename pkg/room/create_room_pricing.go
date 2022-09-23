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

type CreateRoomPricing struct {
	Price float32 `json:"price" binding:"required,min=1"`
	Type  string  `json:"type" binding:"required,Enum=shift_daily_weekly"`
	Room  string  `json:"room" binding:"required,uuid"`
}

func (h handler) CreateRoomPricing(c *gin.Context) {
	user := authorization.ExtractUser(c)

	body := CreateRoomPricing{}

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
	if err := h.DB.Where("id = ?", body.Room).Preload("Building").First(&room).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NotFound("Room not found"))
		return
	}

	if user.ID != room.Building.UserID {
		c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Unauthorized("Unauthorized"))
		return
	}

	pricing := models.RoomPricing{
		Price:  body.Price,
		Type:   body.Type,
		RoomID: body.Room,
	}

	if result := h.DB.Save(&pricing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusCreated, pricing)
}
