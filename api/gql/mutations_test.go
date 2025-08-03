package gql_test

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/constants"
	"blockbuster/api/gql"
	"blockbuster/api/services"
)

func setupTestMemberService() *services.MockMembersService {
	mockSvc := new(services.MockMembersService)
	gql.SetMemberService(mockSvc)
	return mockSvc
}

func TestReturnRentalsField(t *testing.T) {
	mockSvc := setupTestMemberService()
	mockSvc.On("Return", mock.Anything, "alice", []string{"1", "2"}).Return([]string{"Returned 1", "Returned 2"}, 2, nil)

	ctx := context.WithValue(context.Background(), gql.GinContextKey, context.Background())
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.USERNAME:  "alice",
			constants.MOVIE_IDS: []interface{}{"1", "2"},
		},
		Context: ctx,
	}

	res, err := gql.ReturnRentalsField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Returned 1", "Returned 2"}, res)
}

func TestUpdateCartField_Add(t *testing.T) {
	mockSvc := setupTestMemberService()
	mockSvc.On("AddToCart", mock.Anything, "bob", "5").Return(true, nil)

	ctx := context.WithValue(context.Background(), gql.GinContextKey, context.Background())
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.USERNAME: "bob",
			constants.MOVIE_ID: "5",
		},
		Context: ctx,
	}

	res, err := gql.UpdateCartField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, constants.SUCCESS, res)
}

func TestUpdateCartField_Remove_NotModified(t *testing.T) {
	mockSvc := setupTestMemberService()
	mockSvc.On("RemoveFromCart", mock.Anything, "bob", "5").Return(false, nil)

	ctx := context.WithValue(context.Background(), gql.GinContextKey, context.Background())
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.USERNAME:         "bob",
			constants.MOVIE_ID:         "5",
			constants.REMOVE_FROM_CART: true,
		},
		Context: ctx,
	}

	res, err := gql.UpdateCartField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, "did not modify cart for bob", res)
}

func TestCheckoutField(t *testing.T) {
	mockSvc := setupTestMemberService()
	mockSvc.On("Checkout", mock.Anything, "carol", []string{"3"}).Return([]string{"Checked out 3"}, 1, nil)

	ctx := context.WithValue(context.Background(), gql.GinContextKey, context.Background())
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.USERNAME:  "carol",
			constants.MOVIE_IDS: []interface{}{"3"},
		},
		Context: ctx,
	}

	res, err := gql.CheckoutField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Checked out 3"}, res)
}

func TestSetAPIChoiceField(t *testing.T) {
	mockSvc := setupTestMemberService()
	mockSvc.On("SetAPIChoice", mock.Anything, "dave", constants.GRAPHQL_API).Return(nil)

	ctx := context.WithValue(context.Background(), gql.GinContextKey, context.Background())
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.USERNAME:   "dave",
			constants.API_CHOICE: constants.GRAPHQL_API,
		},
		Context: ctx,
	}

	res, err := gql.SetAPIChoiceField.Resolve(params)
	assert.NoError(t, err)
	assert.Contains(t, res.(string), "successfully set dave api choice")
}
