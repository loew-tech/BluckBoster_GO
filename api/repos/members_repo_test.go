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

// --- Mocks for centroid caches ---

type MockCentroidCache struct {
	KNearest    []int
	KNearestErr error
}

// GetMetricsByCentroid implements api_cache.CentroidCacheInterface.
func (m *MockCentroidCache) GetMetricsByCentroid(centroidID int) (data.MovieMetrics, error) {
	panic("unimplemented")
}

// Size implements api_cache.CentroidCacheInterface.
func (m *MockCentroidCache) Size() int {
	return len(m.KNearest)
}

func (m *MockCentroidCache) GetKNearestCentroidsFromMood(mood data.MovieMetrics, k int) ([]int, error) {
	if m.KNearestErr != nil {
		return nil, m.KNearestErr
	}
	return m.KNearest, nil
}

// --- Mocks for centroid to movies caches ---

type MockCentroidsToMoviesCache struct {
	MoviesByCentroid map[int][]string
	Err              error
}

func (m *MockCentroidsToMoviesCache) GetMovieIDsByCentroid(centroidID int) ([]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.MoviesByCentroid[centroidID], nil
}

func (m *MockCentroidsToMoviesCache) GetRandomMovieFromCentroid(centroidID int) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	movies := m.MoviesByCentroid[centroidID]
	if len(movies) == 0 {
		return "", errors.New("no movies")
	}
	return movies[0], nil
}

// --- Mocks for MovieReadRepo ---

type MockReadWriteMovieRepo struct {
	mock.Mock
}

// --- MovieReadRepo methods ---

func (m *MockReadWriteMovieRepo) GetMoviesByPage(ctx context.Context, page string, purpose string) ([]data.Movie, error) {
	args := m.Called(ctx, page, purpose)
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

func setupMemberRepo() (repos.MemberRepoInterface, *MockDynamoClient, *MockReadWriteMovieRepo, *MockCentroidCache, *MockCentroidsToMoviesCache) {
	dynamo := new(MockDynamoClient)
	movieRepo := new(MockReadWriteMovieRepo)
	centroidCache := MockCentroidCache{}
	centroidsToMoviesCache := MockCentroidsToMoviesCache{}
	repo := repos.NewMembersRepo(dynamo, movieRepo, &centroidCache, &centroidsToMoviesCache)
	return repo, dynamo, movieRepo, &centroidCache, &centroidsToMoviesCache
}

const membersTableName = "BluckBoster_members"

func TestGetMemberByUsername_NotFound(t *testing.T) {
	repo, dynamo, _, _, _ := setupMemberRepo()
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
	repo, _, _, _, _ := setupMemberRepo()
	_, err := repo.GetCheckedOutMovies(context.Background(), "")
	assert.EqualError(t, err, "username is required to get checkout moves")
}

func TestSetMemberAPIChoice_Invalid(t *testing.T) {
	repo, _, _, _, _ := setupMemberRepo()
	err := repo.SetMemberAPIChoice(context.Background(), "john", "invalid")
	assert.ErrorContains(t, err, "is not valid api selection")
}

func TestSetMemberAPIChoice_Valid(t *testing.T) {
	repo, dynamo, _, _, _ := setupMemberRepo()
	dynamo.On("UpdateItem", mock.Anything, mock.Anything).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.SetMemberAPIChoice(context.Background(), "john", constants.REST_API)
	assert.NoError(t, err)
}

func TestReturn_MovieError(t *testing.T) {
	repo, _, movieRepo, _, _ := setupMemberRepo()
	movieRepo.On("GetMoviesByID", mock.Anything, []string{"m1"}, constants.NOT_CART).
		Return([]data.Movie{}, errors.New("fetch error"))

	msgs, count, err := repo.Return(context.Background(), "john", []string{"m1"})
	assert.Error(t, err)
	assert.Nil(t, msgs)
	assert.Equal(t, 0, count)
}

func TestGetInitialVotingSlate_Success(t *testing.T) {
	repoIface, _, _, centroidCache, centroidsToMoviesCache := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	// Setup: centroid size > 0, and movies available
	centroidCache.KNearest = []int{1, 2}
	centroidsToMoviesCache.MoviesByCentroid = map[int][]string{
		0: {"a1", "a2"},
		1: {"b1", "b2"},
		2: {"c1", "c2"},
	}

	ctx := context.Background()
	results, err := repo.GetIniitialVotingSlate(ctx)

	assert.NoError(t, err)
	assert.Len(t, results, constants.MAX_MOVIE_SUGGESTIONS)
	// All entries must be non-empty (since we provided movies everywhere)
	for _, mid := range results {
		assert.NotEmpty(t, mid)
	}
}

func TestGetInitialVotingSlate_ErrorOnFetch(t *testing.T) {
	repoIface, _, _, centroidCache, centroidsToMoviesCache := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	// Setup: one centroid exists, but always returns error
	centroidCache.KNearest = []int{1}
	centroidsToMoviesCache.Err = errors.New("fetch failed")

	ctx := context.Background()
	results, err := repo.GetIniitialVotingSlate(ctx)

	// Expect slice length = MAX_MOVIE_SUGGESTIONS
	assert.Len(t, results, constants.MAX_MOVIE_SUGGESTIONS)

	// All slots should be empty since errors prevented movie assignment
	for _, mid := range results {
		assert.Empty(t, mid)
	}
	assert.Error(t, err)
}

func TestGetInitialVotingSlate_NoCentroids(t *testing.T) {
	repoIface, _, _, centroidCache, _ := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	// Zero centroids
	centroidCache.KNearest = []int{}

	ctx := context.Background()
	results, err := repo.GetIniitialVotingSlate(ctx)

	// Expect nil slice and error
	assert.Nil(t, results)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "centroid cache failed to initialize")
}

