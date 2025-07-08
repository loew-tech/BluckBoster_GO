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
	"blockbuster/api/repos"
)

var (
	movieRepo  = repos.NewMovieRepo(repos.GetDynamoClient())
	memberRepo = repos.NewMembersRepo(repos.GetDynamoClient(), repos.NewMovieRepo(repos.GetDynamoClient()))
	movieGraph *MovieGraph
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
		constants.MOVIES:          &graphql.Field{Type: graphql.NewList(MovieType)},
		constants.TOTAL_MOVIES:    &graphql.Field{Type: graphql.Int},
		constants.DIRECTORS:       &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_DIRECTORS: &graphql.Field{Type: graphql.Int},
	},
})

var usernameArg = &graphql.ArgumentConfig{Type: graphql.ID}
var starArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
var movieIDArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
var movieIDsArg = &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String), DefaultValue: []string{}}
var directorArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}

func getQueries() graphql.Fields {
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
				movies, err := movieRepo.GetMoviesByPage(ctx, constants.NOT_FOR_GRAPH, page)
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
				movie, err := movieRepo.GetMovieByID(ctx, movieID, constants.NOT_CART)
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve movie %s from cloud: %w", movieID, err)
					log.Println(errWrap)
					return nil, errWrap
				}
				return movie, nil
			},
		},
		constants.GET_CART: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: usernameArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for getCart query"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				movies, err := memberRepo.GetCartMovies(ctx, username)
				var errorWrap error
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve movies in cart: %w", err)
					log.Println(errWrap)
				}
				return movies, errorWrap
			},
		},
		constants.GET_CHECKEDOUT: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: usernameArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for getCheckedOut query"
					log.Println(msg)
					return nil, errors.New(msg)
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
				}
				return movies, err
			},
		},
		constants.GET_MEMBER: &graphql.Field{
			Type: MemberType,
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: usernameArg,
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
				member, err := memberRepo.GetMemberByUsername(ctx, username, constants.NOT_CART)
				if err != nil {
					errWrap := fmt.Errorf("failed to retrieve member %s from cloud: %w", username, err)
					log.Println(errWrap)
					return nil, errWrap
				}
				return member, nil
			},
		},
		constants.DIRECTED_MOVIES: &graphql.Field{
			Type: graphql.NewList(MovieType),
			Args: graphql.FieldConfigArgument{
				constants.DIRECTOR: directorArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				director := p.Args[constants.DIRECTOR].(string)
				return movieGraph.GetDirectedMovies(director), nil
			},
		},
		constants.DIRECTED_PERFORMERS: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.DIRECTOR: directorArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				director := p.Args[constants.DIRECTOR].(string)
				fmt.Println("Director=", director)
				return movieGraph.GetDirectedActors(director), nil
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

				depth := min(p.Args[constants.DEPTH].(int), 10)
				stars := make(map[string]bool)
				movieIDs := make(map[string]bool)
				directors := make(map[string]bool)
				for _, s := range toSearch {
					if _, found := stars[s]; found {
						continue
					}
					KevinBaconInOut(s, stars, movieIDs, directors, depth)
				}

				return map[string]interface{}{
					constants.STAR:            star,
					constants.STARS:           SetToList(stars),
					constants.TOTAL_STARS:     movieGraph.NumStars,
					constants.MOVIES:          movieGraph.GetMoviesByID(movieIDs),
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
		movie, err := movieRepo.GetMovieByID(ctx, movieID, constants.NOT_CART)
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
	movieIDs := make(map[string]bool)
	directors := make(map[string]bool)
	KevinBaconInOut(star, stars, movieIDs, directors, depth)
	return SetToList(stars), SetToList(movieIDs), SetToList(directors)
}

func KevinBaconInOut(star string, stars map[string]bool, movieIDs map[string]bool, directors map[string]bool, depth int) {
	movieGraph.BFS(star, stars, movieIDs, directors, depth)
}

