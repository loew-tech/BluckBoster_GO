package endpoints

import (
	"blockbuster/api/data"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var movieRepo = data.NewMovieRepo(GetDynamoClient())

func GetMoviesEndpoint(c *gin.Context) {
	movies, err := movieRepo.GetAllMovies()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve movies"})
	} else {
		c.IndentedJSON(http.StatusOK, movies)
	}
}

func GetMovieEndpoint(c *gin.Context) {
	movieID := c.Param("movieID")
	movie, _, err := movieRepo.GetMovieByID(movieID, false)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Failed to retrieve movieID %s", movieID)})
	} else {
		c.IndentedJSON(http.StatusOK, movie)
	}
}

func GetMovieTriviaEndpoint(c *gin.Context) {
	data.GetMovieTrivia(strings.ReplaceAll(c.Param("movieID"), "_", " "))
	c.IndentedJSON(http.StatusNotImplemented, gin.H{"msg": "Trivia Endpoint"})
}
