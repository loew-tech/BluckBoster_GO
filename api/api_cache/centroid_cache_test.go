package api_cache

import (
	"blockbuster/api/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricsByCentroid_Found(t *testing.T) {
	cache := CentroidCache{
		centroids: map[int]data.MovieMetrics{
			1: {Acting: 1.1, Action: 2.2},
		},
	}
	metrics, err := cache.GetMetricsByCentroid(1)
	assert.NoError(t, err)
	assert.Equal(t, 1.1, metrics.Acting)
	assert.Equal(t, 2.2, metrics.Action)
}

func TestGetMetricsByCentroid_NotFound(t *testing.T) {
	cache := CentroidCache{
		centroids: map[int]data.MovieMetrics{},
	}
	_, err := cache.GetMetricsByCentroid(99)
	assert.Error(t, err)
}

func TestGetKNearestCentroidsFromMood_Basic(t *testing.T) {
	cache := CentroidCache{
		centroids: map[int]data.MovieMetrics{
			1: {Acting: 1, Action: 1},
			2: {Acting: 2, Action: 2},
			3: {Acting: 3, Action: 3},
		},
	}
	mood := data.MovieMetrics{Acting: 2.1, Action: 2.1}
	ids, err := cache.GetKNearestCentroidsFromMood(mood, 2)
	assert.NoError(t, err)
	assert.Len(t, ids, 2)
	// Centroid 2 should be closest, then 3 or 1
	assert.Equal(t, 2, ids[0])
}

func TestGetKNearestCentroidsFromMood_KTooLarge(t *testing.T) {
	cache := CentroidCache{
		centroids: map[int]data.MovieMetrics{
			1: {Acting: 1, Action: 1},
			2: {Acting: 2, Action: 2},
		},
	}
	mood := data.MovieMetrics{Acting: 1, Action: 1}
	ids, err := cache.GetKNearestCentroidsFromMood(mood, 5)
	assert.NoError(t, err)
	assert.Len(t, ids, 2)
}

func TestGetKNearestCentroidsFromMood_KZeroOrNegative(t *testing.T) {
	cache := CentroidCache{
		centroids: map[int]data.MovieMetrics{
			1: {Acting: 1, Action: 1},
		},
	}
	mood := data.MovieMetrics{Acting: 1, Action: 1}
	_, err := cache.GetKNearestCentroidsFromMood(mood, 0)
	assert.Error(t, err)
	_, err = cache.GetKNearestCentroidsFromMood(mood, -1)
	assert.Error(t, err)
}
