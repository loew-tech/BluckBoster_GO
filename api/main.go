package main

import (
	"blockbuster/api/endpoints"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello world")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"*"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return origin == "localhost:8080"
        },
        MaxAge: 12 * time.Hour,
    }))

	router.GET("/api/v1/movies", endpoints.GetMovies)
	
	
	router.Run("localhost:8080")
}