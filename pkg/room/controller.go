package room

import (
	"coworking/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := r.Group("/rooms")
	routes.GET("/", h.GetAllRooms)
	routes.GET("/:id", h.ShowRoom)

	authRoutes := routes.Use(middlewares.AuthorizeJWT())
	authRoutes.POST("/pricing", h.CreateRoomPricing)
	authRoutes.POST("/", h.CreateRoom)
	authRoutes.PUT("/:id", h.UpdateRoom)
	authRoutes.DELETE("/:id", h.DeleteRoom)
}
