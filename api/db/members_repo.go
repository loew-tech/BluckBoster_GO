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
	MovieRepo MovieRepo
}

func NewMembersRepo() MemberRepo {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("FAILED TO INSTANTIATE MemberRepo", err)
	}

	return MemberRepo{
		client:    *dynamodb.NewFromConfig(config),
		tableName: membersTableName,
		MovieRepo: NewMovieRepo(),
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

func (r MemberRepo) GetCartMovies(username string) ([]CartMovie, error) {
	ids, err := r.GetCartIDs(username)
	if err != nil {
		log.Printf("Err in fetching cart movie ids for %s\n", username)
		return nil, err
	}
	movies, err := r.MovieRepo.GetMoviesByID(ids)
	if err != nil {
		log.Printf("Failed to get cart. %s\n", err)
		return nil, err
	}
	return movies, nil
}

func (r MemberRepo) ModifyCart(username, movieID, updateKey string) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal data %s", err)
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key:       map[string]types.AttributeValue{USERNAME: name},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cart": &types.AttributeValueMemberSS{
				Value: []string{movieID},
			},
		},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String(fmt.Sprintf("%s cart :cart", updateKey)),
	}

	response, err := r.client.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		log.Printf("Failed to add movie %s to %s cart\n %s\n", movieID, username, err)
		return false, response, err
	}
	return true, response, nil
}
