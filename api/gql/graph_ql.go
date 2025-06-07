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

var (
	initMovieGraphOnce  sync.Once
	movieRepo           = repos.NewMovieRepo(endpoints.GetDynamoClient())
	membersRepo         = repos.NewMembersRepo(endpoints.GetDynamoClient())
	movieGraph          *MovieGraph
	initMovieGraphError error
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
		constants.CHECKED_OUT: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.CART_STRING: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.RENTED:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.TYPE:        &graphql.Field{Type: graphql.String},
	},
})

var KevingBaconType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.KEVING_BACON_TYPE,
	Fields: graphql.Fields{
		constants.STAR:            &graphql.Field{Type: graphql.String},
		constants.STARS:           &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_STARS:     &graphql.Field{Type: graphql.Int},
		constants.MOVIES:          &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_MOVIES:    &graphql.Field{Type: graphql.Int},
		constants.DIRECTORS:       &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_DIRECTORS: &graphql.Field{Type: graphql.Int},
	},
})

var starArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
var movieIDArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
var directorArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}

func getFields() graphql.Fields {
	return graphql.Fields{
		constants.GET_MOVIES: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.PAGE: &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: "A",
				},
				constants.DIRECTOR: directorArg,
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
				if !ok || director == "" {
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
				constants.MOVIE_ID: movieIDArg,
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
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for getMember query"
					log.Println(msg)
					return nil, errors.New(msg)
				}
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
				constants.DIRECTOR: directorArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				director := p.Args[constants.DIRECTOR].(string)
				return movieGraph.GetDirectedMovies(director), nil
			},
		},
		constants.STARREDIN: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.STAR: starArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				if star == "" {
					msg := "star argument is required for starredIn query"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				return movieGraph.GetStarredIn(star), nil
			},
		},
		constants.STARREDWITH: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.STAR: starArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				if star == "" {
					msg := "star argument is required for starredWith query"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				return movieGraph.GetStarredWith(star), nil
			},
		},
		constants.KEVING_BACON: &graphql.Field{
			Type: KevingBaconType,
			Args: graphql.FieldConfigArgument{
				constants.STAR:     starArg,
				constants.MOVIE:    movieIDArg,
				constants.DIRECTOR: directorArg,
				constants.DEPTH:    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				star := p.Args[constants.STAR].(string)
				movieID := p.Args[constants.MOVIE].(string)
				director := p.Args[constants.DIRECTOR].(string)
				toSearch := buildToSearch(p, star, movieID, director)
				if len(toSearch) == 0 {
					msg := "the KevinBacon search requires at least one star, movie, or director"
					log.Println(msg)
					return nil, errors.New(msg)
				}

				depth := min(p.Args[constants.DEPTH].(int), 6)
				stars := make(map[string]bool)
				movies := make(map[string]bool)
				directors := make(map[string]bool)
				for _, s := range toSearch {
					if _, found := stars[s]; s != star && !found {
						stars[s] = true
						continue
					}
					KevinBaconInOut(s, stars, movies, directors, depth)
				}

				return map[string]interface{}{
					constants.STAR:            star,
					constants.STARS:           SetToList(stars),
					constants.TOTAL_STARS:     movieGraph.NumStars,
					constants.MOVIES:          SetToList(movies),
					constants.TOTAL_MOVIES:    movieGraph.NumMovies,
					constants.DIRECTORS:       SetToList(directors),
					constants.TOTAL_DIRECTORS: movieGraph.NumDirectors,
				}, nil
			},
		},
	}
}

func buildToSearch(p graphql.ResolveParams, star string, movieID string, director string) []string {
	var toSearch []string
	if star != "" {
		toSearch = append(toSearch, star)
	}
	if director != "" {
		toSearch = append(toSearch, movieGraph.GetDirectedActors(director)...)
	}
	if movieID != "" {
		ctx, err := getContext(p)
		if err != nil {
			log.Printf("failed to get context: %v", err)
			return toSearch
		}
		movie, _, err := movieRepo.GetMovieByID(ctx, movieID, constants.NOT_CART)
		if err != nil {
			log.Printf("failed to retrieve movie %s from cloud: %v", movieID, err)
			return toSearch
		}
		toSearch = append(toSearch, movie.Cast...)
	}
	return toSearch
}

func KevinBacon(star string, depth int) ([]string, []string, []string) {
	stars := make(map[string]bool)
	movies := make(map[string]bool)
	directors := make(map[string]bool)
	KevinBaconInOut(star, stars, movies, directors, depth)
	return SetToList(stars), SetToList(movies), SetToList(directors)
}

func KevinBaconInOut(star string, stars map[string]bool, movies map[string]bool, directors map[string]bool, depth int) {
	movieGraph.BFS(star, stars, movies, directors, depth)
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
		msg := "gin context not found in resolve params"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	return ctx, nil
}

type contextKeyGin struct{}

var ginContextKey = contextKeyGin{}

func GetGQLHandler() func(*gin.Context) {
	initMovieGraphOnce.Do(func() {
		movieGraph, initMovieGraphError = NewMovieGraph()
		if initMovieGraphError != nil {
			log.Printf("errors encountered while populating movie graph. Some KevinBacon functionality will be affected. See above for individual errors. %v\n", initMovieGraphError)
		}
	})

	schema := getSchema()
	gqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
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
