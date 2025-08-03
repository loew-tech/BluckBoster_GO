package gql_test

import (
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/gql"
	graphsearch "blockbuster/api/graph_search"
)

func TestGetKevinBaconField_Resolve_Valid(t *testing.T) {
	mockGraph := new(graphsearch.MockMovieGraph)
	gql.SetMovieGraph(mockGraph)

	mockGraph.On("GetDirectedActors", "Herbert Ross").Return([]string{"Kevin Bacon"})
	mockGraph.On("GetMovieFromTitle", "Footloose").Return(data.Movie{ID: "1", Title: "Footloose"}, nil)
	mockGraph.On("BFS", "Kevin Bacon", mock.Anything, mock.Anything, mock.Anything, 2).Return()
	mockGraph.On("TotalStars").Return(1)
	mockGraph.On("TotalMovies").Return(1)
	mockGraph.On("TotalDirectors").Return(1)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.STAR:     "Kevin Bacon",
			constants.TITLE:    "Footloose",
			constants.DIRECTOR: "Herbert Ross",
			constants.DEPTH:    2,
		},
	}

	result, err := gql.GetKevinBaconField.Resolve(params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	resp := result.(map[string]interface{})
	assert.Equal(t, "Kevin Bacon", resp[constants.STAR])
	assert.ElementsMatch(t, []string{"Kevin Bacon"}, resp[constants.STARS])
	assert.ElementsMatch(t, []string{"Footloose"}, extractTitles(resp[constants.MOVIES].([]data.Movie)))
	assert.ElementsMatch(t, []string{"Herbert Ross"}, resp[constants.DIRECTORS])
}

func TestGetKevinBaconField_Resolve_EmptyArgs(t *testing.T) {
	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.STAR:     "",
			constants.TITLE:    "",
			constants.DIRECTOR: "",
		},
	}

	result, err := gql.GetKevinBaconField.Resolve(params)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires at least one star")
}

func TestGetKevinBaconField_Resolve_MovieLookupFails(t *testing.T) {
	mockGraph := new(graphsearch.MockMovieGraph)
	gql.SetMovieGraph(mockGraph)

	mockGraph.On("GetDirectedActors", "Spielberg").Return([]string{"Actor A"})
	mockGraph.On("GetMovieFromTitle", "Nonexistent Movie").Return(data.Movie{}, errors.New("not found"))
	mockGraph.On("BFS", "Actor A", mock.Anything, mock.Anything, mock.Anything, 1).Return()
	mockGraph.On("TotalStars").Return(1)
	mockGraph.On("TotalMovies").Return(1)
	mockGraph.On("TotalDirectors").Return(1)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			constants.STAR:     "",
			constants.TITLE:    "Nonexistent Movie",
			constants.DIRECTOR: "Spielberg",
			constants.DEPTH:    1,
		},
	}

	result, err := gql.GetKevinBaconField.Resolve(params)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	resp := result.(map[string]interface{})
	assert.Equal(t, "", resp[constants.STAR])
}

func extractTitles(movies []data.Movie) []string {
	titles := make([]string, len(movies))
	for i, m := range movies {
		titles[i] = m.Title
	}
	return titles
}
