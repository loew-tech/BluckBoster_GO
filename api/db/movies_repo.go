package db

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const movieTableName = "BluckBoster_movies"

type MovieRepo struct {
	client    dynamodb.Client
	tableName string
}

func NewMovieRepo() MovieRepo {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("FAILED TO INSTANTIATE MovieRepo")
	}

	return MovieRepo{
		client:    *dynamodb.NewFromConfig(config),
		tableName: movieTableName,
	}
}

func (r MovieRepo) GetAllMovies() ([]Movie, error) {

	// @TODO: remove debug return
	return TestMovies, nil

	params := &dynamodb.ScanInput{
		TableName: &r.tableName,
	}
	result, err := r.client.Scan(context.TODO(), params)
	if err != nil {
		log.Printf("Query API call failed: %s\n", err)
		return nil, err
	}

	movies := make([]Movie, 0)
	for _, i := range result.Items {
		movie := Movie{}
		err = attributevalue.UnmarshalMap(i, &movie)

		if err != nil {
			log.Printf("Got error unmarshalling: %s\n", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// @TODO: update to sdk v2 when needed
// func (r MovieRepo) QueryMovieByID(id string) (Movie, error) {
// 	queryInput := &dynamodb.QueryInput{
// 		TableName: aws.String(r.tableName),
// 		KeyConditions: map[string]*dynamodb.Condition{
// 			ID: {
// 				ComparisonOperator: aws.String("EQ"),
// 				AttributeValueList: []*dynamodb.AttributeValue{
// 					{
// 						S: aws.String(id),
// 					},
// 				},
// 			},
// 		},
// 	}

// 	result, err := r.svc.Query(queryInput)
// 	if err != nil {
// 		log.Printf("Query API call failed: %s\n", err)
// 		return Movie{}, err
// 	}
// 	if len(result.Items) == 0 {
// 		log.Printf("Could not find movie with id: %s\n", id)
// 		return Movie{}, nil
// 	}

// 	movie := Movie{}
// 	err = dynamodbattribute.UnmarshalMap(result.Items[0], &movie)
// 	if err != nil {
// 		log.Fatalf("Failed to unmarshall data %s\n", err)
// 	}
// 	return movie, nil
// }

func (r MovieRepo) GetMoviesByID(movieIDs []string) ([]CartMovie, error) {

	var keys []map[string]types.AttributeValue
	for _, mid := range movieIDs {
		keys = append(keys, map[string]types.AttributeValue{ID: &types.AttributeValueMemberS{Value: mid}})
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.tableName: {
				Keys:                 keys,
				ProjectionExpression: aws.String("id, title, inventory"),
			},
		},
	}

	result, err := r.client.BatchGetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, err
	}

	movies := make([]CartMovie, 0)
	for _, v := range result.Responses {
		for _, m := range v {
			movie := CartMovie{}
			if err = attributevalue.UnmarshalMap(m, &movie); err != nil {
				log.Printf("Got error unmarshalling: %s", err)
				continue
			}
			movies = append(movies, movie)
		}
	}
	return movies, nil
}

func (r MovieRepo) RentMovie(movie Movie) (bool, error) {

	mid, err := attributevalue.Marshal(movie.ID)
	if err != nil {
		return false, fmt.Errorf("failed to marshal data %s", err)
	}

	updateInventoryAndRentedInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(movieTableName),
		Key:       map[string]types.AttributeValue{ID: mid},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inventory": &types.AttributeValueMemberN{
				Value: strconv.Itoa(movie.Inventory - 1),
			},
			":rented": &types.AttributeValueMemberN{
				Value: strconv.Itoa(movie.Rented + 1),
			},
		},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String("set inventory = :inventory, rented = :rented"),
	}

	_, err = r.client.UpdateItem(context.TODO(), updateInventoryAndRentedInput)
	if err != nil {
		log.Printf("Failed to checkout %s\n%s\n", movie.ID, err)
		return false, err
	}
	return true, nil
}
