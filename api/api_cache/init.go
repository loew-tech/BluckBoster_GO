package centroidcache

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"blockbuster/api/data"
	"blockbuster/api/utils"
)

var (
	initOnce      sync.Once
	centroidCache *CentroidCache
)

func NewCentroidCache() *CentroidCache {
	initOnce.Do(func() {
		centroidCache = &CentroidCache{
			centroids: getCentroidCache(),
		}
	})
	return centroidCache
}

func getCentroidCache() map[int]data.MovieMetrics {
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
	}

	centroidsMap := make(map[int]data.MovieMetrics)
	for _, centroid := range centroids {
		centroidsMap[centroid.ID] = centroid
	}
	return centroidsMap
}
