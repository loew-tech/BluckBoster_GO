package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/handlers"
)

type MockMembersService struct {
	mock.Mock
}

func (m *MockMembersService) GetMember(ctx context.Context, username string, isCart bool) (data.Member, error) {
	args := m.Called(ctx, username, isCart)
	return *args.Get(0).(*data.Member), args.Error(1)
}

func (m *MockMembersService) Login(ctx context.Context, username string) (data.Member, error) {
	args := m.Called(ctx, username)
	return *args.Get(0).(*data.Member), args.Error(1)
}

func (m *MockMembersService) AddToCart(ctx context.Context, username, movieID string) (bool, error) {
	args := m.Called(ctx, username, movieID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMembersService) RemoveFromCart(ctx context.Context, username, movieID string) (bool, error) {
	args := m.Called(ctx, username, movieID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMembersService) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	args := m.Called(ctx, username, movieIDs)
	return args.Get(0).([]string), args.Int(1), args.Error(2)
}

func (m *MockMembersService) GetCartIDs(ctx context.Context, username string) ([]string, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMembersService) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	args := m.Called(ctx, username, movieIDs)
	return args.Get(0).([]string), args.Int(1), args.Error(2)
}

func (m *MockMembersService) GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMembersService) GetCartMovies(ctx context.Context, username string) ([]data.Movie, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMembersService) SetAPIChoice(ctx context.Context, username, choice string) error {
	args := m.Called(ctx, username, choice)
	return args.Error(0)
}

// tests
func TestGetCartMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.GET("/members/:username/cart", h.GetCartMovies)

	mockMovies := []data.Movie{{ID: "movie1"}, {ID: "movie2"}}
	mockService.On("GetCartMovies", mock.Anything, "testuser").Return(mockMovies, nil)

	req, _ := http.NewRequest(http.MethodGet, "/members/testuser/cart", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var movies []data.Movie
	err := json.Unmarshal(resp.Body.Bytes(), &movies)
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
}

func TestReturn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.POST("/members/return", h.Return)

	mockService.On("Return", mock.Anything, "testuser", []string{"movie123"}).Return([]string{"returned movie123"}, 1, nil)

	body := `{"username":"testuser", "movie_ids":["movie123"]}`
	req, _ := http.NewRequest(http.MethodPost, "/members/return", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	assert.Contains(t, resp.Body.String(), "returned movie123")
}

func TestGetCheckedOutMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.GET("/members/:username/checkedout", h.GetCheckedOutMovies)

	mockMovies := []data.Movie{{ID: "movie1"}, {ID: "movie2"}}
	mockService.On("GetCheckedOutMovies", mock.Anything, "testuser").Return(mockMovies, nil)

	req, _ := http.NewRequest(http.MethodGet, "/members/testuser/checkedout", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var movies []data.Movie
	err := json.Unmarshal(resp.Body.Bytes(), &movies)
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
}

func TestSetAPIChoice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.PUT("/members/:username", h.SetAPIChoice)

	mockService.On("SetAPIChoice", mock.Anything, "testuser", constants.REST_API).Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/members/testuser?api_choice=REST", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "API choice set to REST")
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.POST("/members/login", h.Login)

	mockMember := &data.Member{Username: "testuser"}
	mockService.On("Login", mock.Anything, "testuser").Return(mockMember, nil)

	body := `{"username":"testuser"}`
	req, _ := http.NewRequest(http.MethodPost, "/members/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "testuser")
}

func TestAddToCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.PUT("/members/cart", h.AddToCart)

	mockService.On("AddToCart", mock.Anything, "testuser", "movie1").Return(true, nil)

	body := `{"username":"testuser", "movie_id":"movie1"}`
	req, _ := http.NewRequest(http.MethodPut, "/members/cart", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	assert.Contains(t, resp.Body.String(), "success")
}

func TestRemoveFromCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.PUT("/members/cart/remove", h.RemoveFromCart)

	mockService.On("RemoveFromCart", mock.Anything, "testuser", "movie1").Return(true, nil)

	body := `{"username":"testuser", "movie_id":"movie1"}`
	req, _ := http.NewRequest(http.MethodPut, "/members/cart/remove", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	assert.Contains(t, resp.Body.String(), "success")
}

func TestCheckout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.POST("/members/checkout", h.Checkout)

	mockService.On("Checkout", mock.Anything, "testuser", []string{"movie1", "movie2"}).Return([]string{"checked out movie1", "checked out movie2"}, 2, nil)

	body := `{"username":"testuser", "movie_ids":["movie1","movie2"]}`
	req, _ := http.NewRequest(http.MethodPost, "/members/checkout", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	assert.Contains(t, resp.Body.String(), "checked out movie1")
}

func TestGetCartIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.GET("/members/:username/cart/ids", h.GetCartIDs)

	mockService.On("GetCartIDs", mock.Anything, "testuser").Return([]string{"id1", "id2"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/members/testuser/cart/ids", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "id1")
}
