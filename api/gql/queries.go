// File: queries.go
package gql

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/data"
)

var GetMoviesField = &graphql.Field{
	Type: graphql.NewList(MovieType),
	Args: graphql.FieldConfigArgument{
		constants.PAGE: &graphql.ArgumentConfig{
			Type:         graphql.String,
			DefaultValue: constants.DEFAULT_PAGE,
		},
		constants.DIRECTOR: directorArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		page := p.Args[constants.PAGE].(string)
		if !strings.Contains(constants.PAGES, page) {
			return nil, fmt.Errorf("%s is not a valid page: %s", page, constants.PAGES)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		movies, err := movieRepo.GetMoviesByPage(ctx, page, constants.NOT_FOR_GRAPH)
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
}

var GetMovieField = &graphql.Field{
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
		movie, err := movieRepo.GetMovieByID(ctx, movieID, constants.NOT_CART)
		if err != nil {
			errWrap := fmt.Errorf("failed to retrieve movie %s from cloud: %w", movieID, err)
			log.Println(errWrap)
			return nil, errWrap
		}
		return movie, nil
	},
}

var GetCartField = &graphql.Field{
	Type: graphql.NewList(MovieType),
	Args: graphql.FieldConfigArgument{
		constants.USERNAME: usernameArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "getCart")
		if err != nil {
			return nil, err
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		movies, err := memberRepo.GetCartMovies(ctx, username)
		if err != nil {
			errWrap := fmt.Errorf("failed to retrieve movies in cart: %w", err)
			log.Println(errWrap)
			return nil, errWrap
		}
		return movies, nil
	},
}

var GetCheckedOutField = &graphql.Field{
	Type: graphql.NewList(MovieType),
	Args: graphql.FieldConfigArgument{
		constants.USERNAME: usernameArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "getCheckedOut")
		if err != nil {
			return nil, err
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		user, err := memberRepo.GetMemberByUsername(ctx, username, constants.CART)
		if err != nil {
			errWrap := fmt.Errorf("failed to retrieve user  %s: %w", username, err)
			log.Println(errWrap)
			return nil, errWrap
		}
		movies, err := movieRepo.GetMoviesByID(ctx, user.Checkedout, constants.CART)
		if err != nil {
			errWrap := fmt.Errorf("failed to retrieve movies: %w", err)
			log.Println(errWrap)
			return nil, errWrap
		}
		return movies, nil
	},
}

var GetMemberField = &graphql.Field{
	Type: MemberType,
	Args: graphql.FieldConfigArgument{
		constants.USERNAME: usernameArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "getMember")
		if err != nil {
			return nil, err
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		member, err := memberRepo.GetMemberByUsername(ctx, username, constants.NOT_CART)
		if err != nil {
			errWrap := fmt.Errorf("failed to retrieve member %s from cloud: %w", username, err)
			log.Println(errWrap)
			return nil, errWrap
		}
		return member, nil
	},
}

var GetDirectedMoviesField = &graphql.Field{
	Type: graphql.NewList(MovieType),
	Args: graphql.FieldConfigArgument{
		constants.DIRECTOR: directorArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		director := p.Args[constants.DIRECTOR].(string)
		return movieGraph.GetDirectedMovies(director), nil
	},
}

var GetDirectedActorsField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
	Args: graphql.FieldConfigArgument{
		constants.DIRECTOR: directorArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		director := p.Args[constants.DIRECTOR].(string)
		return movieGraph.GetDirectedActors(director), nil
	},
}

var GetStarredInField = &graphql.Field{
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
}

var GetStarredWithField = &graphql.Field{
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
}
