package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	repos "blockbuster/api/repos"
)

var movieRepo = repos.NewMovieRepo(GetDynamoClient())

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
