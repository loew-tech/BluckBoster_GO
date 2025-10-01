package api_cache

import (
	"fmt"
	"math"
	"sort"

	"blockbuster/api/data"
	"blockbuster/api/utils"
)

type CentroidCache struct {
	centroids map[int]data.MovieMetrics
}

func (c *CentroidCache) GetMetricsByCentroid(centroidID int) (data.MovieMetrics, error) {
	if c.centroids == nil {
		return data.MovieMetrics{}, utils.LogError("centroid cache failed to initialize. GetMetricsByCentroid functionality unavailable", nil)
	}

	centroid, ok := c.centroids[centroidID]
	if !ok {
		return data.MovieMetrics{}, fmt.Errorf("centroid %d not found", centroidID)
	}
	return centroid, nil
}

func (c *CentroidCache) GetKNearestCentroidsFromMood(mood data.MovieMetrics, k int) ([]int, error) {
	if c.centroids == nil {
		return nil, utils.LogError("centroid cache failed to initialize. GetKNearestCentroidsFromMood functionality unavailable", nil)
	}
	if k <= 0 {
		return nil, utils.LogError("k must be greater than 0", nil)
	}

	type centroidDist struct {
		id   int
		dist float64
	}
	var dists []centroidDist
	for id, metrics := range c.centroids {
		d := utils.MetricDistance(mood, metrics)
		dists = append(dists, centroidDist{id: id, dist: d})
	}

	sort.Slice(dists, func(i, j int) bool {
		return dists[i].dist < dists[j].dist
	})

	// Return all centroids sorted by distance (or limit to k if you want)
	result := make([]int, int(math.Min(float64(k), float64(len(dists)))))
	for i := 0; i < len(result); i++ {
		result[i] = dists[i].id
	}
	return result, nil
}

func (c *CentroidCache) Size() int {
	return len(c.centroids)
}
