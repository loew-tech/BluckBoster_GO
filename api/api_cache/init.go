package api_cache

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/utils"
)

var (
	initCentroidsCacheOnce         sync.Once
	centroidCache                  *CentroidCache
	initCentroidsToMoviesCacheOnce sync.Once
	centroidToMoviesCache          *CentroidsToMoviesCache
)

func GetDynamoClientCentroidCache() *CentroidCache {
	initCentroidsCacheOnce.Do(func() {
		centroidCache = &CentroidCache{
			centroids: getCentroidsFromDynamo(),
		}
	})
	return centroidCache
}

func getCentroidsFromDynamo() map[int]data.MovieMetrics {
	centroidTableName, client := "BluckBoster_centroids", utils.GetDynamoClient()
	centroidItems, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &centroidTableName,
	})
	if err != nil {
		// @TODO: handle limited functionality from failed cache population
		return nil
	}
	var centroids []data.MovieMetrics
	err = attributevalue.UnmarshalListOfMaps(centroidItems.Items, &centroids)
	if err != nil {
		// @TODO: handle limited functionality from failed cache population
		return nil
	}

	centroidsMap := make(map[int]data.MovieMetrics)
	for _, centroid := range centroids {
		centroidsMap[centroid.ID] = centroid
	}
	return centroidsMap
}

func InitCentroidsToMoviesCache(GetMoviesByPage func(
	ctx context.Context, page string, forGraph bool) ([]data.Movie, error),
) *CentroidsToMoviesCache {
	initCentroidsToMoviesCacheOnce.Do(func() {
		centroidsToMovies := make(map[int][]string)
		for _, page := range constants.PAGES {
			movies, err := GetMoviesByPage(context.TODO(), string(page), true)
			if err != nil {
				utils.LogError("failed to get movies for centroid to movies cache", err)
				continue
			}
			for _, movie := range movies {
				centroidsToMovies[movie.Centroid] = append(centroidsToMovies[movie.Centroid], movie.ID)
			}
		}
		centroidToMoviesCache = &CentroidsToMoviesCache{CentroidToMovieIDs: centroidsToMovies}
	})
	return centroidToMoviesCache
}
