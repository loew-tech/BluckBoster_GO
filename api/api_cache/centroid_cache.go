package centroidcache

import "blockbuster/api/data"

type CentroidCache struct {
	centroids map[int][]data.MovieMetrics
}

func (c *CentroidCache) GetMetricsByCentroid(centroid int) ([]data.MovieMetrics, bool) {
	metrics, exists := c.centroids[centroid]
	return metrics, exists
}

func (c *CentroidCache) GetCentroidsFromMood(mood data.MovieMetrics) []int {
	// @TODO: implement
	return []int{}
}

func (c *CentroidCache) SetCentroids(centroids map[int][]data.MovieMetrics) {
	c.centroids = centroids
}
