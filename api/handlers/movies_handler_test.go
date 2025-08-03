package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/data"
	"blockbuster/api/handlers"
)

type MockMoviesService struct {
	mock.Mock
}

func (m *MockMoviesService) GetMoviesByPage(ctx context.Context, pageKey string) ([]data.Movie, error) {
	args := m.Called(ctx, pageKey)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMoviesService) GetMovie(ctx context.Context, id string) (data.Movie, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(data.Movie), args.Error(1)
}

func (m *MockMoviesService) GetMovies(ctx context.Context, ids []string) ([]data.Movie, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMoviesService) GetTrivia(ctx context.Context, id string) (data.MovieTrivia, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(data.MovieTrivia), args.Error(1)
}

func TestGetMoviesByPage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies", h.GetMoviesByPage)

	expected := []data.Movie{{ID: "m1"}, {ID: "m2"}}
	mockService.On("GetMoviesByPage", mock.Anything, "A").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/movies?page=A", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var movies []data.Movie
	err := json.Unmarshal(resp.Body.Bytes(), &movies)
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
}

func TestGetMoviesByPage_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies", h.GetMoviesByPage)

	req, _ := http.NewRequest(http.MethodGet, "/movies?page=INVALID", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid page key")
}

func TestGetMovie_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies/:movieID", h.GetMovie)

	expected := data.Movie{ID: "m123"}
	mockService.On("GetMovie", mock.Anything, "m123").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/movies/m123", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "m123")
}

func TestGetMovie_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies/:movieID", h.GetMovie)

	mockService.On("GetMovie", mock.Anything, "m999").Return(data.Movie{}, errors.New("movie not found"))

	req, _ := http.NewRequest(http.MethodGet, "/movies/m999", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Movie m999 not found")
}

func TestGetTrivia_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies/:movieID/trivia", h.GetTrivia)

	expected := data.MovieTrivia{Trivia: "Some trivia"}
	mockService.On("GetTrivia", mock.Anything, "m123").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/movies/m123/trivia", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Some trivia")
}

func TestGetTrivia_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMoviesService)
	h := handlers.NewMoviesHandlerWithService(mockService)
	r.GET("/movies/:movieID/trivia", h.GetTrivia)

	expected := data.MovieTrivia{Trivia: ""}
	mockService.On("GetTrivia", mock.Anything, "m123").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/movies/m123/trivia", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Trivia for m123 not found")
}
