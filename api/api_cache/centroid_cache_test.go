package centroidcache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"blockbuster/api/data"
)

func TestSetAndGetMetricsByCentroid(t *testing.T) {

	centroids := map[int]data.MovieMetrics{
		0: {Acting: 10, Action: 20},
		1: {Acting: 30, Action: 40},
	}
	cache := &CentroidCache{centroids: centroids}

	metrics, err := cache.GetMetricsByCentroid(0)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, metrics.Acting)
	assert.Equal(t, 20.0, metrics.Action)

	_, err = cache.GetMetricsByCentroid(2)
	assert.Error(t, err)
}

func TestGetCentroidsFromMood(t *testing.T) {

	centroids := map[int]data.MovieMetrics{
		0: {Acting: 11, Action: 11},
		1: {Acting: 20, Action: 20},
		2: {Acting: 30, Action: 30},
	}
	cache := &CentroidCache{centroids: centroids}

	// Closest to Acting:15, Action:15 should be centroid 1, then 0, then 2
	mood := data.MovieMetrics{Acting: 15, Action: 15}
	ids, err := cache.GetCentroidsFromMood(mood, 2)
	assert.NoError(t, err)
	assert.Len(t, ids, 2)
	assert.Equal(t, 0, ids[0])
	assert.Equal(t, 1, ids[1])

	// k larger than available centroids
	ids, err = cache.GetCentroidsFromMood(mood, 5)
	assert.NoError(t, err)
	assert.Len(t, ids, 3)

	// k = 0 should error
	_, err = cache.GetCentroidsFromMood(mood, 0)
	assert.Error(t, err)
}
