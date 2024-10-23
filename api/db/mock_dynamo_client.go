package db

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MockDynamoDBClient struct {
}

func (mock MockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (
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
