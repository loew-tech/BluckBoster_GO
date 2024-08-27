package main

import (
	"blockbuster/api/endpoints"
	"fmt"

	"github.com/gin-gonic/gin"
)


func main() {
	fmt.Println("hello world")
	router := gin.Default()
	router.GET("/movies", endpoints.GetMovies)
	
	
	router.Run("localhost:8080")
}