package services

import (
	"context"

	"github.com/stretchr/testify/mock"

	"blockbuster/api/data"
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
