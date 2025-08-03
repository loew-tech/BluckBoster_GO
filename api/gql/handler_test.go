package gql

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/data"
	graphsearch "blockbuster/api/graph_search"
)

// --- Mock for MovieGraphInterface ---
type MockMovieGraph struct {
	mock.Mock
}

func (m *MockMovieGraph) BFS(
	startStar string,
	stars map[string]bool,
	movieTitles map[string]bool,
	directors map[string]bool,
	maxDepth int,
) {
	m.Called(startStar, stars, movieTitles, directors, maxDepth)
}

func (m *MockMovieGraph) GetNumDirectors() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMovieGraph) GetNumStars() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMovieGraph) GetNumMovies() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMovieGraph) GetDirectedMovies(director string) []data.Movie {
	args := m.Called()
	return args.Get(0).([]data.Movie)
}

func (m *MockMovieGraph) GetStarredWith(star string) []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockMovieGraph) GetStarredIn(star string) []data.Movie {
	args := m.Called()
	return args.Get(0).([]data.Movie)
}

func (m *MockMovieGraph) GetMovieFromTitle(title string) (data.Movie, error) {
	args := m.Called(title)
	return args.Get(0).(data.Movie), args.Error(1)
}

func (m *MockMovieGraph) GetStarredInMovies(star string) ([]data.Movie, error) {
	args := m.Called(star)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMovieGraph) GetDirectedActors(director string) []string {
	args := m.Called(director)
	return args.Get(0).([]string)
}

func (m *MockMovieGraph) TotalStars() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMovieGraph) TotalMovies() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMovieGraph) TotalDirectors() int {
	args := m.Called()
	return args.Int(0)
}

// --- Test setup override ---
func setTestMovieGraph(mockGraph graphsearch.MovieGraphInterface) {
	initMovieGraphOnce = sync.Once{} // reset the sync.Once
	movieGraph = mockGraph
	initMovieGraphErr = nil
}

func TestGetGQLHandler_ValidRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/graphql", GetGQLHandler())

	mockGraph := new(MockMovieGraph)
	mockGraph.On("GetGraph").Return(nil)
	setTestMovieGraph(mockGraph)

	body := bytes.NewBufferString(`{"query":"{ dummy }"}`)
	req, err := http.NewRequest(http.MethodPost, "/graphql", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)        // GraphQL returns 200 even on errors
	assert.Contains(t, w.Body.String(), "errors") // No schema so it errors
}

func TestGetGQLHandler_ContextInjected(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/graphql", func(c *gin.Context) {
		h := GetGQLHandler()
		h(c)
		// Assert context value was injected
		ctx := c.Request.Context()
		val := ctx.Value(ginContextKey)
		assert.Equal(t, c, val)
	})

	setTestMovieGraph(&MockMovieGraph{})

	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewBufferString(`{"query":"{ dummy }"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)
}

func TestGetGQLHandler_InitError(t *testing.T) {
	// Arrange
	initMovieGraphOnce = sync.Once{} // Reset sync.Once
	movieGraph = nil
	initMovieGraphErr = errors.New("init failed") // Simulate failure

	// Just make sure it still returns a handler
	h := GetGQLHandler()
	assert.NotNil(t, h)
}

func TestGinContextKey_Isolated(t *testing.T) {
	// Ensures ginContextKey doesn't collide
	ctx := context.WithValue(context.Background(), ginContextKey, "test")
	val := ctx.Value(ginContextKey)
	assert.Equal(t, "test", val)
}
