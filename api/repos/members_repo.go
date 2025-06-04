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
	MovieRepo *MovieRepo
}

func NewMembersRepo(client *dynamodb.Client) *MemberRepo {
	return &MemberRepo{
		client:    client,
		tableName: membersTableName,
		MovieRepo: NewMovieRepo(client),
	}
}

func (r *MemberRepo) GetMemberByUsername(ctx context.Context, username string, cartOnly bool) (data.Member, error) {
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
		errWrap := fmt.Errorf("err fetching user from cloud: %w", err)
		log.Println(errWrap)
		return member, errWrap
	}
	if result.Item == nil {
		errWrap := fmt.Errorf("failed to get user %s", username)
		log.Println(errWrap)
		return member, errWrap
	}

	err = attributevalue.UnmarshalMap(result.Item, &member)
	if err != nil {
		errWrap := fmt.Errorf("failed to unmarshall data %w", err)
		log.Println(errWrap)
		return member, errWrap
	}
	return member, nil
}

func (r *MemberRepo) GetCartMovies(ctx context.Context, username string) ([]data.CartMovie, error) {
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		log.Printf("Err in fetching cart movie ids for %s\n", username)
		return nil, err
	}

	var movies []data.CartMovie
	if 0 < len(user.Cart) {
		_, movies, err = r.MovieRepo.GetMoviesByID(ctx, user.Cart, constants.CART)
		if err != nil {
			errWrap := fmt.Errorf("failed to get movies in cart. %w", err)
			log.Println(errWrap)
			return nil, errWrap
		}
	}
	return movies, nil
}

func (r *MemberRepo) ModifyCart(ctx context.Context, username, movieID, updateKey string, checkingOut bool) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	name, err := attributevalue.Marshal(username)
	if err != nil {
		errWrap := fmt.Errorf("failed to marshal username %s: %w", username, err)
		log.Println(errWrap)
		return false, nil, errWrap
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

func (r *MemberRepo) updateMember(ctx context.Context, username string, updateInput *dynamodb.UpdateItemInput) (
	bool, *dynamodb.UpdateItemOutput, error,
) {
	response, err := r.client.UpdateItem(ctx, updateInput)
	if err != nil {
		errWrap := fmt.Errorf("failed to update member %s: %w", username, err)
		log.Println(errWrap)
		return false, response, errWrap
	}
	return true, response, nil
}

func (r *MemberRepo) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	var rented int
	var messages []string
	user, err := r.GetMemberByUsername(ctx, username, constants.CART)
	if err != nil {
		errWrap := fmt.Errorf("failed to retrieve user %s from cloud. err=%w", username, err)
		log.Println(errWrap)
		return nil, rented, errWrap
	}
	if data.MemberTypes[user.Type] < len(movieIDs)+len(user.Checkedout) {
		return nil, rented, nil
	}

	movies, _, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		errWrap := fmt.Errorf("failed to retrieve movies from cloud: %w", err)
		log.Println(errWrap)
		return nil, rented, errWrap
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

		err := r.checkoutMovie(ctx, user, movie)
		if err != nil {
			messages = append(messages, fmt.Errorf("failed to rent %s\n:Err: %w", movie.ID, err).Error())
			continue
		}
		rented++
	}
	return messages, rented, nil
}

func (r *MemberRepo) checkoutMovie(ctx context.Context, user data.Member, movie data.Movie) error {
	success, err := r.MovieRepo.Rent(ctx, movie)
	if err != nil {
		errWrap := fmt.Errorf("failed to checkout movie %s: %w", movie.Title, err)
		log.Println(errWrap)
		return errWrap
	}
	if !success {
		errWrap := fmt.Errorf("failed to checkout %s", movie.Title)
		log.Println(errWrap)
		return errWrap
	}

	modified, _, err := r.ModifyCart(ctx, user.Username, movie.ID, constants.DELETE, constants.CHECKOUT)
	if !modified || err != nil {
		r.MovieRepo.Return(ctx, movie)
		errWrap := fmt.Errorf("failed to remove %s from %s cart: %w", movie.Title, user.Username, err)
		log.Println(errWrap)
		return errWrap
	}
	return nil
}

func (r *MemberRepo) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	movies, _, err := r.MovieRepo.GetMoviesByID(ctx, movieIDs, constants.NOT_CART)
	if err != nil {
		errWrap := fmt.Errorf("err returning movies. Failed to fetch movies from cloud: %w", err)
		log.Println(errWrap)
		return nil, 0, errWrap
	}

	var returned int
	var messages []string
	name, err := attributevalue.Marshal(username)
	if err != nil {
		errWrap := fmt.Errorf("failed to marshal username %s: %w", username, err)
		log.Println(errWrap)
		return nil, 0, errWrap
	}
	for _, movie := range movies {
		updateInput, err := r.getReturnInput(movie, name)
		if err != nil {
			errWrap := fmt.Errorf("failed to get return input for movie %s\n:Err: %w", movie.ID, err)
			log.Println(errWrap)
			messages = append(messages, errWrap.Error())
			continue
		}
		ok, response, err := r.updateMember(ctx, username, updateInput)
		if err != nil || !ok {
			errWrap := fmt.Errorf("failed to return movie %s\nResponse: %v\nErr: %w", movie.ID, response, err)
			log.Println(errWrap)
			messages = append(messages, errWrap.Error())
			continue
		}

		ok, err = r.MovieRepo.Return(ctx, movie)
		if err != nil || !ok {
			errWrap := fmt.Errorf("failed to update movie table %s\nResponse: %v\nErr: %w", movie.ID, response, err)
			log.Println(errWrap)
			messages = append(messages, errWrap.Error())
			continue
		}
		returned++
	}
	return messages, returned, nil
}

func (r *MemberRepo) getReturnInput(movie data.Movie, name types.AttributeValue) (*dynamodb.UpdateItemInput, error) {

	updateExpr := "DELETE checked_out :checked_out ADD rented :rented"
	expressionAttrs := map[string]types.AttributeValue{
		":checked_out": &types.AttributeValueMemberSS{
			Value: []string{movie.ID},
		},
		":rented": &types.AttributeValueMemberSS{
			Value: []string{movie.ID},
		},
	}
	return &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       map[string]types.AttributeValue{constants.USERNAME: name},
		ExpressionAttributeValues: expressionAttrs,
		ReturnValues:              types.ReturnValueUpdatedNew,
		UpdateExpression:          &updateExpr,
	}, nil
}
