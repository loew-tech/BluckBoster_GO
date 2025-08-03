package graphsearch

import (
	"github.com/stretchr/testify/mock"
	
	"blockbuster/api/data"
)

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