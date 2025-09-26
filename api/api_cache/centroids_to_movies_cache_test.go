package api_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomMovieFromCentroid_Found(t *testing.T) {
	cache := CentroidsToMoviesCache{
		CentroidToMovieIDs: map[int][]string{
			1: {"movieA", "movieB", "movieC"},
		},
	}
	// Should always return a movie from the list
	movie, err := cache.GetRandomMovieFromCentroid(1)
	assert.NoError(t, err)
	assert.Contains(t, []string{"movieA", "movieB", "movieC"}, movie)
}

func TestGetRandomMovieFromCentroid_NotFound(t *testing.T) {
	cache := CentroidsToMoviesCache{
		CentroidToMovieIDs: map[int][]string{},
	}
	movie, err := cache.GetRandomMovieFromCentroid(99)
	assert.Error(t, err)
	assert.Empty(t, movie)
}
