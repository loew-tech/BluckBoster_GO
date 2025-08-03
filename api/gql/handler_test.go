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

	graphsearch "blockbuster/api/graph_search"
)

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

	mockGraph := new(graphsearch.MockMovieGraph)
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
		val := ctx.Value(GinContextKey)
		assert.Equal(t, c, val)
	})

	setTestMovieGraph(&graphsearch.MockMovieGraph{})

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
	ctx := context.WithValue(context.Background(), GinContextKey, "test")
	val := ctx.Value(GinContextKey)
	assert.Equal(t, "test", val)
}
