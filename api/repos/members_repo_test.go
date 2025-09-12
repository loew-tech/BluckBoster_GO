package repos_test

import (
	"context"
	"errors"
	"testing"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReadWriteMovieRepo struct {
	mock.Mock
}

// --- MovieReadRepo methods ---

func (m *MockReadWriteMovieRepo) GetMoviesByPage(ctx context.Context, page string, forGraph bool) ([]data.Movie, error) {
	args := m.Called(ctx, page, forGraph)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockReadWriteMovieRepo) GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error) {
	args := m.Called(ctx, movieID, forCart)
	return args.Get(0).(data.Movie), args.Error(1)
}

func (m *MockReadWriteMovieRepo) GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error) {
	args := m.Called(ctx, movieIDs, forCart)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockReadWriteMovieRepo) GetMovieMetrics(ctx context.Context, movieID string) (data.MovieMetrics, error) {
	args := m.Called(ctx, movieID)
	return args.Get(0).(data.MovieMetrics), args.Error(1)
}

func (m *MockReadWriteMovieRepo) GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error) {
	args := m.Called(ctx, movieID)
	return args.Get(0).(data.MovieTrivia), args.Error(1)
}

// --- MovieInventoryRepo methods ---

func (m *MockReadWriteMovieRepo) Rent(ctx context.Context, movie data.Movie) (bool, error) {
	args := m.Called(ctx, movie)
	return args.Bool(0), args.Error(1)
}

func (m *MockReadWriteMovieRepo) Return(ctx context.Context, movie data.Movie) (bool, error) {
	args := m.Called(ctx, movie)
	return args.Bool(0), args.Error(1)
}

func setupMemberRepo() (repos.MemberRepoInterface, *MockDynamoClient, *MockReadWriteMovieRepo) {
	dynamo := new(MockDynamoClient)
	movieRepo := new(MockReadWriteMovieRepo)
	repo := repos.NewMembersRepo(dynamo, movieRepo)
	return repo, dynamo, movieRepo
}

const membersTableName = "BluckBoster_members"

func TestGetMemberByUsername_NotFound(t *testing.T) {
	repo, dynamo, _ := setupMemberRepo()
	input := &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"username": &types.AttributeValueMemberS{Value: "notfound"}},
		TableName: aws.String(membersTableName),
	}
	dynamo.On("GetItem", mock.Anything, input).
		Return(&dynamodb.GetItemOutput{Item: nil}, nil)

	_, err := repo.GetMemberByUsername(context.Background(), "notfound", false)
	assert.Error(t, err)
}

func TestGetCheckedOutMovies_EmptyUsername(t *testing.T) {
	repo, _, _ := setupMemberRepo()
	_, err := repo.GetCheckedOutMovies(context.Background(), "")
	assert.EqualError(t, err, "username is required to get checkout moves")
}

func TestSetMemberAPIChoice_Invalid(t *testing.T) {
	repo, _, _ := setupMemberRepo()
	err := repo.SetMemberAPIChoice(context.Background(), "john", "invalid")
	assert.ErrorContains(t, err, "is not valid api selection")
}

func TestSetMemberAPIChoice_Valid(t *testing.T) {
	repo, dynamo, _ := setupMemberRepo()
	dynamo.On("UpdateItem", mock.Anything, mock.Anything).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.SetMemberAPIChoice(context.Background(), "john", constants.REST_API)
	assert.NoError(t, err)
}

func TestReturn_MovieError(t *testing.T) {
	repo, _, movieRepo := setupMemberRepo()
	movieRepo.On("GetMoviesByID", mock.Anything, []string{"m1"}, constants.NOT_CART).
		Return([]data.Movie{}, errors.New("fetch error"))

	msgs, count, err := repo.Return(context.Background(), "john", []string{"m1"})
	assert.Error(t, err)
	assert.Nil(t, msgs)
	assert.Equal(t, 0, count)
}

func TestUpdateMood_Success(t *testing.T) {
	repo, _, mockMovieRepo := setupMemberRepo()

	// Fake movie metrics for 2 movies
	m1 := data.MovieMetrics{Acting: 80, Action: 60, Cinematography: 70}
	m2 := data.MovieMetrics{Acting: 90, Action: 40, Cinematography: 75}

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m1").Return(m1, nil)
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m2").Return(m2, nil)

	cur := data.MovieMetrics{Acting: 70, Action: 50}
	result, err := repo.UpdateMood(context.Background(), cur, 1, []string{"m1", "m2"})

	assert.NoError(t, err)
	assert.True(t, result.Acting > 0)
	assert.True(t, result.Cinematography > 0)
	// sanity check average effect
	assert.GreaterOrEqual(t, result.Acting, cur.Acting)
	assert.Equal(t, (cur.Action+m1.Action+m2.Action)/3, result.Action)
	assert.Equal(t, (m1.Cinematography+m2.Cinematography)/3, result.Cinematography)
	mockMovieRepo.AssertExpectations(t)
}

func TestUpdateMood_ErrorFromMovieRepo(t *testing.T) {
	repo, _, mockMovieRepo := setupMemberRepo()

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "badID").Return(data.MovieMetrics{}, errors.New("metrics not found"))

	cur := data.MovieMetrics{Acting: 50, Action: 50, Cinematography: 50}
	result, err := repo.UpdateMood(context.Background(), cur, 1, []string{"badID"})

	// Even with error, function returns averaged current mood
	assert.NoError(t, err) // errors.Join([]) → nil
	assert.NotZero(t, result.Acting)
	mockMovieRepo.AssertExpectations(t)
}
