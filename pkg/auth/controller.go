package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ginn/pkg/middlewares"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/auth")
	routes.POST("/login", h.Login)
	routes.POST("/register", h.Register)
	routes.GET("/profile", middlewares.AuthorizeJWT(), h.Profile)
}
