package authorization

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"coworking/pkg/common/db"
	"coworking/pkg/common/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func getSecret() string {
	return viper.Get("SECRET").(string)
}

func GenerateToken(email string) string {
	claims := &jwt.RegisteredClaims{
		Issuer:    "ISSUER",
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(getSecret()))
	if err != nil {
		panic(err)
	}
	return t
}

func ValidateToken(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(getSecret()), nil
	})

	if err != nil {
		return err
	}

	return nil

}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenID(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecret()), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["sub"].(string), nil
	}

	return "", nil
}

func ExtractUser(c *gin.Context) models.User {
	h := db.Init(viper.Get("DB_URL").(string))
	id, err := ExtractTokenID(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		c.Abort()

		return models.User{}
	}

	user := models.User{}
	h.First(&user, "id = ?", id)

	if user.ID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		c.Abort()

		return models.User{}
	}

	return user
}
