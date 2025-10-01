package services

import (
	"blockbuster/api/data"
	"context"

	"github.com/stretchr/testify/mock"
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

func (m *MockMembersService) GetIniitialVotingSlate(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(0)
}

func (m *MockMembersService) IterateRecommendationVoting(ctx context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, []string, error) {
	args := m.Called(ctx, currentMood, iteration, movieIDs)
	return args.Get(0).(data.MovieMetrics), args.Get(1).([]string), args.Error(0)
}

func (m *MockMembersService) GetVotingFinalPicks(ctx context.Context, mood data.MovieMetrics) ([]string, error) {
	args := m.Called(ctx, mood)
	return args.Get(0).([]string), args.Error(0)
}

func (m *MockMembersService) UpdateMood(ctx context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, error) {
	args := m.Called(ctx, currentMood, iteration, movieIDs)
	return args.Get(0).(data.MovieMetrics), args.Error(1)
}
