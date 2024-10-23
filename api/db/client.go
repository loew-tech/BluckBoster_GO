package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBAPI interface {
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)

	GetItem(ctx context.Context, params *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)

	Scan(ctx context.Context, params *dynamodb.ScanInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)

	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}
