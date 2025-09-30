package api_cache

import "blockbuster/api/data"

type CentroidCacheInterface interface {
	GetMetricsByCentroid(centroidID int) (data.MovieMetrics, error)
	GetKNearestCentroidsFromMood(mood data.MovieMetrics, k int) ([]int, error)
	Size() int
}

type CentroidsToMoviesCacheInterface interface {
	GetMovieIDsByCentroid(centroid int) ([]string, error)
	GetRandomMovieFromCentroid(centroid int) (string, error)
}
