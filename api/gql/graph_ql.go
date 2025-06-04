package gql

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/endpoints"
	"blockbuster/api/repos"
)

var movieRepo = repos.NewMovieRepo(endpoints.GetDynamoClient())
var membersRepo = repos.NewMembersRepo(endpoints.GetDynamoClient())

var PAGES = []string{"#", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var DIRECTED = make(map[string][]data.Movie)
var STARRED_WITH = make(map[string][]string)
var STARRED_IN = make(map[string][]data.Movie)

var MovieType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MOVIE_TYPE,
	Fields: graphql.Fields{
		constants.ID:        &graphql.Field{Type: graphql.String},
		constants.INVENTORY: &graphql.Field{Type: graphql.Int},
		constants.RATING:    &graphql.Field{Type: graphql.String},
		constants.REVIEW:    &graphql.Field{Type: graphql.String},
		constants.RENTED:    &graphql.Field{Type: graphql.Int},
		constants.SYNOPSIS:  &graphql.Field{Type: graphql.String},
		constants.TRIVIA:    &graphql.Field{Type: graphql.String},
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
		constants.FIRSTNAME:   &graphql.Field{Type: graphql.String},
		constants.LASTNAME:    &graphql.Field{Type: graphql.String},
		constants.CART_STRING: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.CHECKED_OUT: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.RENTED:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.TYPE:        &graphql.Field{Type: graphql.String},
	},
})

func populateCaches(ctx context.Context) {
	if len(DIRECTED) != 0 && len(STARRED_IN) != 0 && len(STARRED_WITH) != 0 {
		return
	}
	for _, page := range PAGES {
		movies, err := movieRepo.GetMoviesByPage(ctx, page)
		if err != nil {
			log.Fatalf("Err fetching movies for page %s. Err: %s", page, err)
		}
		for _, movie := range movies {
			DIRECTED[movie.Director] = append(DIRECTED[movie.Director], movie)
			for _, actor := range movie.Cast {
				STARRED_IN[actor] = append(STARRED_IN[actor], movie)
				for _, coStar := range movie.Cast {
					if actor == coStar {
						continue
					}
					STARRED_WITH[actor] = append(STARRED_WITH[actor], coStar)
				}
			}
		}
	}
}

func getFields() graphql.Fields {
	return graphql.Fields{
		constants.GET_MOVIES: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.PAGE: &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "A",
				},
				constants.DIRECTOR: &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page := p.Args[constants.PAGE].(string)
				movies, _ := movieRepo.GetMoviesByPage(p.Context, page)
				director, ok := p.Args[constants.DIRECTOR]
				if !ok {
					return movies, nil
				}
				moviesDirected := make([]data.Movie, 0)
				for _, movie := range movies {
					if movie.Director == director {
						moviesDirected = append(moviesDirected, movie)
					}
				}
				return moviesDirected, nil
			},
		},
		constants.GET_MOVIE: &graphql.Field{
			Type: MovieType,
			Args: graphql.FieldConfigArgument{
				constants.MOVIE_ID: &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				movieID := p.Args[constants.MOVIE_ID].(string)
				movie, _, err := movieRepo.GetMovieByID(p.Context, movieID, constants.NOT_CART)
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
				_, member, err := membersRepo.GetMemberByUsername(p.Context, username, constants.NOT_CART)
				if err != nil {
					log.Fatalf("Failed to retrieve member from cloud. err: %s", err)
				}
				return member, nil
			},
		},
		constants.DIRECTED_BY: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.DIRECTOR: &graphql.ArgumentConfig{Type: graphql.ID},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				director := p.Args[constants.DIRECTOR].(string)
				if len(DIRECTED) == 0 {
					populateCaches(p.Context)
				}
				return DIRECTED[director], nil
			},
		},
		constants.STARREDIN: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.STAR: &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				if len(STARRED_IN) == 0 {
					populateCaches(p.Context)
				}
				return STARRED_IN[star], nil
			},
		},
		constants.STARREDWITH: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.STAR: &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				if len(STARRED_WITH) == 0 {
					populateCaches(p.Context)
				}
				return STARRED_WITH[star], nil
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

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // your frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	return func(c *gin.Context) {
		corsHandler.Handler(gqlHandler).ServeHTTP(c.Writer, c.Request)
	}
}
