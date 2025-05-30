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
		constants.MOVIES: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args["page"].(string)
				if !ok {
					page = "A"
				}
				movies, _ := movieRepo.GetMoviesByPage(page)
				return movies, nil
			},
		},
		constants.MEMBER: &graphql.Field{
			Type: MemberType,
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username := p.Args[constants.USERNAME].(string)
				fmt.Printf("** Tick username= %s\n\n", username)
				_, member, err := membersRepo.GetMemberByUsername(username, false)
				if err != nil {
					log.Fatalf("Failed to retrieve member from cloud. err: %s", err)
				}
				fmt.Println(member.Username, member.Cart, member.Type)
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
			movies(page: $page){
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
		query member($username: ID) {
			member(username: $username) {
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
	fmt.Println("\nGoob Bye GQL")
}
