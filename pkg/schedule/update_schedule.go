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

type UpdateScheduleRequestBody struct {
	Price float32 `json:"price" binding:"required,gte=0"`
	From  string  `json:"from" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	To    string  `json:"to" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

func (h handler) UpdateSchedule(c *gin.Context) {
	user := authorization.ExtractUser(c)
	id := c.Param("id")
	body := UpdateScheduleRequestBody{}

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

	var schedule models.Schedule
	if err := h.DB.Where("id = ?", id).Preload("Room").Preload("Room.Building").First(&schedule).Error; err != nil {
		c.JSON(http.StatusNotFound, exceptions.NotFound("Schedule not found"))
		return
	}

	if schedule.Room.Building.UserID != user.ID {
		c.JSON(http.StatusForbidden, exceptions.Unauthorized("You are not allowed to update this schedule"))
		return
	}

	if schedule.Status != models.ScheduleStatusOpen {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Room already scheduled"))
		return
	}

	from, _ := time.Parse("2006-01-02T15:04:05Z07:00", body.From)
	to, _ := time.Parse("2006-01-02T15:04:05Z07:00", body.To)

	if schedule.From.Sub(from) != 0 || schedule.To.Sub(to) != 0 {
		var overlappingSchedules []models.Schedule
		if err := h.DB.Where("room_id = ? AND (\"from\", \"to\") overlaps (?, ?) AND id != ?", schedule.Room.ID, from, to, schedule.ID).Find(&overlappingSchedules).Error; err != nil {
			c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
			return
		}

		if len(overlappingSchedules) > 0 {
			c.JSON(http.StatusBadRequest, exceptions.BadRequest("There is already a schedule for this room"))
			return
		}
	}

	schedule.Price = body.Price
	schedule.From = from
	schedule.To = to

	if result := h.DB.Save(&schedule); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusOK, schedule)
}
