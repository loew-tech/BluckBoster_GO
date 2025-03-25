package data

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const movieTableName = "BluckBoster_movies"

type MovieRepo struct {
	client    dynamodb.Client
	tableName string
}

func NewMovieRepo(client *dynamodb.Client) MovieRepo {
	return MovieRepo{
		client:    *client,
		tableName: movieTableName,
	}
}

func (r MovieRepo) GetAllMovies() ([]Movie, error) {
	params := &dynamodb.ScanInput{
		TableName: &r.tableName,
	}
	result, err := r.client.Scan(context.TODO(), params)
	if err != nil {
		log.Printf("Query API call failed: %s\n", err)
		return nil, err
	}

	movies := make([]Movie, 0)
	for _, i := range result.Items {
		movie := Movie{}
		err = attributevalue.UnmarshalMap(i, &movie)

		if err != nil {
			log.Printf("Got error unmarshalling: %s\n", err)
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (r MovieRepo) GetMovieByID(movieID string, forCart bool) (Movie, CartMovie, error) {
	expr := "#i, title, inventory"
	exprAttrNames := map[string]string{"#i": "id"}
	if !forCart {
		expr = fmt.Sprintf("%s, #c, director, rented, rating, review, synopsis, #y", expr)
		exprAttrNames["#c"], exprAttrNames["#y"] = CAST, YEAR
	}
	input := &dynamodb.GetItemInput{
		Key:                      map[string]types.AttributeValue{ID: &types.AttributeValueMemberS{Value: movieID}},
		TableName:                aws.String(r.tableName),
		ProjectionExpression:     &expr,
		ExpressionAttributeNames: exprAttrNames,
	}

	result, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return Movie{}, CartMovie{}, err
	}

	movie, cartMovie := Movie{}, CartMovie{}
	if !forCart {
		err = attributevalue.UnmarshalMap(result.Item, &movie)
	} else {
		err = attributevalue.UnmarshalMap(result.Item, &cartMovie)
	}

	if err != nil {
		log.Printf("Failed to unmarhal movie")
	}
	return movie, cartMovie, err
}

func (r MovieRepo) GetMoviesByID(movieIDs []string, forCart bool) ([]Movie, []CartMovie, error) {
	if len(movieIDs) == 0 {
		return make([]Movie, 0), make([]CartMovie, 0), nil
	}
	keys := make([]map[string]types.AttributeValue, 0)
	for _, mid := range movieIDs {
		keys = append(keys, map[string]types.AttributeValue{ID: &types.AttributeValueMemberS{Value: mid}})
	}
	expr := "#i, title, inventory"
	exprAttrNames := map[string]string{"#i": "id"}
	if !forCart {
		expr = fmt.Sprintf("%s, #c, director, rented, rating, review, synopsis, #y", expr)
		exprAttrNames["#c"], exprAttrNames["#y"] = CAST, YEAR
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.tableName: {
				Keys:                     keys,
				ProjectionExpression:     aws.String(expr),
				ExpressionAttributeNames: exprAttrNames,
			},
		},
	}

	result, err := r.client.BatchGetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Err fetching movies from cloud: %s\n", err)
		return nil, nil, err
	}

	movies, cartMovies := make([]Movie, 0), make([]CartMovie, 0)
	for _, v := range result.Responses {
		for _, m := range v {
			if forCart {
				// @TODO: refactor to remove deep leveling
				cartMovie := CartMovie{}
				if err = attributevalue.UnmarshalMap(m, &cartMovie); err != nil {
					log.Printf("Got error unmarshalling: %s", err)
					continue
				}
				cartMovies = append(cartMovies, cartMovie)
			} else {
				movie := Movie{}
				if err = attributevalue.UnmarshalMap(m, &movie); err != nil {
					log.Printf("Got error unmarshalling: %s", err)
					continue
				}
				movies = append(movies, movie)
			}
		}
	}
	return movies, cartMovies, nil
}

func (r MovieRepo) Rent(movie Movie) (bool, error) {
	input, err := getUpdateInventoryInput(movie, -1)
	if err != nil {
		return false, fmt.Errorf("failed to generate input for update call %s", err)
	}
	return r.updateInventory(movie, input)
}

func (r MovieRepo) Return(movie Movie) (bool, error) {
	input, err := getUpdateInventoryInput(movie, 1)
	if err != nil {
		return false, fmt.Errorf("failed to generate input for update call %s", err)
	}
	return r.updateInventory(movie, input)
}

func getUpdateInventoryInput(movie Movie, inventoryInc int) (*dynamodb.UpdateItemInput, error) {
	mid, err := attributevalue.Marshal(movie.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data %s", err)
	}

	return &dynamodb.UpdateItemInput{
		TableName: aws.String(movieTableName),
		Key:       map[string]types.AttributeValue{ID: mid},
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

func (r MovieRepo) updateInventory(movie Movie, input *dynamodb.UpdateItemInput) (bool, error) {
	response, err := r.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to update movie item %s\nResp %v\nErr: %s\n", movie.ID, response, err)
		return false, err
	}
	return true, nil
}
