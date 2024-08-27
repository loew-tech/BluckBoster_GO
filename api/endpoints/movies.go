package endpoints

import (
	"blockbuster/api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var repo = db.NewMovieRepo()

func GetMovies(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, repo.GetAllMovies())
}