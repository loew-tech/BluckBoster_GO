package repos

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/utils"
)

const movieTableName = "BluckBoster_movies"

type DynamoMovieRepo struct {
	client    DynamoClientInterface
	tableName string
}

func NewDynamoMovieRepo(client DynamoClientInterface) *DynamoMovieRepo {
	return &DynamoMovieRepo{
		client:    client,
		tableName: movieTableName,
	}
}

func (r *DynamoMovieRepo) GetMoviesByPage(ctx context.Context, page string, forGraph bool) ([]data.Movie, error) {
	expr := "#i, title, #c, director"
	exprAttrNames := map[string]string{
		"#i": constants.ID,
		"#c": constants.CAST,
	}
	if !forGraph {
		expr += ", inventory, rented, rating, #y"
		exprAttrNames["#y"] = constants.YEAR
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		IndexName: aws.String(constants.PAGINATE_KEY_INDEX),
		KeyConditions: map[string]types.Condition{
			constants.PAGINATE_KEY: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: page},
				},
			},
		},
		ProjectionExpression:     &expr,
		ExpressionAttributeNames: exprAttrNames,
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, utils.LogError(fmt.Sprintf("querying movies by page %s", page), err)
	}

	var movies []data.Movie
	err = attributevalue.UnmarshalListOfMaps(result.Items, &movies)
	if err != nil {
		return nil, utils.LogError("unmarshalling movies from query response", err)
	}

	return movies, nil
}

func (r *DynamoMovieRepo) GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error) {
	input, err := r.getMovieByIDInput(movieID, forCart)
	if err != nil {
		return data.Movie{}, utils.LogError("creating GetItemInput", err)
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return data.Movie{}, utils.LogError("fetching movie from DynamoDB", err)
	}

	return r.getMovieByIDResult(result)
}

func (r *DynamoMovieRepo) getMovieByIDInput(movieID string, forCart bool) (*dynamodb.GetItemInput, error) {
	if movieID == "" {
		return nil, errors.New("movieID cannot be empty")
	}
	expr := "#i, title, inventory"
	exprAttrNames := map[string]string{"#i": constants.ID}
	if !forCart {
		expr = fmt.Sprintf("%s, #c, director, rented, rating, review, synopsis, trivia, #y", expr)
		exprAttrNames["#c"] = constants.CAST
		exprAttrNames["#y"] = constants.YEAR
	}

	return &dynamodb.GetItemInput{
		TableName:                aws.String(r.tableName),
		Key:                      map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: movieID}},
		ProjectionExpression:     &expr,
		ExpressionAttributeNames: exprAttrNames,
	}, nil
}

func (r *DynamoMovieRepo) getMovieByIDResult(result *dynamodb.GetItemOutput) (data.Movie, error) {
	if result.Item == nil {
		return data.Movie{}, utils.LogError("movie not found", nil)
	}
	var movie data.Movie
	err := attributevalue.UnmarshalMap(result.Item, &movie)
	if err != nil {
		return data.Movie{}, utils.LogError("unmarshalling movie", err)
	}
	return movie, nil
}

func (r *DynamoMovieRepo) GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error) {
	if len(movieIDs) == 0 {
		return []data.Movie{}, nil
	}
	if len(movieIDs) > 10 {
		return nil, utils.LogError("batch size exceeds DynamoDB 10-item limit", nil)
	}

	keys := make([]map[string]types.AttributeValue, 0, len(movieIDs))
	for _, id := range movieIDs {
		keys = append(keys, map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: id}})
	}

	expr := "#i, title, inventory"
	exprAttrNames := map[string]string{"#i": constants.ID}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.tableName: {
				Keys:                     keys,
				ProjectionExpression:     &expr,
				ExpressionAttributeNames: exprAttrNames,
			},
		},
	}

	result, err := r.client.BatchGetItem(ctx, input)
	if err != nil {
		return nil, utils.LogError("batch fetching movies", err)
	}

	var movies []data.Movie
	for _, v := range result.Responses {
		var batch []data.Movie
		if err := attributevalue.UnmarshalListOfMaps(v, &batch); err != nil {
			return nil, utils.LogError("unmarshalling batch movies", err)
		}
		movies = append(movies, batch...)
	}
	return movies, nil
}

func (r *DynamoMovieRepo) GetMovieMetrics(ctx context.Context, movieID string) (data.MovieMetrics, error) {
	if movieID == "" {
		return data.MovieMetrics{}, utils.LogError("movieID cannot be empty", nil)
	}

	input := &dynamodb.GetItemInput{
		Key:                  map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: movieID}},
		TableName:            aws.String(r.tableName),
		ProjectionExpression: aws.String(constants.METRICS),
	}
	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return data.MovieMetrics{}, utils.LogError("fetching movie metrics from DynamoDB", err)
	}

	// look only at the "mets" map
	attr, ok := result.Item[constants.METRICS]
	if !ok {
		return data.MovieMetrics{}, utils.LogError("mets attribute not found", nil)
	}

	var metrics data.MovieMetrics
	if err := attributevalue.Unmarshal(attr, &metrics); err != nil {
		return data.MovieMetrics{}, utils.LogError("unmarshalling movie metrics", err)
	}
	return metrics, nil
}

func (r *DynamoMovieRepo) Rent(ctx context.Context, movie data.Movie) (bool, error) {
	input := getUpdateInventoryInput(movie, constants.RENT_MOVIE_INC)
	return r.updateInventory(ctx, movie, input)
}

func (r *DynamoMovieRepo) Return(ctx context.Context, movie data.Movie) (bool, error) {
	input := getUpdateInventoryInput(movie, constants.RETURN_MOVIE_INC)
	return r.updateInventory(ctx, movie, input)
}

func getUpdateInventoryInput(movie data.Movie, inventoryDelta int) *dynamodb.UpdateItemInput {
	mid := &types.AttributeValueMemberS{Value: movie.ID}

	return &dynamodb.UpdateItemInput{
		TableName: aws.String(movieTableName),
		Key:       map[string]types.AttributeValue{constants.ID: mid},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inventory": &types.AttributeValueMemberN{Value: strconv.Itoa(movie.Inventory + inventoryDelta)},
			":rented":    &types.AttributeValueMemberN{Value: strconv.Itoa(movie.Rented - inventoryDelta)},
		},
		UpdateExpression: aws.String("SET inventory = :inventory, rented = :rented"),
		ReturnValues:     types.ReturnValueUpdatedNew,
	}
}

func (r *DynamoMovieRepo) updateInventory(ctx context.Context, movie data.Movie, input *dynamodb.UpdateItemInput) (bool, error) {
	_, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return false, utils.LogError(fmt.Sprintf("updating inventory for movie %s", movie.ID), err)
	}
	return true, nil
}

func (r *DynamoMovieRepo) GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error) {
	input := &dynamodb.GetItemInput{
		Key:             map[string]types.AttributeValue{constants.ID: &types.AttributeValueMemberS{Value: movieID}},
		TableName:       &r.tableName,
		AttributesToGet: []string{constants.TRIVIA},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return data.MovieTrivia{}, utils.LogError("fetching movie trivia", err)
	}

	var trivia data.MovieTrivia
	if err := attributevalue.UnmarshalMap(result.Item, &trivia); err != nil {
		return data.MovieTrivia{}, utils.LogError("unmarshalling movie trivia", err)
	}
	return trivia, nil
}
