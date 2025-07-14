package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"blockbuster/api/constants"
	"blockbuster/api/services"
)

type MoviesHandler struct {
	service services.MoviesServiceInterface
}

func NewMoviesHandler() *MoviesHandler {
	return &MoviesHandler{
		service: services.GetMovieService(),
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
	movies, err := h.service.GetMoviesByPage(c.Request.Context(), pageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch movies"})
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
	movie, err := h.service.GetMovie(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Movie %s not found", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("err fetching movie wid id %s", id)})
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
	trivia, err := h.service.GetTrivia(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("err occurred fetching trivia for %s", id)})
	}
	if trivia.Trivia == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": fmt.Sprintf("Trivia for %s not found", id)})
		return
	}
	c.JSON(http.StatusOK, trivia)
}
