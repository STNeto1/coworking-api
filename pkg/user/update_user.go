package user

import (
	"net/http"

	"coworking/pkg/common/authorization"
	"coworking/pkg/common/exceptions"
	"coworking/pkg/common/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h handler) UpdateUser(c *gin.Context) {
	user := authorization.ExtractUser(c)

	body := UpdateUserRequestBody{}

	// // getting request's body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Invalid request body"))
		return
	}

	if body.Name != "" {
		user.Name = body.Name
	}

	if body.Email != "" && body.Email != user.Email {
		existingUser := models.User{}
		h.DB.First(&existingUser, "email = ?", body.Email)
		if existingUser.ID != "" {
			c.JSON(http.StatusBadRequest, exceptions.BadRequest("Email already in use"))
			return
		}

		user.Email = body.Email
	}

	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
			return
		}
		user.Password = string(hashedPassword)
	}

	if result := h.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("Internal error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
