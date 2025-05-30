package gql

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"blockbuster/api/constants"
	"blockbuster/api/endpoints"
	"blockbuster/api/repos"
)

var movieRepo = repos.NewMovieRepo(endpoints.GetDynamoClient())
var membersRepo = repos.NewMembersRepo(endpoints.GetDynamoClient())

var MovieType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MOVIE_TYPE,
	Fields: graphql.Fields{
		constants.ID:        &graphql.Field{Type: graphql.String},
		constants.INVENTORY: &graphql.Field{Type: graphql.Int},
		constants.RATING:    &graphql.Field{Type: graphql.String},
		constants.RENTED:    &graphql.Field{Type: graphql.Int},
		constants.YEAR:      &graphql.Field{Type: graphql.String},
		constants.CAST:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.DIRECTOR:  &graphql.Field{Type: graphql.String},
		constants.TITLE:     &graphql.Field{Type: graphql.String},
	},
})

var MemberType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MEMBER_TYPE,
	Fields: graphql.Fields{
		constants.USERNAME:    &graphql.Field{Type: graphql.String},
		constants.FIRSTNAME:   &graphql.Field{Type: graphql.Int},
		constants.LASTNAME:    &graphql.Field{Type: graphql.String},
		constants.CART_STRING: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.CHECKED_OUT: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.RENTED:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.TYPE:        &graphql.Field{Type: graphql.String},
	},
})

func Foo() {
	fields := graphql.Fields{
		"GetMovies": &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args["page"].(string)
				if !ok {
					page = "A"
				}
				return movieRepo.GetMoviesByPage(page)
			},
		},
		"GetMovie": &graphql.Field{
			Type: MovieType,
			Args: graphql.FieldConfigArgument{
				"movieID": &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movieID := p.Args["movieID"].(string)
				movie, _, err := movieRepo.GetMovieByID(movieID, constants.NOT_CART)
				if err != nil {
					log.Fatalf("Failed to retrieve movie with ID %s from cloud. Err: %s\n", movieID, err)
				}
				return movie, nil
			},
		},
		"GetMember": &graphql.Field{
			Type: MemberType,
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username := p.Args[constants.USERNAME].(string)
				_, member, err := membersRepo.GetMemberByUsername(username, false)
				if err != nil {
					log.Fatalf("Failed to retrieve member from cloud. err: %s", err)
				}
				return member, nil
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
		query movies($page: String) {
			GetMovies(page: $page){
				id
				inventory
			}
		}
	`

	variableValues := map[string]interface{}{
		"page": "Z",
	}

	params := graphql.Params{Schema: schema, RequestString: query, VariableValues: variableValues}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
	fmt.Print("\n======\n\n")

	queryMember := `
		query Member($username: ID) {
			GetMember(username: $username) {
				username
				last_name
				member_type
				cart
				checked_out
			}
		}
	`

	variableValuesMember := map[string]interface{}{
		"username": "sea_captain",
	}
	paramsMember := graphql.Params{Schema: schema, RequestString: queryMember, VariableValues: variableValuesMember}
	rMember := graphql.Do(paramsMember)
	if len(rMember.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", rMember.Errors)
	}
	rJSONMember, _ := json.Marshal(rMember)
	fmt.Printf("2.\n%s \n", rJSONMember)

	handler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	http.Handle("/graphql", handler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("\nGoob Bye GQL")
}
