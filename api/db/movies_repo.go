package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TableName = "BluckBoster_movies"
const ID = "id"

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
		tableName: TableName,
	}
}


func (r MovieRepo) GetAllMovies() []Movie {
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
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		movies = append(movies, movie)
	}

	return movies
}

func (r MovieRepo) QueryMovieByID(id string) Movie {

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
	if result.Items == nil {
		log.Println("Could not find item")
	}

	println("Result=", result)
	movie := Movie{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &movie)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
	}
	return movie
}
