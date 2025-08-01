package repos_test

import (
	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Implement mocked methods as needed

func TestGetMoviesByPage_Success(t *testing.T) {
	mockClient := new(MockDynamoClient)
	repo := reposTestWrapper(mockClient)

	// setup fake response
	fakeOutput := &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				constants.ID:   &types.AttributeValueMemberS{Value: "1"},
				"title":        &types.AttributeValueMemberS{Value: "Movie 1"},
				"inventory":    &types.AttributeValueMemberN{Value: "5"},
				"rented":       &types.AttributeValueMemberN{Value: "2"},
				"rating":       &types.AttributeValueMemberS{Value: "PG-13"},
				"review":       &types.AttributeValueMemberS{Value: "foo-bar"},
				constants.YEAR: &types.AttributeValueMemberS{Value: "2020"},
				constants.CAST: &types.AttributeValueMemberSS{Value: []string{"Actor 1", "Actor 2"}},
				"director":     &types.AttributeValueMemberS{Value: "Director A"},
			},
		},
	}

	mockClient.On("Query", mock.Anything, mock.Anything).
		Return(fakeOutput, nil)

	ctx := context.Background()
	movies, err := repo.GetMoviesByPage(ctx, "A", false)
	assert.NoError(t, err)
	assert.Len(t, movies, 1)
	assert.Equal(t, "Movie 1", movies[0].Title)
}

func TestGetMovieByID_EmptyID(t *testing.T) {
	repo := reposTestWrapper(new(MockDynamoClient))
	ctx := context.Background()
	_, err := repo.GetMovieByID(ctx, "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "movieID cannot be empty")
}

func TestGetMoviesByID_BatchLimitExceeded(t *testing.T) {
	repo := reposTestWrapper(new(MockDynamoClient))
	ids := make([]string, 11)
	ctx := context.Background()
	_, err := repo.GetMoviesByID(ctx, ids, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "batch size")
}

func TestRentMovie_Success(t *testing.T) {
	mockClient := new(MockDynamoClient)
	repo := reposTestWrapper(mockClient)

	mockClient.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	ctx := context.Background()
	movie := data.Movie{ID: "123", Inventory: 5, Rented: 3}
	ok, err := repo.Rent(ctx, movie)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestReturnMovie_Error(t *testing.T) {
	mockClient := new(MockDynamoClient)
	repo := reposTestWrapper(mockClient)

	someErr := errors.New("some dynamodb error")
	mockClient.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, someErr)

	ctx := context.Background()
	movie := data.TestMovies[0]
	ok, err := repo.Return(ctx, movie)
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "updating inventory")
}

func reposTestWrapper(client repos.DynamoClientInterface) *repos.DynamoMovieRepo {
	return repos.NewDynamoMovieRepo(client)
}
