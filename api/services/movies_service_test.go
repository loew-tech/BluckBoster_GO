package services_test

import (
	"context"
	"errors"
	"testing"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMovieRepo struct {
	mock.Mock
}

func (m *MockMovieRepo) GetMoviesByPage(ctx context.Context, page string, forGraph bool) ([]data.Movie, error) {
	args := m.Called(ctx, page, forGraph)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMovieRepo) GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error) {
	args := m.Called(ctx, movieID, forCart)
	return args.Get(0).(data.Movie), args.Error(1)
}

func (m *MockMovieRepo) GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error) {
	args := m.Called(ctx, movieIDs, forCart)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMovieRepo) GetMovieMetrics(ctx context.Context, movieID string) (data.MovieMetrics, error) {
	args := m.Called(ctx, movieID)
	return args.Get(0).(data.MovieMetrics), args.Error(1)
}

func (m *MockMovieRepo) GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error) {
	args := m.Called(ctx, movieID)
	return args.Get(0).(data.MovieTrivia), args.Error(1)
}

func setupMockMovieService() (*services.MoviesService, *MockMovieRepo) {
	repo := new(MockMovieRepo)
	service := services.NewMovieserviceWithRepo(repo)
	return service, repo
}

func TestGetMoviesByPage_Success(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMoviesByPage", mock.Anything, "A", constants.NOT_FOR_GRAPH).
		Return([]data.Movie{{ID: "1"}, {ID: "2"}}, nil)

	movies, err := service.GetMoviesByPage(context.Background(), "A")
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
}

func TestGetMoviesByPage_Error(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMoviesByPage", mock.Anything, "B", constants.NOT_FOR_GRAPH).
		Return([]data.Movie{}, errors.New("db error"))

	_, err := service.GetMoviesByPage(context.Background(), "B")
	assert.Error(t, err)
}

func TestGetMovie_Success(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMovieByID", mock.Anything, "m1", constants.NOT_CART).
		Return(data.Movie{ID: "m1"}, nil)

	movie, err := service.GetMovie(context.Background(), "m1")
	assert.NoError(t, err)
	assert.Equal(t, "m1", movie.ID)
}

func TestGetMovie_Error(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMovieByID", mock.Anything, "bad", constants.NOT_CART).
		Return(data.Movie{}, errors.New("not found"))

	_, err := service.GetMovie(context.Background(), "bad")
	assert.Error(t, err)
}

func TestGetMovies_Success(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMoviesByID", mock.Anything, []string{"m1", "m2"}, constants.CART).
		Return([]data.Movie{{ID: "m1"}, {ID: "m2"}}, nil)

	movies, err := service.GetMovies(context.Background(), []string{"m1", "m2"})
	assert.NoError(t, err)
	assert.Len(t, movies, 2)
}

func TestGetMovies_Error(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetMoviesByID", mock.Anything, []string{"bad"}, constants.CART).
		Return([]data.Movie{}, errors.New("db fail"))

	_, err := service.GetMovies(context.Background(), []string{"bad"})
	assert.Error(t, err)
}

func TestGetTrivia_Success(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetTrivia", mock.Anything, "m1").
		Return(data.MovieTrivia{Trivia: "fact"}, nil)

	trivia, err := service.GetTrivia(context.Background(), "m1")
	assert.NoError(t, err)
	assert.Equal(t, trivia.Trivia, "fact")
}

func TestGetTrivia_Error(t *testing.T) {
	service, repo := setupMockMovieService()
	repo.On("GetTrivia", mock.Anything, "bad").
		Return(data.MovieTrivia{}, errors.New("fail"))

	_, err := service.GetTrivia(context.Background(), "bad")
	assert.Error(t, err)
}
