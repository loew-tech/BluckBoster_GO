package gql_test

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/gql"
	graphsearch "blockbuster/api/graph_search"
	"blockbuster/api/services"
)

var (
	mockMemberService = new(services.MockMembersService)
	mockMovieService  = new(services.MockMoviesService)
	mockMovieGraph    = new(graphsearch.MockMovieGraph)
)

func setupTestContext() context.Context {
	return context.WithValue(context.Background(), gql.GinContextKey, context.Background())
}

func TestGetMoviesField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	mockMovieService.On("GetMoviesByPage", mock.Anything, constants.DEFAULT_PAGE).Return([]data.Movie{{ID: "1", Title: "Test Movie"}}, nil)
	gql.SetMovieGraph(mockMovieGraph)

	params := graphql.ResolveParams{
		Args:    map[string]interface{}{constants.PAGE: constants.DEFAULT_PAGE},
		Context: setupTestContext(),
	}

	resp, err := gql.GetMoviesField.Resolve(params)
	assert.NoError(t, err)
	assert.Len(t, resp.([]data.Movie), 1)
}

func TestGetMovieField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	mockMovieService.On("GetMovie", mock.Anything, mock.Anything).Return(data.Movie{ID: "1", Title: "Test Movie"}, nil)
	gql.SetMovieGraph(mockMovieGraph)

	params := graphql.ResolveParams{
		Args:    map[string]interface{}{constants.MOVIE_ID: "1"},
		Context: setupTestContext(),
	}
	resp, err := gql.GetMovieField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.(data.Movie).ID)
}

func TestGetCartField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	mockMemberService.On("GetCartMovies", mock.Anything, "test").Return([]data.Movie{{ID: "1", Title: "Test Movie"}}, nil)
	gql.SetMovieService(mockMovieService)
	gql.SetMovieGraph(mockMovieGraph)

	params := graphql.ResolveParams{
		Args:    map[string]interface{}{constants.USERNAME: "test"},
		Context: setupTestContext(),
	}
	resp, err := gql.GetCartField.Resolve(params)
	assert.NoError(t, err)
	assert.Len(t, resp.([]data.Movie), 1)
}

func TestGetCheckedOutField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	mockMemberService.On("GetMember", mock.Anything, mock.Anything, mock.Anything).Return(&data.Member{Username: "bob"}, nil)
	mockMemberService.On("GetCheckedOutMovies", mock.Anything, "bob").Return([]data.Movie{{ID: "1", Title: "Test Movie"}}, nil)
	gql.SetMovieService(mockMovieService)
	mockMovieService.On("GetMovies", mock.Anything, mock.Anything).Return([]data.Movie{{ID: "1", Title: "Test Movie"}}, nil)
	gql.SetMovieGraph(mockMovieGraph)

	params := graphql.ResolveParams{
		Args:    map[string]interface{}{constants.USERNAME: "bob"},
		Context: setupTestContext(),
	}
	resp, err := gql.GetCheckedOutField.Resolve(params)
	assert.NoError(t, err)
	assert.Len(t, resp.([]data.Movie), 1)
}

func TestGetMemberField(t *testing.T) {
	params := graphql.ResolveParams{
		Args:    map[string]interface{}{constants.USERNAME: "bob"},
		Context: setupTestContext(),
	}
	resp, err := gql.GetMemberField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, "bob", resp.(data.Member).Username)
}

func TestGetDirectedMoviesField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	gql.SetMovieGraph(mockMovieGraph)
	mockMovieGraph.On("GetDirectedMovies", mock.Anything).Return([]data.Movie{{ID: "1", Title: "Jaws", Director: "Spielberg"}}, nil)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{constants.DIRECTOR: "Spielberg"},
	}
	resp, err := gql.GetDirectedMoviesField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, "Spielberg", resp.([]data.Movie)[0].Director)
}

func TestGetDirectedActorsField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	gql.SetMovieGraph(mockMovieGraph)
	mockMovieGraph.On("GetDirectedActors", mock.Anything).Return([]string{"Actor A", "Actor B"}, nil)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{constants.DIRECTOR: "Spielberg"},
	}
	resp, err := gql.GetDirectedActorsField.Resolve(params)
	assert.NoError(t, err)
	assert.Contains(t, resp.([]string), "Actor A")
}

func TestGetStarredInField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	gql.SetMovieGraph(mockMovieGraph)
	mockMovieGraph.On("GetStarredIn", mock.Anything).Return([]data.Movie{{ID: "1", Title: "Jaws", Director: "Spielberg", Cast: []string{"Tom Hanks"}}}, nil)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{constants.STAR: "Tom Hanks"},
	}
	resp, err := gql.GetStarredInField.Resolve(params)
	assert.NoError(t, err)
	assert.Equal(t, "Tom Hanks", resp.([]data.Movie)[0].Cast[0])
}

func TestGetStarredWithField(t *testing.T) {
	gql.SetMemberService(mockMemberService)
	gql.SetMovieService(mockMovieService)
	gql.SetMovieGraph(mockMovieGraph)
	mockMovieGraph.On("GetStarredWith", mock.Anything).Return([]string{"Actor A", "Actor B"}, nil)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{constants.STAR: "Tom Hanks"},
	}
	resp, err := gql.GetStarredWithField.Resolve(params)
	assert.NoError(t, err)
	assert.Contains(t, resp.([]string), "Actor B")
}
