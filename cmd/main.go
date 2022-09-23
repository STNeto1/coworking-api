package main

import (
	"coworking/pkg/auth"
	"coworking/pkg/building"
	"coworking/pkg/common/config"
	"coworking/pkg/common/db"
	"coworking/pkg/common/validators"
	"coworking/pkg/room"
	"coworking/pkg/user"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func main() {
	production := os.Getenv("RAILWAY_ENVIRONMENT") == "production"

	if !production {
		viper.SetConfigFile("./pkg/common/envs/.env")
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	conf, err := config.LoadConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error loading config: %w", err))
	}

	r := gin.Default()
	h := db.Init(conf.DBUrl)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("Enum", validators.Enum)
	}

	auth.RegisterRoutes(r, h)
	user.RegisterRoutes(r, h)
	building.RegisterRoutes(r, h)
	room.RegisterRoutes(r, h)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Gin API!",
		})
	})

	r.Run()
}
