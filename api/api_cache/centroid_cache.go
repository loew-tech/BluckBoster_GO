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
	centroid, ok := c.centroids[centroidID]
	if !ok {
		return data.MovieMetrics{}, fmt.Errorf("centroid %d not found", centroidID)
	}
	return centroid, nil
}

func (c *CentroidCache) GetKNearestCentroidsFromMood(mood data.MovieMetrics, k int) ([]int, error) {
	if k <= 0 {
		return nil, utils.LogError("k must be greater than 0", nil)
	}
	type centroidDist struct {
		id   int
		dist float64
	}

	// Euclidean distance between two MovieMetrics
	distance := func(a, b data.MovieMetrics) float64 {
		sum := 0.0
		sum += (a.Acting - b.Acting) * (a.Acting - b.Acting)
		sum += (a.Action - b.Action) * (a.Action - b.Action)
		sum += (a.Cinematography - b.Cinematography) * (a.Cinematography - b.Cinematography)
		sum += (a.Comedy - b.Comedy) * (a.Comedy - b.Comedy)
		sum += (a.Directing - b.Directing) * (a.Directing - b.Directing)
		sum += (a.Drama - b.Drama) * (a.Drama - b.Drama)
		sum += (a.Fantasy - b.Fantasy) * (a.Fantasy - b.Fantasy)
		sum += (a.Horror - b.Horror) * (a.Horror - b.Horror)
		sum += (a.Romance - b.Romance) * (a.Romance - b.Romance)
		sum += (a.StoryTelling - b.StoryTelling) * (a.StoryTelling - b.StoryTelling)
		sum += (a.Suspense - b.Suspense) * (a.Suspense - b.Suspense)
		sum += (a.Writing - b.Writing) * (a.Writing - b.Writing)
		return sum
	}

	var dists []centroidDist
	for id, metrics := range c.centroids {
		d := distance(mood, metrics)
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
