package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	baseRoute := "api"

	// Update CORS policies
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"*"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = []string{"*"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	r.GET(baseRoute, func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Connected to server!"})
	})

	r.Run(":8000") // Defaults to localhost:8080 when no port given
}