func getMutations() graphql.Fields {
	return graphql.Fields{
		constants.RETURN_RENTALS: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.USERNAME:  usernameArg,
				constants.MOVIE_IDS: movieIDsArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for returnRentals mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				movieIDs, ok := p.Args[constants.MOVIE_IDS].([]interface{})
				if !ok || len(movieIDs) == 0 {
					msg := "movieIds argument is required for returnRentals mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				ids := make([]string, len(movieIDs))
				for i, v := range movieIDs {
					ids[i], _ = v.(string)
				}
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				messages, _, err := memberRepo.Return(ctx, username, ids)
				if err != nil {
					errWrap := fmt.Errorf("failed to return rentals for user %s: %w", username, err)
					log.Println(errWrap)
					return nil, errWrap
				}
				return messages, nil
			},
		},
		constants.UPDATE_CART: &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				constants.USERNAME:         usernameArg,
				constants.MOVIE_ID:         movieIDArg,
				constants.REMOVE_FROM_CART: &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for checkoutString mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				movieID, ok := p.Args[constants.MOVIE_ID].(string)
				if !ok || len(movieID) == 0 {
					msg := "movieIds argument is required for checkoutString mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				shouldRemoveFromCart, _ := p.Args[constants.REMOVE_FROM_CART].(bool)
				ctx, err := getContext(p)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				action := constants.ADD
				if shouldRemoveFromCart {
					action = constants.DELETE
				}
				act, direction := "adding", "to"
				if action == constants.DELETE {
					act, direction = "removing", "from"
				}
				inserted, _, err := memberRepo.ModifyCart(ctx, username, movieID, action, false)
				if err != nil {
					wrapErr := fmt.Errorf("error %s %s %s %s cart. Err: %w", act, movieID, direction, username, err)
					log.Println(wrapErr)
					return nil, wrapErr
				} else if !inserted {
					return fmt.Sprintf("Failed to %s from %s cart", movieID, username), nil
				}

				return "success", nil
			},
		},
		constants.CHECKOUT_STRING: &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Args: graphql.FieldConfigArgument{
				constants.USERNAME:  usernameArg,
				constants.MOVIE_IDS: movieIDsArg,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for checkout mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				movieIDs, ok := p.Args[constants.MOVIE_IDS].([]interface{})
				if !ok || len(movieIDs) == 0 {
					msg := "movieIds argument is required for checkout mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				ids := make([]string, len(movieIDs))
				for i, v := range movieIDs {
					ids[i], _ = v.(string)
				}
				ctx, err := getContext(p)
				if err != nil {
					return nil, err
				}
				messages, _, err := memberRepo.Checkout(ctx, username, ids)
				errWrap := err
				if err != nil {
					errWrap = fmt.Errorf("failed to checkout for user %s: %w", username, err)
					log.Println(errWrap)
				}
				return messages, errWrap
			},
		},
		constants.SET_API_CHOICE: &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				constants.USERNAME: usernameArg,
				constants.API_CHOICE: &graphql.ArgumentConfig{
					Type:         graphql.String,
					DefaultValue: constants.REST_API,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args[constants.USERNAME].(string)
				if !ok || username == "" {
					msg := "username argument is required for checkoutString mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				apiChoice, ok := p.Args[constants.API_CHOICE].(string)
				if !ok || apiChoice == "" {
					msg := "apiChoice argument is required for checkoutString mutation"
					log.Println(msg)
					return nil, errors.New(msg)
				}
				if apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API {
					msg := fmt.Sprintf("apiChoice must be either %s or %s, got %s", constants.REST_API, constants.GRAPHQL_API, apiChoice)
					log.Println(msg)
					return nil, errors.New(msg)
				}
				ctx, err := getContext(p)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				err = memberRepo.SetMemberAPIChoice(ctx, username, apiChoice)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				return fmt.Sprintf("successfully set %s api choice to %s", username, apiChoice), nil
			},
		},
	}
}

func getSchema() graphql.Schema {
	queryFiels := getQueries()
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFiels}
	mutationFields := getMutations()
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
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

var (
	ginContextKey       = contextKeyGin{}
	initMovieGraphOnce  sync.Once
	initMovieGraphError error
)

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
