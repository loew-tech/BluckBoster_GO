package gql

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/endpoints"
	"blockbuster/api/repos"
)

var movieRepo = repos.NewMovieRepo(endpoints.GetDynamoClient())

var MovieType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MOVIE_TYPE,
	Fields: graphql.Fields{
		constants.ID:        &graphql.Field{Type: graphql.String},
		constants.INVENTORY: &graphql.Field{Type: graphql.Int},
		constants.RATING:    &graphql.Field{Type: graphql.String},
		constants.RENTED:    &graphql.Field{Type: graphql.Int},
		constants.YEAR:      &graphql.Field{Type: graphql.String},
		constants.CAST:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.TITLE:     &graphql.Field{Type: graphql.String},
	},
})

func Foo() {
	fields := graphql.Fields{
		constants.TITLE: &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movies, _ := movieRepo.GetMoviesByPage("A")
				return movies, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	query := `
	{
		title
	}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
	fmt.Println("\nGoob Bye GQL")
}
