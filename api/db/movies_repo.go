package db

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const movieTableName = "BluckBoster_movies"

type MovieRepo struct {
	svc       dynamodb.DynamoDB
	tableName string
}

func NewMovieRepo() MovieRepo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return MovieRepo{
		svc:       *dynamodb.New(sess),
		tableName: movieTableName,
	}
}

func (r MovieRepo) GetAllMovies() ([]Movie, error) {
	params := &dynamodb.ScanInput{
		TableName: &r.tableName,
	}
	result, err := r.svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s\n", err)
		return nil, err
	}

	movies := make([]Movie, 0)
	for _, i := range result.Items {
		movie := Movie{}
		err = dynamodbattribute.UnmarshalMap(i, &movie)

		if err != nil {
			log.Printf("Got error unmarshalling: %s\n", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r MovieRepo) QueryMovieByID(id string) (Movie, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			ID: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
		},
	}

	result, err := r.svc.Query(queryInput)
	if err != nil {
		log.Printf("Query API call failed: %s\n", err)
		return Movie{}, err
	}
	if len(result.Items) == 0 {
		log.Printf("Could not find movie with id: %s\n", id)
		return Movie{}, nil
	}

	movie := Movie{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &movie)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
	}
	return movie, nil
}

func (r MovieRepo) GetMoviesByID(movieIDs []string) ([]MovieIdAndTitle, []error, error) {
	var keys []map[string]*dynamodb.AttributeValue
	for _, mid := range movieIDs {
		m := map[string]*dynamodb.AttributeValue{
			ID: {
				S: aws.String(mid),
			},
		}
		keys = append(keys, m)
	}

	proj := expression.NamesList(expression.Name(ID), expression.Name(TITLE))
	expr, _ := expression.NewBuilder().WithProjection(proj).Build()
	fmt.Println("$$ expr.Projection=", expr.Projection())
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			r.tableName: {
				Keys:                 keys,
				ProjectionExpression: aws.String("id, title"),
			},
		},
	}

	movies, errors := make([]MovieIdAndTitle, 0), make([]error, 0)
	result, err := r.svc.BatchGetItem(input)

	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, nil, err
	}

	for _, v := range result.Responses {
		for _, m := range v {
			movie := MovieIdAndTitle{}
			err = dynamodbattribute.UnmarshalMap(m, &movie)

			if err != nil {
				log.Printf("Got error unmarshalling: %s", err)
				errors = append(errors, err)
			}
			movies = append(movies, movie)

		}
	}
	return movies, errors, nil
}
