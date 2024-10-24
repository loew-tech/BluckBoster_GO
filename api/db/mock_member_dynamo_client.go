package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MockDynamoMemberClient struct{}

func (mock MockDynamoMemberClient) BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	return &dynamodb.BatchGetItemOutput{}, nil
}

func (mock MockDynamoMemberClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput,
	optFns ...func(*dynamodb.Options)) (
	*dynamodb.GetItemOutput,
	error,
) {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"username":   &types.AttributeValueMemberS{Value: "sea_captain"},
			"first_name": &types.AttributeValueMemberS{Value: "Sea"},
			"last_name":  &types.AttributeValueMemberS{Value: "Captain"},
			"type":       &types.AttributeValueMemberS{Value: "advance"},
			"cart":       &types.AttributeValueMemberSS{Value: TestMovieIDs},
		},
	}, nil
}

func (mock MockDynamoMemberClient) Scan(ctx context.Context, params *dynamodb.ScanInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{}, nil
}

func (mock MockDynamoMemberClient) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, nil
}
