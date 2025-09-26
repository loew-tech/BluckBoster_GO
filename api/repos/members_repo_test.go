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
	movieRepo.On("GetMoviesByPage", mock.Anything, mock.Anything, mock.Anything).Return([]data.Movie{}, nil)
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

func TestUpdateMood_AllSuccess(t *testing.T) {
	repo, _, mockMovieRepo := setupMemberRepo()

	m1 := data.MovieMetrics{Acting: 2, Action: 3}
	m2 := data.MovieMetrics{Acting: 4, Action: 5}

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m1").Return(m1, nil)
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m2").Return(m2, nil)

	ctx := context.Background()
	currentMood := data.MovieMetrics{Acting: 1, Action: 1}
	iteration := 2
	movieIDs := []string{"m1", "m2"}

	result, err := repo.UpdateMood(ctx, currentMood, iteration, movieIDs)
	assert.NoError(t, err)
	assert.InDelta(t, 2.0, result.Acting, 0.01)
	assert.InDelta(t, 2.5, result.Action, 0.01)
	mockMovieRepo.AssertExpectations(t)
}

func TestUpdateMood_SomeErrors(t *testing.T) {
	repoIface, _, mockMovieRepo := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	m1 := data.MovieMetrics{Acting: 2, Action: 3}

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m1").Return(m1, nil)
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m2").Return(data.MovieMetrics{}, errors.New("not found"))

	ctx := context.Background()
	currentMood := data.MovieMetrics{Acting: 1, Action: 1}
	iteration := 1
	movieIDs := []string{"m1", "m2"}

	result, err := repo.UpdateMood(ctx, currentMood, iteration, movieIDs)
	assert.Error(t, err)
	assert.InDelta(t, 1.5, result.Acting, 0.01)
	assert.InDelta(t, 2.0, result.Action, 0.01)
	mockMovieRepo.AssertExpectations(t)
}

func TestUpdateMood_AllErrors(t *testing.T) {
	repoIface, _, mockMovieRepo := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m1").Return(data.MovieMetrics{}, errors.New("not found"))
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m2").Return(data.MovieMetrics{}, errors.New("not found"))

	ctx := context.Background()
	currentMood := data.MovieMetrics{Acting: 5, Action: 5}
	iteration := 2
	movieIDs := []string{"m1", "m2"}

	result, err := repo.UpdateMood(ctx, currentMood, iteration, movieIDs)
	assert.Error(t, err)
	assert.InDelta(t, 5.0, result.Acting, 0.01)
	assert.InDelta(t, 5.0, result.Action, 0.01)
	mockMovieRepo.AssertExpectations(t)
}

func TestUpdateMood_NoMovies(t *testing.T) {
	repoIface, _, mockMovieRepo := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	ctx := context.Background()
	currentMood := data.MovieMetrics{Acting: 7, Action: 8}
	iteration := 3
	movieIDs := []string{}

	result, err := repo.UpdateMood(ctx, currentMood, iteration, movieIDs)
	assert.NoError(t, err)
	assert.InDelta(t, 7.0, result.Acting, 0.01)
	assert.InDelta(t, 8.0, result.Action, 0.01)
	mockMovieRepo.AssertExpectations(t)
}
