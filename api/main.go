package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/gql"
	"blockbuster/api/handlers"
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

	// === Create and register handlers ===
	membersHandler := handlers.NewMembersHandler()
	moviesHandler := handlers.NewMoviesHandler()

	// === register routes ===
	api := router.Group(constants.REST_ROUTER_GROUP)
	membersHandler.RegisterRoutes(api)
	moviesHandler.RegisterRoutes(api)

	// === GraphQL endpoint ===
	router.POST(constants.GRAPHQL_ENDPOINT, gqlHandler)

	router.Run(LOCAL_HOST)
}
