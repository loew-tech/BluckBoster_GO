package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

// func (r MovieRepo) GetMoviesByID(movieIDs []string) ([]MovieIdAndTitle, []error, error) {
// 	var keys []map[string]*dynamodb.AttributeValue
// 	for _, mid := range movieIDs {
// 		m := map[string]*dynamodb.AttributeValue{
// 			ID: {
// 				S: aws.String(mid),
// 			},
// 		}
// 		keys = append(keys, m)
// 	}

// 	input := &dynamodb.BatchGetItemInput{
// 		RequestItems: map[string]*dynamodb.KeysAndAttributes{
// 			r.tableName: {
// 				Keys:                 keys,
// 				ProjectionExpression: aws.String("id, title"),
// 			},
// 		},
// 	}

// 	result, err := r.svc.BatchGetItem(input)
// 	if err != nil {
// 		log.Printf("Err fetching movies from cloud: %s\n", err)
// 		return nil, nil, err
// 	}

// 	movies, errors := make([]MovieIdAndTitle, 0), make([]error, 0)
// 	for _, v := range result.Responses {
// 		for _, m := range v {
// 			movie := MovieIdAndTitle{}
// 			err = dynamodbattribute.UnmarshalMap(m, &movie)

// 			if err != nil {
// 				log.Printf("Got error unmarshalling: %s", err)
// 				errors = append(errors, err)
// 			}
// 			movies = append(movies, movie)

// 		}
// 	}
// 	return movies, errors, nil
// }
