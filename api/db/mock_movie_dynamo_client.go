package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MockMovieDynamoClient struct{}

func (mock MockMovieDynamoClient) BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	return &dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]types.AttributeValue{
			"Items": []map[string]types.AttributeValue{
				map[string]types.AttributeValue{
					"cast": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "Kevin Spacey"},
						&types.AttributeValueMemberS{Value: "Russell Crowe"},
						&types.AttributeValueMemberS{Value: "Guy Pearce"},
						&types.AttributeValueMemberS{Value: "James Cromwell"},
					}},
					"director":  &types.AttributeValueMemberS{Value: "Curtis Hanson"},
					"id":        &types.AttributeValueMemberS{Value: "l.a._confidential_1997"},
					"inventory": &types.AttributeValueMemberN{Value: "5"},
					"rented":    &types.AttributeValueMemberN{Value: "0"},
					"rating":    &types.AttributeValueMemberS{Value: "99%"},
					"review":    &types.AttributeValueMemberS{Value: "foo"},
					"synopsis":  &types.AttributeValueMemberS{Value: "bar"},
					"title":     &types.AttributeValueMemberS{Value: "L.A. Confidential"},
					"year":      &types.AttributeValueMemberS{Value: "1997"},
				},
				map[string]types.AttributeValue{
					"cast": &types.AttributeValueMemberL{Value: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "Ingrid Bergman"},
						&types.AttributeValueMemberS{Value: "Russell Crowe"},
						&types.AttributeValueMemberS{Value: "Paul Henreid"},
						&types.AttributeValueMemberS{Value: "Claude Rains"},
					}},
					"director":  &types.AttributeValueMemberS{Value: "Michael Curtiz"},
					"id":        &types.AttributeValueMemberS{Value: "casablanca_1942"},
					"inventory": &types.AttributeValueMemberN{Value: "4"},
					"rented":    &types.AttributeValueMemberN{Value: "0"},
					"rating":    &types.AttributeValueMemberS{Value: "99%"},
					"review":    &types.AttributeValueMemberS{Value: " An undisputed masterpiece and perhaps Hollywood's quintessential statement on love and romance, "},
					"synopsis":  &types.AttributeValueMemberS{Value: "Rick Blaine (Humphrey Bogart), who owns a nightclub in Casablanca, discovers his old flame Ilsa (Ingrid Bergman) is in town..."},
					"title":     &types.AttributeValueMemberS{Value: "Casablanca"},
					"year":      &types.AttributeValueMemberS{Value: "1942"},
				},
			},
		},
	}, nil
}

func (mock MockMovieDynamoClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput,
	optFns ...func(*dynamodb.Options)) (
	*dynamodb.GetItemOutput,
	error,
) {
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"cast": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "Kevin Spacey"},
				&types.AttributeValueMemberS{Value: "Russell Crowe"},
				&types.AttributeValueMemberS{Value: "Guy Pearce"},
				&types.AttributeValueMemberS{Value: "James Cromwell"},
			}},
			"director":  &types.AttributeValueMemberS{Value: "Curtis Hanson"},
			"id":        &types.AttributeValueMemberS{Value: "l.a._confidential_1997"},
			"inventory": &types.AttributeValueMemberN{Value: "5"},
			"rented":    &types.AttributeValueMemberN{Value: "0"},
			"rating":    &types.AttributeValueMemberS{Value: "99%"},
			"review":    &types.AttributeValueMemberS{Value: "foo"},
			"synopsis":  &types.AttributeValueMemberS{Value: "bar"},
			"title":     &types.AttributeValueMemberS{Value: "L.A. Confidential"},
			"year":      &types.AttributeValueMemberS{Value: "1997"},
		},
	}, nil
}

func (mock MockMovieDynamoClient) Scan(ctx context.Context, params *dynamodb.ScanInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			map[string]types.AttributeValue{
				"cast": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "Kevin Spacey"},
					&types.AttributeValueMemberS{Value: "Russell Crowe"},
					&types.AttributeValueMemberS{Value: "Guy Pearce"},
					&types.AttributeValueMemberS{Value: "James Cromwell"},
				}},
				"director":  &types.AttributeValueMemberS{Value: "Curtis Hanson"},
				"id":        &types.AttributeValueMemberS{Value: "l.a._confidential_1997"},
				"inventory": &types.AttributeValueMemberN{Value: "5"},
				"rented":    &types.AttributeValueMemberN{Value: "0"},
				"rating":    &types.AttributeValueMemberS{Value: "99%"},
				"review":    &types.AttributeValueMemberS{Value: "foo"},
				"synopsis":  &types.AttributeValueMemberS{Value: "bar"},
				"title":     &types.AttributeValueMemberS{Value: "L.A. Confidential"},
				"year":      &types.AttributeValueMemberS{Value: "1997"},
			},
			map[string]types.AttributeValue{
				"cast": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "Ingrid Bergman"},
					&types.AttributeValueMemberS{Value: "Russell Crowe"},
					&types.AttributeValueMemberS{Value: "Paul Henreid"},
					&types.AttributeValueMemberS{Value: "Claude Rains"},
				}},
				"director":  &types.AttributeValueMemberS{Value: "Michael Curtiz"},
				"id":        &types.AttributeValueMemberS{Value: "casablanca_1942"},
				"inventory": &types.AttributeValueMemberN{Value: "4"},
				"rented":    &types.AttributeValueMemberN{Value: "0"},
				"rating":    &types.AttributeValueMemberS{Value: "99%"},
				"review":    &types.AttributeValueMemberS{Value: " An undisputed masterpiece and perhaps Hollywood's quintessential statement on love and romance, "},
				"synopsis":  &types.AttributeValueMemberS{Value: "Rick Blaine (Humphrey Bogart), who owns a nightclub in Casablanca, discovers his old flame Ilsa (Ingrid Bergman) is in town..."},
				"title":     &types.AttributeValueMemberS{Value: "Casablanca"},
				"year":      &types.AttributeValueMemberS{Value: "1942"},
			},
		},
	}, nil
}

func (mock MockMovieDynamoClient) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, nil
}
