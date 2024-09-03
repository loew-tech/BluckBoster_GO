package endpoints

import (
	"blockbuster/api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var movieRepo = db.NewMovieRepo()

func GetMoviesEndpoint(c *gin.Context) {
	movies, err := movieRepo.GetAllMovies()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Failed to retrieve movies"})
	} else {
		c.IndentedJSON(http.StatusOK, movies)
	}
}
