package centroidcache

import (
	"blockbuster/api/data"
	"fmt"
)

type CentroidCache struct {
	centroids map[int]data.MovieMetrics
}

func (c *CentroidCache) GetMetricsByCentroid(centroid int) (data.MovieMetrics, error) {
	if centroid < 0 || centroid >= len(c.centroids) {
		return data.MovieMetrics{}, fmt.Errorf("centroid %d not found", centroid)
	}
	return c.centroids[centroid], nil
}

func (c *CentroidCache) GetCentroidsFromMood(mood data.MovieMetrics) []int {
	// @TODO: implement
	return []int{}
}

func (c *CentroidCache) SetCentroids(centroids map[int]data.MovieMetrics) {
	c.centroids = centroids
}
