package schedule

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

	routes := r.Group("/schedules")

	authRoutes := routes.Use(middlewares.AuthorizeJWT())
	authRoutes.POST("/", h.CreateSchedule)
	authRoutes.PUT("/:id", h.UpdateSchedule)
}
