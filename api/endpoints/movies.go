package endpoints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	repos "blockbuster/api/repos"
)

var movieRepo = repos.NewMovieRepo(GetDynamoClient())

const VALID_KEYS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ#"

func GetMoviesByPageEndpoint(c *gin.Context) {
	pageKey := strings.ToUpper(c.DefaultQuery("page", "A"))
	if !strings.Contains(VALID_KEYS, pageKey) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("No movies associated with page key: %s", pageKey)})
		return
	}
	movies, err := movieRepo.GetMoviesByPage(c, constants.NOT_FOR_GRAPH, pageKey)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve movies"})
	} else {
		c.IndentedJSON(http.StatusOK, movies)
	}
}

func GetMovieEndpoint(c *gin.Context) {
	movieID := c.Param("movieID")
	movie, err := movieRepo.GetMovieByID(c, movieID, false)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Failed to retrieve movieID %s", movieID)})
	} else {
		c.IndentedJSON(http.StatusOK, movie)
	}
}

func GetTriviaEndpoint(c *gin.Context) {
	movieID := c.Param("movieID")
	trivia, err := movieRepo.GetTrivia(c, movieID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Failed to retrieve movieID %s", movieID)})
	} else {
		c.IndentedJSON(http.StatusOK, trivia)
	}
}
