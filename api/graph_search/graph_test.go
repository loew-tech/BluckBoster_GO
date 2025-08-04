package graphsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"blockbuster/api/data"
)

func createSampleGraph() *MovieGraph {
	movie1 := data.Movie{ID: "1", Title: "A", Director: "Spielberg", Cast: []string{"Alice", "Bob"}}
	movie2 := data.Movie{ID: "2", Title: "B", Director: "Spielberg", Cast: []string{"Bob", "Charlie"}}
	movie3 := data.Movie{ID: "3", Title: "C", Director: "Nolan", Cast: []string{"Alice", "David"}}

	return &MovieGraph{
		directedMovies: map[string][]data.Movie{
			"Spielberg": {movie1, movie2},
			"Nolan":     {movie3},
		},
		starredIn: map[string][]data.Movie{
			"Alice":   {movie1, movie3},
			"Bob":     {movie1, movie2},
			"Charlie": {movie2},
			"David":   {movie3},
		},
		starredWith: map[string]map[string]bool{
			"Alice":   {"Bob": true, "David": true},
			"Bob":     {"Alice": true, "Charlie": true},
			"Charlie": {"Bob": true},
			"David":   {"Alice": true},
		},
		movieTitleToMovie: map[string]data.Movie{
			"A": movie1,
			"B": movie2,
			"C": movie3,
		},
		NumDirectors: 2,
		NumMovies:    3,
		NumStars:     4,
	}
}

func TestBFS(t *testing.T) {
	graph := createSampleGraph()
	stars := make(map[string]bool)
	movies := make(map[string]bool)
	directors := make(map[string]bool)

	graph.BFS("Alice", stars, movies, directors, 2)

	assert.True(t, stars["Alice"])
	assert.True(t, stars["Bob"])
	assert.True(t, stars["Charlie"])
	assert.True(t, stars["David"])
	assert.True(t, movies["A"])
	assert.True(t, movies["B"])
	assert.True(t, movies["C"])
	assert.True(t, directors["Spielberg"])
	assert.True(t, directors["Nolan"])
}

func TestGetDirectedMovies(t *testing.T) {
	graph := createSampleGraph()
	movies := graph.GetDirectedMovies("Spielberg")
	assert.Len(t, movies, 2)
}

func TestGetDirectedActors(t *testing.T) {
	graph := createSampleGraph()
	actors := graph.GetDirectedActors("Spielberg")
	assert.ElementsMatch(t, []string{"Alice", "Bob", "Charlie"}, actors)
}

func TestGetStarredIn(t *testing.T) {
	graph := createSampleGraph()
	movies := graph.GetStarredIn("Alice")
	assert.Len(t, movies, 2)
}

func TestGetStarredWith(t *testing.T) {
	graph := createSampleGraph()
	coStars := graph.GetStarredWith("Alice")
	assert.ElementsMatch(t, []string{"Bob", "David"}, coStars)
}

func TestGetMovieFromTitle(t *testing.T) {
	graph := createSampleGraph()
	movie, err := graph.GetMovieFromTitle("A")
	assert.NoError(t, err)
	assert.Equal(t, "1", movie.ID)
}

func TestGetIDFromTitle(t *testing.T) {
	graph := createSampleGraph()
	id, err := graph.GetIDFromTitle("C")
	assert.NoError(t, err)
	assert.Equal(t, "3", id)
}

func TestGetIDFromTitle_Invalid(t *testing.T) {
	graph := createSampleGraph()
	_, err := graph.GetIDFromTitle("Z")
	assert.Error(t, err)
}

func TestTotalCounts(t *testing.T) {
	graph := createSampleGraph()
	assert.Equal(t, 4, graph.TotalStars())
	assert.Equal(t, 3, graph.TotalMovies())
	assert.Equal(t, 2, graph.TotalDirectors())
}
