// File: kevin_bacon.go
package gql

import (
	"errors"
	"log"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/repos"
)

var GetKevinBaconField = &graphql.Field{
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
			constants.TOTAL_STARS:     movieGraph.TotalStars(),
			constants.MOVIES:          movieGraph.GetMoviesByID(movieIDs),
			constants.TOTAL_MOVIES:    movieGraph.TotalMovies(),
			constants.DIRECTORS:       SetToList(directors),
			constants.TOTAL_DIRECTORS: movieGraph.TotalDirectors(),
		}, nil
	},
}

func KevinBacon(start string, depth int) ([]string, []string, []string) {
	stars := make(map[string]bool)
	movieIDs := make(map[string]bool)
	directors := make(map[string]bool)
	KevinBaconInOut(start, stars, movieIDs, directors, depth)
	return SetToList(stars), SetToList(movieIDs), SetToList(directors)
}

func KevinBaconInOut(star string, stars, movieIDs, directors map[string]bool, depth int) {
	movieGraph.BFS(star, stars, movieIDs, directors, depth)
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
		movie, err := repos.NewMovieRepoWithDynamo().GetMovieByID(ctx, movieID, constants.NOT_CART)
		if err != nil {
			log.Printf("failed to retrieve movie %s from cloud: %v", movieID, err)
			return toSearch
		}
		toSearch = append(toSearch, movie.Cast...)
	}
	return toSearch
}
