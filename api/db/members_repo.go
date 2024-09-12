package db

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const membersTableName = "BluckBoster_members"

type MemberRepo struct {
	client    dynamodb.Client
	tableName string
}

func NewMembersRepo() MemberRepo {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("FAILED TO INSTANTIATE MemberRepo")
	}

	return MemberRepo{
		client:    *dynamodb.NewFromConfig(config),
		tableName: membersTableName,
	}
}

func (r MemberRepo) GetMemberByUsername(username string) (bool, Member, error) {
	// @TODO: never actually return error (it's always nil)
	keyEx := expression.Key(USERNAME).Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Fatalf("Failed to build query expression")
	}
	queryInput := &dynamodb.QueryInput{
		TableName:                 &r.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := r.client.Query(context.TODO(), queryInput)
	if err != nil {
		log.Fatalf("Got error calling svc.Query %s\n", err)
	}
	if len(result.Items) == 0 {
		log.Printf("Could not find member with username: %s\n", username)
		return false, Member{}, nil
	}
	member := Member{}
	err = attributevalue.UnmarshalMap(result.Items[0], &member)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
	}
	return true, member, nil
}

func (r MemberRepo) cartContains(username, movieID string) (bool, int, error) {
	ids, err := r.GetCartIDs(username)
	if err != nil {
		log.Printf("Failed checking to see if %s is in %s cart\n", movieID, username)
		return false, -1, err
	}
	for i, id := range ids {
		if id == movieID {
			return true, i, nil
		}
	}
	return false, -1, nil
}

func (r MemberRepo) GetCartIDs(username string) ([]string, error) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		log.Printf("Failed to marshal %s\n", username)
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key:             map[string]types.AttributeValue{USERNAME: name},
		TableName:       &r.tableName,
		AttributesToGet: []string{CART},
	}

	response, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, err
	}

	cart := Cart{}
	err = attributevalue.UnmarshalMap(response.Item, &cart)
	if err != nil {
		log.Printf("Err unmarshalling movies from: %s\n", err)
		return nil, err
	}
	return cart.Cart, err
}

func (r MemberRepo) AddToCart(username, movieID string) (
	bool, *dynamodb.UpdateItemOutput, error,
) {

	name, movie, errN, errM := MarshallUsernameAndMovieID(username, movieID)
	if errN != nil || errM != nil {
		return false, nil, fmt.Errorf("failed to marshal data errName=%s errMovie=%s", errN, errM)
	}

	found, _, err := r.cartContains(username, movieID)
	if err != nil || found {
		return false, nil, err
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key:       map[string]types.AttributeValue{USERNAME: name},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cart": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{movie},
			},
			":empty_cart": &types.AttributeValueMemberL{
				Value: make([]types.AttributeValue, 0),
			},
		},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String("SET cart = list_append(if_not_exists(cart, :empty_cart), :cart)"),
	}

	response, err := r.client.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		log.Printf("Failed to add movie %s to %s cart\n %s\n", movieID, username, err)
		return false, response, err
	}
	return true, response, nil
}

func (r MemberRepo) RemoveFromCart(username, movieID string) (
	bool, *dynamodb.UpdateItemOutput, error,
) {

	name, _, errN, errM := MarshallUsernameAndMovieID(username, movieID)
	if errN != nil || errM != nil {
		return false, nil, fmt.Errorf("failed to marshal data errName=%s errMovie=%s", errN, errM)
	}

	found, index, err := r.cartContains(username, movieID)
	if err != nil || !found {
		return false, nil, err
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:        aws.String(r.tableName),
		Key:              map[string]types.AttributeValue{USERNAME: name},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String(fmt.Sprintf("REMOVE cart[%v]", index)),
	}

	response, err := r.client.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		log.Printf("Failed to add movie %s to %s cart\n %s\n", movieID, username, err)
		return false, response, err
	}
	return true, response, nil
}
