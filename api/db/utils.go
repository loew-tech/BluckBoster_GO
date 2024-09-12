package db

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MarshallUsernameAndMovieID(username, movieID string) (types.AttributeValue, types.AttributeValue, error, error) {
	name, errN := attributevalue.Marshal(username)
	if errN != nil {
		log.Printf("Cannot marshall %s\n", username)
	}
	movie, errM := attributevalue.Marshal(movieID)
	if errM != nil {
		log.Printf("Cannot marshall %s\n", movieID)
	}

	return name, movie, errN, errM
}
