// File: schema.go
package gql

import (
	"log"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
)

func getQueries() graphql.Fields {
	return graphql.Fields{
		constants.GET_MOVIES:          GetMoviesField,
		constants.GET_MOVIE:           GetMovieField,
		constants.GET_CART:            GetCartField,
		constants.GET_CHECKEDOUT:      GetCheckedOutField,
		constants.GET_MEMBER:          GetMemberField,
		constants.DIRECTED_MOVIES:     GetDirectedMoviesField,
		constants.DIRECTED_PERFORMERS: GetDirectedActorsField,
		constants.STARREDIN:           GetStarredInField,
		constants.STARREDWITH:         GetStarredWithField,
		constants.KEVING_BACON:        GetKevinBaconField,
	}
}

func getMutations() graphql.Fields {
	return graphql.Fields{
		constants.RETURN_RENTALS:  ReturnRentalsField,
		constants.UPDATE_CART:     UpdateCartField,
		constants.CHECKOUT_STRING: CheckoutField,
		constants.SET_API_CHOICE:  SetAPIChoiceField,
	}
}

func getSchema() graphql.Schema {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: getQueries(),
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: getMutations(),
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	return schema
}
