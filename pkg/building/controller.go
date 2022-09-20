package building

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

	routes := r.Group("/buildings")
	routes.GET("/", h.GetAllBuildings)

	authRoutes := routes.Use(middlewares.AuthorizeJWT())
	authRoutes.POST("/", h.CreateBuilding)
	authRoutes.GET("/user", h.GetUserBuildings)
	authRoutes.PUT("/:id", h.UpdateBuilding)
	authRoutes.DELETE("/:id", h.DeleteBuilding)
}
