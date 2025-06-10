package repos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"blockbuster/api/constants"
	"blockbuster/api/data"
)

const movieTableName = "BluckBoster_movies"

type MovieRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewMovieRepo(client *dynamodb.Client) *MovieRepo {
	return &MovieRepo{
		client:    client,
		tableName: movieTableName,
	}
}

func (r *MovieRepo) GetMoviesByPage(ctx context.Context, forGraph bool, page string) ([]data.Movie, error) {
	var projectionExpr *string
	expr := "#i, title, #c, director"

	exprAttrNames := map[string]string{
		"#i": constants.ID,
		"#c": constants.CAST,
	}
	if !forGraph {
		expr = fmt.Sprintf("%s, inventory, rented, #y", expr)
		exprAttrNames["#y"] = constants.YEAR
	}
	projectionExpr = &expr

	input := &dynamodb.QueryInput{
		TableName: &r.tableName,
		IndexName: aws.String(constants.PAGINATE_KEY_INDEX),
		KeyConditions: map[string]types.Condition{
			constants.PAGINATE_KEY: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: page},
				},
			},
		},
		ProjectionExpression:     projectionExpr,
		ExpressionAttributeNames: exprAttrNames,
	}

	response, err := r.client.Query(ctx, input)
	if err != nil {
		errWrap := fmt.Errorf("failed querying movies by page %s: %w", page, err)
		log.Println(errWrap)
		return nil, errWrap
	}

	var movies []data.Movie
	err = attributevalue.UnmarshalListOfMaps(response.Items, &movies)
	if err != nil {
		errWrap := fmt.Errorf("err unmarshalling movies from query response: %w", err)
		log.Print(errWrap)
		return nil, errWrap
	}

	return movies, nil
}

func (r *MovieRepo) GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error) {
	input, err := r.getMovieByIDInput(movieID, forCart)
	if input == nil || err != nil {
		errWrap := fmt.Errorf("failed to create GetItemInput for movieID %s: %w", movieID, err)
		log.Print(errWrap)
		return data.Movie{}, errWrap
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return data.Movie{}, err
	}
	return r.getMovieByIDResult(result, forCart)
}

func (r *MovieRepo) getMovieByIDInput(movieID string, forCart bool) (*dynamodb.GetItemInput, error) {
	if movieID == "" {
		return nil, errors.New("movieID cannot be empty")
	}
	expr := "#i, title, inventory"
	exprAttrNames := map[string]string{"#i": "id"}
	if !forCart {
		expr = fmt.Sprintf("%s, #c, director, rented, rating, review, synopsis, trivia, #y", expr)
		exprAttrNames["#c"], exprAttrNames["#y"] = constants.CAST, constants.YEAR
	}
	return &dynamodb.GetItemInput{
		Key:                      map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: movieID}},
		TableName:                aws.String(r.tableName),
		ProjectionExpression:     &expr,
		ExpressionAttributeNames: exprAttrNames,
	}, nil
}

func (r *MovieRepo) getMovieByIDResult(result *dynamodb.GetItemOutput, forCart bool) (data.Movie, error) {
	if result.Item == nil {
		return data.Movie{}, errors.New("movie not found")
	}

	var movie data.Movie
	err := attributevalue.UnmarshalMap(result.Item, &movie)
	if err != nil {
		log.Printf("Failed to unmarshal movie: %s", err)
		return data.Movie{}, err
	}
	return movie, nil
}

func (r *MovieRepo) GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error) {
	if len(movieIDs) == 0 {
		return make([]data.Movie, 0), nil
	}
	keys := make([]map[string]types.AttributeValue, 0, len(movieIDs))
	for _, mid := range movieIDs {
		keys = append(keys, map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: mid}})
	}
	var projectionExpr *string
	var exprAttrNames map[string]string
	if forCart {
		expr := "#i, title, inventory"
		projectionExpr = &expr
		exprAttrNames = make(map[string]string)
		exprAttrNames["#i"] = constants.ID
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.tableName: {
				Keys:                     keys,
				ProjectionExpression:     projectionExpr,
				ExpressionAttributeNames: exprAttrNames,
			},
		},
	}

	result, err := r.client.BatchGetItem(ctx, input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, err
	}

	var movies []data.Movie
	for _, v := range result.Responses {
		var batch []data.Movie
		if err = attributevalue.UnmarshalListOfMaps(v, &batch); err != nil {
			log.Printf("Got error unmarshalling movies: %s", err)
			return nil, err
		}
		movies = append(movies, batch...)
	}
	return movies, nil
}

func (r *MovieRepo) Rent(ctx context.Context, movie data.Movie) (bool, error) {
	input, err := getUpdateInventoryInput(movie, -1)
	if err != nil {
		return false, fmt.Errorf("failed to generate input for update call %w", err)
	}
	return r.updateInventory(ctx, movie, input)
}

func (r *MovieRepo) Return(ctx context.Context, movie data.Movie) (bool, error) {
	input, err := getUpdateInventoryInput(movie, 1)
	if err != nil {
		return false, fmt.Errorf("failed to generate input for update call %w", err)
	}
	return r.updateInventory(ctx, movie, input)
}

func getUpdateInventoryInput(movie data.Movie, inventoryInc int) (*dynamodb.UpdateItemInput, error) {
	mid, err := attributevalue.Marshal(movie.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data %w", err)
	}

	return &dynamodb.UpdateItemInput{
		TableName: aws.String(movieTableName),
		Key:       map[string]types.AttributeValue{constants.ID: mid},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inventory": &types.AttributeValueMemberN{
				Value: strconv.Itoa(movie.Inventory + inventoryInc),
			},
			":rented": &types.AttributeValueMemberN{
				Value: strconv.Itoa(movie.Rented - inventoryInc),
			},
		},
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String("set inventory = :inventory, rented = :rented"),
	}, nil
}

func (r *MovieRepo) updateInventory(ctx context.Context, movie data.Movie, input *dynamodb.UpdateItemInput) (bool, error) {
	response, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		log.Printf("Failed to update movie item %s\nResp %v\nErr: %s\n", movie.ID, response, err)
		return false, err
	}
	return true, nil
}

func (r *MovieRepo) GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error) {
	input := &dynamodb.GetItemInput{
		Key:             map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: movieID}},
		TableName:       &r.tableName,
		AttributesToGet: []string{constants.TRIVIA},
	}

	response, err := r.client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Failed to get trivia for %s\nResp %v\nErr: %s\n", movieID, response, err)
		return data.MovieTrivia{}, err
	}

	trivia := data.MovieTrivia{}
	if err = attributevalue.UnmarshalMap(response.Item, &trivia); err != nil {
		log.Printf("Got error unmarshalling: %s", err)
		return data.MovieTrivia{}, err
	}
	return trivia, nil
}
