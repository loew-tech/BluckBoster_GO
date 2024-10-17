package endpoints

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetDynamoClient() *dynamodb.Client {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("FAILED TO INSTANTIATE MemberRepo", err)
	}
	return dynamodb.NewFromConfig(config)
}
