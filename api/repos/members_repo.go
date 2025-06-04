package repos

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"blockbuster/api/constants"
	"blockbuster/api/data"
)

const membersTableName = "BluckBoster_members"

type MemberRepo struct {
	client    *dynamodb.Client
	tableName string
	MovieRepo MovieRepo
}

func NewMembersRepo(client *dynamodb.Client) MemberRepo {
	return MemberRepo{
		client:    client,
		tableName: membersTableName,
		MovieRepo: NewMovieRepo(client),
	}
}

func (r MemberRepo) GetMemberByUsername(ctx context.Context, username string, cartOnly bool) (bool, data.Member, error) {
	member := data.Member{}
	var projectionExpr *string
	var exprAttrNames map[string]string
	if cartOnly {
		expr := "username, cart, checked_out, #t"
		projectionExpr = &expr
		exprAttrNames = map[string]string{
			"#t": constants.TYPE,
		}
	}
	input := &dynamodb.GetItemInput{
		Key:                      map[string]types.AttributeValue{constants.USERNAME: &types.AttributeValueMemberS{Value: username}},
		TableName:                &r.tableName,
		ProjectionExpression:     projectionExpr,
		ExpressionAttributeNames: exprAttrNames,
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Err fetching user from cloud: %s\n", err)
		return false, member, err
	}
	if result.Item == nil {
		log.Printf("Failed to get user %s\n", username)
		return false, member, nil
	}

	err = attributevalue.UnmarshalMap(result.Item, &member)
	if err != nil {
		log.Fatalf("Failed to unmarshall data %s\n", err)
		return false, member, err
	}
	return true, member, nil
}

func (r MemberRepo) GetCartMovies(ctx context.Context, username string) ([]data.CartMovie, error) {
	_, user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		log.Printf("Err in fetching cart movie ids for %s\n", username)
		return nil, err
	}

	movies := make([]data.CartMovie, 0)
	if 0 < len(user.Cart) {
		_, movies, err = r.MovieRepo.GetMoviesByID(ctx, user.Cart, constants.CART)
		if err != nil {
			log.Printf("Failed to get movies in cart. %s\n", err)
			return nil, err
		}
	}
	return movies, nil
}

func (r MemberRepo) ModifyCart(ctx context.Context, username, movieID, updateKey string, checkingOut bool) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal data %s", err)
	}

	updateExpr := fmt.Sprintf("%s cart :cart", updateKey)
	expressionAttrs := map[string]types.AttributeValue{
		":cart": &types.AttributeValueMemberSS{
			Value: []string{movieID},
		},
	}
	if checkingOut {
		updateExpr = fmt.Sprintf("%s ADD checked_out :checked_out", updateExpr)
		expressionAttrs[":checked_out"] = &types.AttributeValueMemberSS{Value: []string{movieID}}
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: expressionAttrs,
		ReturnValues:              types.ReturnValueUpdatedNew,
		UpdateExpression:          &updateExpr,
	}
	return r.updateMember(ctx, username, updateInput)
}

func (r MemberRepo) updateMember(ctx context.Context, username string, updateInput *dynamodb.UpdateItemInput) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	response, err := r.client.UpdateItem(ctx, updateInput)
	if err != nil {
		log.Printf("Failed to update %s\n%s\n", username, err)
		return false, response, err
	}
	return true, response, nil
}

func (r MemberRepo) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	rented, messages := 0, make([]string, 0)
	found, user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil || !found {
		return nil, rented, fmt.Errorf("failed to retrieve user from cloud. UserFound=%v err=%s", found, err)
	}
	if data.MemberTypes[user.Type] < len(movieIDs)+len(user.Checkedout) {
		return nil, rented, nil
	}

	movies, _, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		return nil, rented, fmt.Errorf("failed to retrieve movies from cloud %s", err)
	}

	for _, movie := range movies {
		if movie.Inventory < 0 {
			messages = append(messages, fmt.Sprintf("%s is out of stock and could not be rented", movie.Title))
			continue
		}

		contains, _ := data.SliceContains(user.Checkedout, movie.ID)
		if contains {
			messages = append(messages, fmt.Sprintf("%s is currently checked out by %s", movie.Title, user.Username))
			continue
		}
		contains, _ = data.SliceContains(user.Cart, movie.ID)
		if !contains {
			messages = append(messages, fmt.Sprintf("%s is not in %s cart", movie.Title, user.Username))
			continue
		}

		success, err := r.checkoutMovie(ctx, user, movie)
		if err != nil {
			messages = append(messages, fmt.Sprintf("Failed to rent %s", movie.ID))
			continue
		}
		if success {
			rented++
		}
	}
	return messages, rented, nil
}

func (r MemberRepo) checkoutMovie(ctx context.Context, user data.Member, movie data.Movie) (bool, error) {
	success, err := r.MovieRepo.Rent(ctx, movie)
	if err != nil {
		return false, fmt.Errorf("err checking %s\n%s", movie.Title, err)
	}
	if !success {
		return false, fmt.Errorf("failed to checkout %s", movie.Title)
	}

	modified, _, err := r.ModifyCart(ctx, user.Username, movie.ID, constants.DELETE, constants.CHECKOUT)
	if !modified || err != nil {
		r.MovieRepo.Return(ctx, movie)
		return false, fmt.Errorf("failed to remove %s from %s cart\n%s", movie.Title, user.Username, err)
	}
	return true, nil
}

func (r MemberRepo) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, 0, err
	}

	movies, _, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		log.Print("Err returning movies. Failed to fetch movies from cloud")
		return nil, 0, err
	}

	messages, returned := make([]string, 0), 0
	for _, movie := range movies {
		updateExpr := "DELETE checked_out :checked_out ADD rented :rented"
		expressionAttrs := map[string]types.AttributeValue{
			":checked_out": &types.AttributeValueMemberSS{
				Value: []string{movie.ID},
			},
			":rented": &types.AttributeValueMemberSS{
				Value: []string{movie.ID},
			},
		}
		updateInput := &dynamodb.UpdateItemInput{
			TableName:                 &r.tableName,
			Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
			ExpressionAttributeValues: expressionAttrs,
			ReturnValues:              types.ReturnValueUpdatedNew,
			UpdateExpression:          &updateExpr,
		}
		ok, response, err := r.updateMember(ctx, username, updateInput)
		if err != nil || !ok {
			msg := fmt.Sprintf("Failed to return movie %s\nResponse: %v\nErr: %s\n", movie.ID, response, err)
			log.Print(msg)
			messages = append(messages, msg)
			continue
		}

		ok, err = r.MovieRepo.Return(ctx, movie)
		if err != nil || !ok {
			msg := fmt.Sprintf("Failed to update movie table %s\nResponse: %v\nErr: %s\n", movie.ID, response, err)
			log.Print(msg)
			messages = append(messages, msg)
			continue
		}
		returned++
	}
	return messages, returned, nil
}
