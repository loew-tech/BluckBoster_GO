package gql

import (
	"log"

	"github.com/gin-gonic/gin"
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

func getFields() graphql.Fields {
	return graphql.Fields{
		constants.GET_MOVIES: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.PAGE: &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, ok := p.Args[constants.PAGE].(string)
				if !ok {
					page = "A"
				}
				return movieRepo.GetMoviesByPage(page)
			},
		},
		constants.GET_MOVIE: &graphql.Field{
			Type: MovieType,
			Args: graphql.FieldConfigArgument{
				constants.MOVIE_ID: &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movieID := p.Args[constants.MOVIE_ID].(string)
				movie, _, err := movieRepo.GetMovieByID(movieID, constants.NOT_CART)
				if err != nil {
					log.Fatalf("Failed to retrieve movie with ID %s from cloud. Err: %s\n", movieID, err)
				}
				return movie, nil
			},
		},
		constants.GET_MEMBER: &graphql.Field{
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
}

func getSchema() graphql.Schema {
	fields := getFields()
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}

func GetGQLHandler() func(*gin.Context) {

	schema := getSchema()
	gqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	return func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	}
}
