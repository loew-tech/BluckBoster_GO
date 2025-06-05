package gql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

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
var (
	cacheMu      sync.RWMutex
	DIRECTED     = make(map[string][]data.Movie)
	STARRED_WITH = make(map[string][]string)
	STARRED_IN   = make(map[string][]data.Movie)
)

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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				movies, err := movieRepo.GetMoviesByPage(ctx, page)
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve movies for page %s: %w", page, err)
					log.Println(errWrap)
					return nil, errWrap
				}
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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				movie, _, err := movieRepo.GetMovieByID(ctx, movieID, constants.NOT_CART)
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve movie %s from cloud: %w", movieID, err)
					log.Println(errWrap)
					return nil, errWrap
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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				member, err := membersRepo.GetMemberByUsername(ctx, username, constants.NOT_CART)
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve member %s from cloud: %w", username, err)
					log.Println(errWrap)
					return nil, errWrap
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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				if len(DIRECTED) == 0 {
					populateCaches(ctx)
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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				if len(DIRECTED) == 0 {
					populateCaches(ctx)
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
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				if len(DIRECTED) == 0 {
					populateCaches(ctx)
				}
				return STARRED_WITH[star], nil
			},
		},
		constants.KEVING_BACON: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.STAR:  &graphql.ArgumentConfig{Type: graphql.String},
				constants.DEPTH: &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				depth := min(p.Args[constants.DEPTH].(int), 6)
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				if len(DIRECTED) == 0 {
					populateCaches(ctx)
				}

				stars := make(map[string]bool)
				movies := make(map[string]bool)
				directors := make(map[string]bool)

				bfs(star, stars, movies, directors, depth)

				result := make([]string, 0, len(stars))
				for s := range stars {
					if s != star {
						result = append(result, s)
					}
				}
				return result, nil
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

func getContext(p graphql.ResolveParams) (*gin.Context, error) {
	ctx, ok := p.Context.Value(ginContextKey).(*gin.Context)
	if !ok {
		log.Println("Gin context not found in resolve params")
		return nil, errors.New("missing Gin context")
	}
	return ctx, nil
}

func populateCaches(ctx context.Context) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if len(DIRECTED) != 0 && len(STARRED_IN) != 0 && len(STARRED_WITH) != 0 {
		return
	}

	for _, page := range PAGES {
		movies, err := movieRepo.GetMoviesByPage(ctx, page)
		if err != nil {
			log.Printf("Err fetching movies for page %s. Err: %s\n", page, err)
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

func bfs(startStar string, stars map[string]bool, movies map[string]bool, directors map[string]bool, depth int) {
	toSearch := []string{startStar}
	exploredDepth := 0
	for len(toSearch) > 0 && exploredDepth < depth {
		exploredDepth++
		nextSearch := make([]string, 0)
		for _, star := range toSearch {
			if _, found := stars[star]; found {
				continue
			}
			stars[star] = true

			for _, movie := range STARRED_IN[star] {
				if _, found := movies[movie.ID]; !found {
					movies[movie.ID] = true
					for _, coStar := range movie.Cast {
						if _, found := stars[coStar]; !found {
							stars[coStar] = true
							nextSearch = append(nextSearch, coStar)
						}
					}
				}
				directors[movie.Director] = true
			}
		}
		toSearch = nextSearch
	}
}

type ginContextKeyType struct{}

var ginContextKey = ginContextKeyType{}

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
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		corsHandler.Handler(gqlHandler).ServeHTTP(c.Writer, c.Request)
	}
}
