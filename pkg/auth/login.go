package auth

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h handler) Login(c *gin.Context) {
	body := LoginRequestBody{}

	// getting request's body
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

	user := models.User{}
	h.DB.First(&user, "email = ?", body.Email)
	if user.ID == "" {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Invalid credentials"))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Invalid credentials"))
		return
	}

	token := authorization.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
