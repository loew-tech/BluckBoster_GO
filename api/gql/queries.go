package gql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/utils"
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
			return nil, fmt.Errorf("invalid page key: %s. Must be one of %s", page, constants.PAGES)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		movies, err := movieRepo.GetMoviesByPage(ctx, page, constants.NOT_FOR_GRAPH)
		if err != nil {
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve movies for page %s", page), err)
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
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve movie %s from cloud", movieID), err)
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
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve movies in cart for user %s", username), err)
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
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve user %s", username), err)
		}
		movies, err := movieRepo.GetMoviesByID(ctx, user.Checkedout, constants.CART)
		if err != nil {
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve checked out movies for user %s", username), err)
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
			return nil, utils.LogError(fmt.Sprintf("failed to retrieve member %s from cloud", username), err)
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
			return nil, utils.LogError("star argument is required for starredIn query", errors.New("missing star argument"))
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
			return nil, utils.LogError("star argument is required for starredWith query", errors.New("missing star argument"))
		}
		return movieGraph.GetStarredWith(star), nil
	},
}
