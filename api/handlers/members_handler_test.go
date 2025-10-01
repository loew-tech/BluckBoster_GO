package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/handlers"
	"blockbuster/api/services"
)

func TestGetCartMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
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
	mockService := new(services.MockMembersService)
	h := handlers.NewMembersHandlerWithService(mockService)
	r.GET("/members/:username/cart/ids", h.GetCartIDs)

	mockService.On("GetCartIDs", mock.Anything, "testuser").Return([]string{"id1", "id2"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/members/testuser/cart/ids", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "id1")
}

func TestGetVotingFinalPicksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockReturn     []string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			requestBody: gin.H{
				"current_mood": gin.H{"acting": 50, "action": 40, "cinematography": 30},
			},
			mockReturn:     []string{"m1", "m2", "m3"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"movies":["m1","m2","m3"]`,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"current_mood": "oops"}`, // wrong type
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"msg":"Invalid request body"`,
		},
		{
			name: "service error",
			requestBody: gin.H{
				"current_mood": gin.H{"acting": 50, "action": 40, "cinematography": 30},
			},
			mockReturn:     nil,
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"msg":"db error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(services.MockMembersService)
			handler := handlers.NewMembersHandlerWithService(mockSvc)

			// Only set expectation if valid JSON and expecting service call
			if tt.mockReturn != nil || tt.mockError != nil {
				mockSvc.On("GetVotingFinalPicks", mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError)
			}

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)
			r.POST("/members/mood/final_picks", handler.GetVotingFinalPicks)

			var reqBody []byte
			switch b := tt.requestBody.(type) {
			case string:
				reqBody = []byte(b)
			default:
				reqBody, _ = json.Marshal(b)
			}
			req, _ := http.NewRequest(http.MethodPost, "/members/mood/final_picks", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestUpdateMoodHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           interface{}
		mockReturn     data.MovieMetrics
		mockError      error
		expectedStatus int
	}{
		{
			name: "success",
			body: gin.H{
				"current_mood": gin.H{"acting": 50, "action": 30, "cinematography": 70},
				"iteration":    2,
				"movie_ids":    []string{"m1", "m2"},
			},
			mockReturn:     data.MovieMetrics{Acting: 60, Action: 40, Cinematography: 80},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "bad request - invalid json",
			body: `{"current_mood": "oops"}`, // wrong type
			// no mock expectations (won’t call service)
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: gin.H{
				"current_mood": gin.H{"acting": 10, "action": 20, "cinematography": 30},
				"iteration":    1,
				"movie_ids":    []string{"m1"},
			},
			mockReturn:     data.MovieMetrics{},
			mockError:      errors.New("db down"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(services.MockMembersService)
			handler := handlers.NewMembersHandlerWithService(mockSvc)

			// Only set expectation if valid body and not bad JSON
			if tt.expectedStatus != http.StatusBadRequest {
				mockSvc.On("UpdateMood", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockError)
			}

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)
			r.PUT("/members/mood", handler.UpdateMood)

			var reqBody []byte
			switch b := tt.body.(type) {
			case string: // raw invalid JSON
				reqBody = []byte(b)
			default:
				reqBody, _ = json.Marshal(b)
			}
			req, _ := http.NewRequest(http.MethodPut, "/members/mood", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp data.MovieMetrics
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn, resp)
			}
		})
	}
}
