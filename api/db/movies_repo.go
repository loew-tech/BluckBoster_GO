package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const movieTableName   = "BluckBoster_movies"

type MovieRepo struct {
	svc dynamodb.DynamoDB
	tableName string
}

func NewMovieRepo() MovieRepo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return MovieRepo{
		svc: *dynamodb.New(sess),
		tableName: movieTableName,
	}
}

func (r MovieRepo) GetAllMovies() ([]Movie, error) {
	// @TODO: err is always nil
	params := &dynamodb.ScanInput{
		TableName: &r.tableName,
	}
	result, err := r.svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s\n", err)
	}

	movies := make([]Movie, 0)
	for _, i := range result.Items {
		movie := Movie{}
		err = dynamodbattribute.UnmarshalMap(i, &movie)

		if err != nil {
			log.Panicf("Got error unmarshalling: %s", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r MovieRepo) QueryMovieByID(id string) (Movie, error) {
	// @TODO: never actually return error (it's always nil)
	queryInput  := &dynamodb.QueryInput{
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
		log.Fatalf("Got error calling svc.Query %s\n", err)
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
