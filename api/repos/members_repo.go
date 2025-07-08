package repos

import (
	"context"
	"errors"
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
	MovieRepo *MovieRepo
}

func NewMembersRepo(client *dynamodb.Client, movieRepo *MovieRepo) *MemberRepo {
	return &MemberRepo{
		client:    client,
		tableName: membersTableName,
		MovieRepo: movieRepo,
	}
}

func (r *MemberRepo) GetMemberByUsername(ctx context.Context, username string, cartOnly bool) (data.Member, error) {
	member := data.Member{}

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			constants.USERNAME: &types.AttributeValueMemberS{Value: username},
		},
		TableName: &r.tableName,
	}

	if cartOnly {
		expr := "username, cart, checked_out, #t"
		input.ProjectionExpression = &expr
		input.ExpressionAttributeNames = map[string]string{"#t": constants.TYPE}
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return member, logError("fetching user from cloud", err)
	}
	if result.Item == nil {
		return member, logError(fmt.Sprintf("user %s not found", username), nil)
	}

	err = attributevalue.UnmarshalMap(result.Item, &member)
	if err != nil {
		return member, logError("unmarshalling user data", err)
	}
	return member, nil
}

func (r *MemberRepo) GetCartMovies(ctx context.Context, username string) ([]data.Movie, error) {
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		return nil, logError("fetching cart movie IDs", err)
	}
	if len(user.Cart) == 0 {
		return []data.Movie{}, nil
	}
	return r.MovieRepo.GetMoviesByID(ctx, user.Cart, constants.CART)
}

func (r *MemberRepo) ModifyCart(ctx context.Context, username, movieID, updateKey string, checkingOut bool) (bool, *dynamodb.UpdateItemOutput, error) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		return false, nil, logError("marshalling username", err)
	}

	expr, attrs := buildCartUpdateExpr(movieID, updateKey, checkingOut)
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: attrs,
		ReturnValues:              types.ReturnValueUpdatedNew,
		UpdateExpression:          &expr,
	}
	return r.updateMember(ctx, username, updateInput)
}

func (r *MemberRepo) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		return nil, 0, logError("retrieving user", err)
	}

	if data.MemberTypes[user.Type] < len(movieIDs)+len(user.Checkedout) {
		return []string{"member limit exceeded"}, 0, nil
	}

	movies, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		return nil, 0, logError("retrieving movies", err)
	}

	return r.performCheckout(ctx, user, movies)
}

func (r *MemberRepo) performCheckout(ctx context.Context, user data.Member, movies []data.Movie) ([]string, int, error) {
	var rented int
	var messages []string

	for _, movie := range movies {
		if movie.Inventory < 0 {
			messages = append(messages, fmt.Sprintf("%s is out of stock", movie.Title))
			continue
		}
		if contains(user.Checkedout, movie.ID) {
			messages = append(messages, fmt.Sprintf("%s is already checked out", movie.Title))
			continue
		}
		if !contains(user.Cart, movie.ID) {
			messages = append(messages, fmt.Sprintf("%s is not in cart", movie.Title))
			continue
		}
		if err := r.checkoutMovie(ctx, user, movie); err != nil {
			messages = append(messages, err.Error())
			continue
		}
		rented++
	}
	return messages, rented, nil
}

func (r *MemberRepo) checkoutMovie(ctx context.Context, user data.Member, movie data.Movie) error {
	ok, err := r.MovieRepo.Rent(ctx, movie)
	if err != nil || !ok {
		return logError(fmt.Sprintf("renting %s", movie.Title), err)
	}

	ok, _, err = r.ModifyCart(ctx, user.Username, movie.ID, constants.DELETE, constants.CHECKOUT)
	if err != nil || !ok {
		r.MovieRepo.Return(ctx, movie)
		return logError(fmt.Sprintf("removing %s from cart", movie.Title), err)
	}
	return nil
}

func (r *MemberRepo) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	var messages []string
	var returned int

	movies, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		return nil, 0, logError("fetching movies for return", err)
	}

	name, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, 0, logError("marshalling username", err)
	}

	for _, movie := range movies {
		updateInput, err := r.getReturnInput(movie, name)
		if err != nil {
			messages = append(messages, logError(fmt.Sprintf("preparing return for %s", movie.Title), err).Error())
			continue
		}

		ok, _, err := r.updateMember(ctx, username, updateInput)
		if err != nil || !ok {
			messages = append(messages, logError(fmt.Sprintf("returning %s", movie.Title), err).Error())
			continue
		}

		ok, err = r.MovieRepo.Return(ctx, movie)
		if err != nil || !ok {
			messages = append(messages, logError(fmt.Sprintf("updating inventory for %s", movie.Title), err).Error())
			continue
		}
		returned++
	}
	return messages, returned, nil
}

func (r *MemberRepo) SetMemberAPIChoice(ctx context.Context, username, apiChoice string) error {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		errWrap := fmt.Errorf("failed to marshal username %s: %w", username, err)
		log.Println(errWrap)
		return errWrap
	}
	updateExpr := "SET api_choice = :api_choice"
	expressionAttrs := map[string]types.AttributeValue{
		":api_choice": &types.AttributeValueMemberS{
			Value: apiChoice,
		},
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: expressionAttrs,
		UpdateExpression:          &updateExpr,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}
	_, err = r.client.UpdateItem(ctx, updateInput)
	if err != nil {
		errWrap := fmt.Errorf("failed to update member %s api choice: %w", username, err)
		log.Println(errWrap)
		return errWrap
	}
	return nil
}

func (r *MemberRepo) updateMember(ctx context.Context, username string, input *dynamodb.UpdateItemInput) (bool, *dynamodb.UpdateItemOutput, error) {
	response, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return false, response, logError("updating member", err)
	}
	return true, response, nil
}

func (r *MemberRepo) getReturnInput(movie data.Movie, name types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
	expr := "DELETE checked_out :checked_out ADD rented :rented"
	attrs := map[string]types.AttributeValue{
		":checked_out": &types.AttributeValueMemberSS{Value: []string{movie.ID}},
		":rented":      &types.AttributeValueMemberSS{Value: []string{movie.ID}},
	}
	return &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: attrs,
		UpdateExpression:          &expr,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}, nil
}

func buildCartUpdateExpr(movieID, updateKey string, checkingOut bool) (string, map[string]types.AttributeValue) {
	expr := fmt.Sprintf("%s cart :cart", updateKey)
	attrs := map[string]types.AttributeValue{
		":cart": &types.AttributeValueMemberSS{Value: []string{movieID}},
	}
	if checkingOut {
		expr += " ADD checked_out :checked_out"
		attrs[":checked_out"] = &types.AttributeValueMemberSS{Value: []string{movieID}}
	}
	return expr, attrs
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func logError(msg string, err error) error {
	if err == nil {
		err = errors.New(msg)
	} else {
		err = fmt.Errorf("%s: %w", msg, err)
	}
	log.Println(err)
	return err
}
