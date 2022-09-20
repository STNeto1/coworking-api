package building

import (
	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type addressRequestBody struct {
	PostalCode    string `json:"postal_code" binding:"required"`
	State         string `json:"state" binding:"required"`
	City          string `json:"city" binding:"required"`
	Street        string `json:"street" binding:"required"`
	Number        string `json:"number" binding:"required"`
	GoogleMapsURL string `json:"google_maps_url" binding:"required,url"`
}

type CreateBuildingRequestBody struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Address     *addressRequestBody `json:"address" binding:"required"`
}

func (h handler) CreateBuilding(c *gin.Context) {
	user := authorization.ExtractUser(c)
	body := CreateBuildingRequestBody{}

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

	h.DB.Transaction(func(tx *gorm.DB) error {
		address := models.Address{
			PostalCode:    body.Address.PostalCode,
			State:         body.Address.State,
			City:          body.Address.City,
			Street:        body.Address.Street,
			Number:        body.Address.Number,
			GoogleMapsURL: body.Address.GoogleMapsURL,
		}

		if err := tx.Create(&address).Error; err != nil {
			return err
		}

		building := models.Building{
			Name:        body.Name,
			Description: body.Description,
			AddressID:   address.ID,
			UserID:      user.ID,
		}

		if err := tx.Create(&building).Error; err != nil {
			return err
		}

		return nil
	})

	c.JSON(http.StatusCreated, gin.H{})
}
