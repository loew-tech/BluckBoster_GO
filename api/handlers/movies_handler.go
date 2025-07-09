package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

type MoviesHandler struct {
	repo repos.ReadWriteMovieRepo
}

func NewMoviesHandler() *MoviesHandler {
	return &MoviesHandler{
		repo: repos.NewMovieRepoWithDynamo(),
	}
}

func (h *MoviesHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/movies", h.GetMoviesByPage)
	rg.GET("/movies/:movieID", h.GetMovie)
	rg.GET("/movies/:movieID/trivia", h.GetTrivia)
}

func (h *MoviesHandler) GetMoviesByPage(c *gin.Context) {
	pageKey := strings.ToUpper(c.DefaultQuery(constants.PAGE, constants.DEFAULT_PAGE))
	if pageKey == "" {
		log.Println("Missing 'page' query parameter")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing page query parameter"})
		return
	}
	if !strings.Contains(constants.PAGES, pageKey) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid page key: %s", pageKey)})
		return
	}
	log.Printf("Fetching movies for page key: %s", pageKey)
	movies, err := h.repo.GetMoviesByPage(c, constants.NOT_FOR_GRAPH, pageKey)
	if err != nil {
		utils.LogError("Failed to fetch movies", err)
		c.JSON(http.StatusNotFound, gin.H{"msg": "Failed to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MoviesHandler) GetMovie(c *gin.Context) {
	id := c.Param(constants.MOVIE_ID)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing movieID parameter"})
		return
	}
	log.Printf("Fetching movie by ID: %s", id)
	movie, err := h.repo.GetMovieByID(c, id, false)
	if err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve movie %s", id), err)
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Movie %s not found", id)})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *MoviesHandler) GetTrivia(c *gin.Context) {
	id := c.Param(constants.MOVIE_ID)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing movieID parameter"})
		return
	}
	log.Printf("Fetching trivia for movie ID: %s", id)
	trivia, err := h.repo.GetTrivia(c, id)
	if err != nil {
		utils.LogError(fmt.Sprintf("Failed to retrieve trivia for movie %s", id), err)
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Trivia for %s not found", id)})
		return
	}
	c.JSON(http.StatusOK, trivia)
}
