package schedule

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookScheduleRequestBody struct {
	Schedule string `json:"schedule" binding:"required,uuid"`
}

func (h handler) BookSchedule(c *gin.Context) {
	user := authorization.ExtractUser(c)
	body := BookScheduleRequestBody{}

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
	if err := h.DB.Where("id = ?", body.Schedule).Preload("Room").Preload("Room.Building").First(&schedule).Error; err != nil {
		c.JSON(http.StatusNotFound, exceptions.NotFound("Schedule not found"))
		return
	}

	if schedule.Room.Building.UserID == user.ID {
		c.JSON(http.StatusForbidden, exceptions.Unauthorized("You are not allowed to book your own room"))
		return
	}

	if schedule.Status != models.ScheduleStatusOpen {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Room already scheduled"))
		return
	}

	schedule.Status = models.ScheduleStatusReserved
	schedule.UserID.Scan(user.ID)
	// TODO stripe payment and etc

	if result := h.DB.Save(&schedule); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusCreated, schedule)
}
