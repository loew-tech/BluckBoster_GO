package api_cache

import (
	"fmt"
	"math/rand"

	"blockbuster/api/utils"
)

type CentroidsToMoviesCache struct {
	CentroidToMovieIDs map[int][]string
}

func (ctm *CentroidsToMoviesCache) GetMovieIDsByCentroid(centroid int) ([]string, error) {
	if movieIDs, ok := ctm.CentroidToMovieIDs[centroid]; ok {
		return movieIDs, nil
	}
	return nil, utils.LogError(fmt.Sprintf("cannot find movies for centroid id %v", centroid), nil)
}

func (ctm *CentroidsToMoviesCache) GetRandomMovieFromCentroid(centroid int) (string, error) {
	if movieIDs, ok := ctm.CentroidToMovieIDs[centroid]; ok {
		return movieIDs[rand.Intn(len(movieIDs))], nil
	}
	return "", utils.LogError(fmt.Sprintf("failed to retrieve random movie from centroid %d", centroid), nil)
}
