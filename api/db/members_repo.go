package db

import (
	"blockbuster/api/utils"
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
	// @TODO: switch from Query to GetItem
	member := Member{}
	keyEx := expression.Key(USERNAME).Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Fatalf("Failed to build query expression")
		return false, member, err
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
		return false, member, err
	}
	if len(result.Items) == 0 {
		log.Printf("Could not find member with username: %s\n", username)
		return false, member, nil
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &member)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
		return true, member, err
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

func (r MemberRepo) Checkout(username string, movieIDs []string) (
	bool, *dynamodb.UpdateItemOutput, []string, int, error,
) {
	rented := 0
	messages := make([]string, 0)

	found, user, err := r.GetMemberByUsername(username)
	if err != nil || !found {
		return false, nil, messages, rented, fmt.Errorf("failed to retrieve user from cloud. UserFound=%v err=%s", found, err)
	}

	movies, err := r.MovieRepo.GetMoviesByID(movieIDs)
	if err != nil {
		return false, nil, messages, rented, fmt.Errorf("failed to retrieve movies from cloud %s", err)
	}

	for _, movie := range movies {
		if movie.Inventory < 0 {
			messages = append(messages, fmt.Sprintf("%s is out of stock and could not be rented", movie.Title))
		}
		contains, _ := utils.SliceContains(user.CheckedOut, movie.ID)
		if contains {
			messages = append(messages, fmt.Sprintf("%s is currently checked out by %s", movie.Title, user.Username))
		}

	}
	return false, nil, messages, rented, nil
}

func (r MemberRepo) CheckoutMovie(user Member, movie Movie) (bool, error) {
	success, err := r.MovieRepo.RentMovie(movie)
	if !success {
		return false, fmt.Errorf("err checking %s\n%s", movie.Title, err)
	}
	if !success {
		return false, fmt.Errorf("failed to checkout %s", movie.Title)
	}

	name, err := attributevalue.Marshal(user.Username)
	if err != nil {
		return false, fmt.Errorf("failed to marshal data %s", err)
	}

	updateCheckedOutInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(membersTableName),
		Key:       map[string]types.AttributeValue{USERNAME: name},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":checked_out": &types.AttributeValueMemberSS{
				Value: []string{movie.ID},
			},
		},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String("ADD checked_out :checked_out"),
	}

	_, err = r.client.UpdateItem(context.TODO(), updateCheckedOutInput)
	if err != nil {
		log.Printf("Failed for user %s to checkout %s\n%s\n", movie.ID, user.Username, err)
		return false, err
	}

	modified, _, err := r.ModifyCart(user.Username, movie.ID, DELETE)
	if !modified || err != nil {
		return false, fmt.Errorf("failed to remove %s from %s cart\n%s", movie.Title, user.Username, err)
	}
	return true, nil
}
