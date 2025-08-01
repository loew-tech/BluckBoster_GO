package services_test

import (
	"context"
	"errors"
	"testing"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMemberRepo struct {
	mock.Mock
}

func (m *MockMemberRepo) GetMemberByUsername(ctx context.Context, username string, forCart bool) (data.Member, error) {
	args := m.Called(ctx, username, forCart)
	return args.Get(0).(data.Member), args.Error(1)
}

func (m *MockMemberRepo) GetCartMovies(ctx context.Context, username string) ([]data.Movie, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMemberRepo) ModifyCart(ctx context.Context, username, movieID, action string, checkingOut bool) (bool, error) {
	args := m.Called(ctx, username, movieID, action, checkingOut)
	return args.Bool(0), args.Error(1)
}

func (m *MockMemberRepo) Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	args := m.Called(ctx, username, movieIDs)
	return args.Get(0).([]string), args.Int(1), args.Error(2)
}

func (m *MockMemberRepo) Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error) {
	args := m.Called(ctx, username, movieIDs)
	return args.Get(0).([]string), args.Int(1), args.Error(2)
}

func (m *MockMemberRepo) GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]data.Movie), args.Error(1)
}

func (m *MockMemberRepo) SetMemberAPIChoice(ctx context.Context, username, apiChoice string) error {
	args := m.Called(ctx, username, apiChoice)
	return args.Error(0)
}

func setupMockService() (*services.MembersService, *MockMemberRepo) {
	repo := new(MockMemberRepo)
	service := services.NewMemberServiceWithRepo(repo)
	return service, repo
}

func TestLogin_Success(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetMemberByUsername", mock.Anything, "john", constants.NOT_CART).
		Return(data.Member{Username: "john"}, nil)

	member, err := service.Login(context.Background(), "john")
	assert.NoError(t, err)
	assert.Equal(t, "john", member.Username)
}

func TestLogin_Failure(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetMemberByUsername", mock.Anything, "unknown", constants.NOT_CART).
		Return(data.Member{}, errors.New("not found"))

	_, err := service.Login(context.Background(), "unknown")
	assert.Error(t, err)
}

func TestGetMember_Success(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetMemberByUsername", mock.Anything, "alice", true).
		Return(data.Member{Username: "alice"}, nil)

	member, err := service.GetMember(context.Background(), "alice", true)
	assert.NoError(t, err)
	assert.Equal(t, "alice", member.Username)
}

func TestGetMember_Failure(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetMemberByUsername", mock.Anything, "fail", false).
		Return(data.Member{}, errors.New("db error"))

	_, err := service.GetMember(context.Background(), "fail", false)
	assert.Error(t, err)
}

func TestAddToCart_Success(t *testing.T) {
	service, repo := setupMockService()
	repo.On("ModifyCart", mock.Anything, "bob", "m1", constants.ADD, constants.NOT_CHECKOUT).
		Return(true, nil)

	ok, err := service.AddToCart(context.Background(), "bob", "m1")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestRemoveFromCart_Error(t *testing.T) {
	service, repo := setupMockService()
	repo.On("ModifyCart", mock.Anything, "bob", "m1", constants.DELETE, constants.NOT_CHECKOUT).
		Return(false, errors.New("failed"))

	ok, err := service.RemoveFromCart(context.Background(), "bob", "m1")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestCheckout_Success(t *testing.T) {
	service, repo := setupMockService()
	repo.On("Checkout", mock.Anything, "john", []string{"m1", "m2"}).
		Return([]string{"ok"}, 2, nil)

	msgs, count, err := service.Checkout(context.Background(), "john", []string{"m1", "m2"})
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, []string{"ok"}, msgs)
}

func TestReturn_Failure(t *testing.T) {
	service, repo := setupMockService()
	repo.On("Return", mock.Anything, "jane", []string{"m1"}).
		Return([]string{}, 0, errors.New("db error"))

	msgs, count, err := service.Return(context.Background(), "jane", []string{"m1"})
	assert.Error(t, err)
	assert.Empty(t, msgs)
	assert.Equal(t, 0, count)
}

func TestGetCheckedOutMovies(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetCheckedOutMovies", mock.Anything, "john").
		Return([]data.Movie{{ID: "m1"}}, nil)

	movies, err := service.GetCheckedOutMovies(context.Background(), "john")
	assert.NoError(t, err)
	assert.Len(t, movies, 1)
}

func TestGetCartIDs(t *testing.T) {
	service, repo := setupMockService()
	repo.On("GetMemberByUsername", mock.Anything, "john", constants.CART).
		Return(data.Member{Cart: []string{"m1", "m2"}}, nil)

	cart, err := service.GetCartIDs(context.Background(), "john")
	assert.NoError(t, err)
	assert.Equal(t, []string{"m1", "m2"}, cart)
}

func TestSetAPIChoice(t *testing.T) {
	service, repo := setupMockService()
	repo.On("SetMemberAPIChoice", mock.Anything, "john", "omdb").Return(nil)

	err := service.SetAPIChoice(context.Background(), "john", "omdb")
	assert.NoError(t, err)
}

func TestSetAPIChoice_Error(t *testing.T) {
	service, repo := setupMockService()
	repo.On("SetMemberAPIChoice", mock.Anything, "john", "bad").Return(errors.New("fail"))

	err := service.SetAPIChoice(context.Background(), "john", "bad")
	assert.Error(t, err)
}
