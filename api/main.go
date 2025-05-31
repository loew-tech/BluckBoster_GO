package main

import (
	"fmt"
	"time"

	"blockbuster/api/endpoints"
	"blockbuster/api/gql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const LOCAL_HOST = "localhost:8080"

func main() {
	fmt.Println("hello world")

	gqlHandler := gql.GetGQLHandler()

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

	router.GET("/api/v1/movies", endpoints.GetMoviesByPageEndpoint)
	router.GET("/api/v1/movies/:movieID", endpoints.GetMovieEndpoint)
	router.GET("/api/v1/movies/:movieID/trivia", endpoints.GetTriviaEndpoint)
	router.GET("/api/v1/members/:username", endpoints.GetMemberEndpoint)
	router.POST("/api/v1/members/login", endpoints.MemberLoginEndpoint)
	router.GET("/api/v1/members/:username/cart", endpoints.GetCartMoviesEndpoint)
	// @TODO: remove this once FE cors issue is resolved
	router.GET("/api/v1/members/cart/:username", endpoints.GetCartMoviesEndpoint)
	router.PUT("/api/v1/members/cart", endpoints.AddToCartEndpoint)
	router.PUT("/api/v1/members/cart/remove", endpoints.RemoveFromCartEndpoint)
	router.GET("/api/v1/members/:username/cart/ids", endpoints.GetCartIDsEndpoint)
	router.POST("/api/v1/members/checkout", endpoints.CheckoutEndpoint)
	router.GET("api/v1/members/:username/checkedout", endpoints.GetCheckedOutMovies)
	router.POST("/api/v1/members/return", endpoints.ReturnEndpoint)
	router.POST("/graphql/v1", gqlHandler)

	router.Run(LOCAL_HOST)
}
