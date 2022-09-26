package schedule

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateScheduleRequestBody struct {
	Room  string  `json:"room" binding:"required,uuid"`
	Price float32 `json:"price" binding:"required,gte=0"`
	From  string  `json:"from" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	To    string  `json:"to" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func (h handler) CreateSchedule(c *gin.Context) {
	user := authorization.ExtractUser(c)
	body := CreateScheduleRequestBody{}

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
		c.JSON(http.StatusNotFound, exceptions.NotFound("Room not found"))
		return
	}

	if room.Building.UserID != user.ID {
		c.JSON(http.StatusForbidden, exceptions.Unauthorized("You are not allowed to create schedule for this room"))
		return
	}

	from, _ := time.Parse("2006-01-02T15:04:05Z07:00", body.From)
	to, _ := time.Parse("2006-01-02T15:04:05Z07:00", body.To)
	var overlappingSchedules []models.Schedule
	if err := h.DB.Where("room_id = ? AND (\"from\", \"to\") overlaps (?, ?)", body.Room, from, to).Find(&overlappingSchedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	if len(overlappingSchedules) > 0 {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("There is already a schedule for this room"))
		return
	}

	schedule := models.Schedule{
		RoomID: room.ID,
		Price:  body.Price,
		From:   from,
		To:     to,
		Status: models.ScheduleStatusOpen,
	}

	if result := h.DB.Create(&schedule); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusCreated, schedule)
}
