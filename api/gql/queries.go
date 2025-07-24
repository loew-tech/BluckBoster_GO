package gql

import (
	"fmt"
	"net/http"
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
			return nil, getFormattedError(fmt.Sprintf("invalid page key: %s. Must be one of %s", page, constants.PAGES), http.StatusBadRequest)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		movies, err := movieService.GetMoviesByPage(ctx, page)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		movie, err := movieService.GetMovie(ctx, movieID)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		movies, err := memberService.GetCartMovies(ctx, username)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		user, err := memberService.GetMember(ctx, username, constants.CART)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
		}
		movies, err := movieService.GetMovies(ctx, user.Checkedout)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		member, err := memberService.GetMember(ctx, username, constants.NOT_CART)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError("'star' argument is required for starredIn query", http.StatusBadRequest)
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
			return nil, getFormattedError("'star' argument is required for starredWith query", http.StatusBadRequest)
		}
		return movieGraph.GetStarredWith(star), nil
	},
}
