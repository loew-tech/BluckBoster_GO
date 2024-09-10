package main

import (
	"blockbuster/api/endpoints"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const LOCAL_HOST = "localhost:8080"

func main() {
	fmt.Println("hello world")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == LOCAL_HOST
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/api/v1/movies", endpoints.GetMoviesEndpoint)
	router.GET("/api/v1/members/:username", endpoints.GetMemberEndpoint)
	router.POST("/api/v1/members/login", endpoints.MemberLoginEndpoint)
	router.PUT("/api/v1/members/cart", endpoints.AddToCartEndpoint)
	router.GET("/api/v1/members/cart/ids/:username", endpoints.GetCartIDsEndpoint)

	router.Run(LOCAL_HOST)
}