func TestUpdateMood_AllSuccess(t *testing.T) {
	repo, _, mockMovieRepo, _, _ := setupMemberRepo()

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
	repoIface, _, mockMovieRepo, _, _ := setupMemberRepo()
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
	repoIface, _, mockMovieRepo, _, _ := setupMemberRepo()
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
	repoIface, _, mockMovieRepo, _, _ := setupMemberRepo()
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

func TestGetVotingFinalPicks_Success(t *testing.T) {
	repoIface, _, mockMovieRepo, centroidCache, centroidsToMoviesCache := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	// Mock centroid cache returning centroid IDs
	centroidCache.KNearest = []int{1, 2}

	// Mock centroid-to-movies mapping
	centroidsToMoviesCache.MoviesByCentroid = map[int][]string{
		1: {"m1", "n1"},
		2: {"m2", "n2"},
	}

	// Mock metrics for the movies
	m1 := data.MovieMetrics{Acting: 1, Action: 1}
	m2 := data.MovieMetrics{Acting: 5, Action: 5}

	n1 := data.MovieMetrics{Acting: 100, Action: 100}
	n2 := data.MovieMetrics{Acting: 100, Action: 100}

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m1").Return(m1, nil)
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "m2").Return(m2, nil)

	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "n1").Return(n1, nil)
	mockMovieRepo.On("GetMovieMetrics", mock.Anything, "n2").Return(n2, nil)

	ctx := context.Background()
	mood := data.MovieMetrics{Acting: 2, Action: 2}

	results, err := repo.GetVotingFinalPicks(ctx, mood)
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.Contains(t, []string{"m1", "m2"}, results[0]) // should pick a valid nearest neighbor
	mockMovieRepo.AssertExpectations(t)
}

func TestGetVotingFinalPicks_CentroidError(t *testing.T) {
	repoIface, _, _, centroidCache, _ := setupMemberRepo()
	repo := repoIface.(*repos.MemberRepo)

	centroidCache.KNearestErr = errors.New("centroid fail")

	ctx := context.Background()
	mood := data.MovieMetrics{Acting: 1, Action: 2}

	results, err := repo.GetVotingFinalPicks(ctx, mood)
	assert.Error(t, err)
	assert.Nil(t, results)
}
