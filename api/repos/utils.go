package repos

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dynamoClient *dynamodb.Client
	once         sync.Once
)

func GetDynamoClient() *dynamodb.Client {
	once.Do(func() {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalln("FAILED TO INSTANTIATE MemberRepo", err)
		}
		dynamoClient = dynamodb.NewFromConfig(cfg)
	})
	return dynamoClient
}
