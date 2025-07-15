package gql

import (
	"errors"
	"log"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/data"
)

var GetKevinBaconField = &graphql.Field{
	Type: KevingBaconType,
	Args: graphql.FieldConfigArgument{
		constants.STAR:     starArg,
		constants.TITLE:    movieTitleArg,
		constants.DIRECTOR: directorArg,
		constants.DEPTH:    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		star := p.Args[constants.STAR].(string)
		movieTitle := p.Args[constants.TITLE].(string)
		director := p.Args[constants.DIRECTOR].(string)
		toSearch := buildToSearch(p, star, movieTitle, director)
		if len(toSearch) == 0 {
			msg := "the KevinBacon search requires at least one star, title, or director"
			log.Println(msg)
			return nil, errors.New(msg)
		}

		depth := min(p.Args[constants.DEPTH].(int), 10)
		stars := make(map[string]bool)
		movieTitles := make(map[string]bool)
		directors := make(map[string]bool)
		for _, s := range toSearch {
			if _, found := stars[s]; found {
				continue
			}
			KevinBaconInOut(s, stars, movieTitles, directors, depth)
		}

		var movies []data.Movie
		for title := range movieTitles {
			movie, err := movieGraph.GetMovieFromTitle(title)
			if err != nil {
				continue
			}
			movies = append(movies, movie)
			log.Printf("\t\tlen(movies)=%v\n", len(movies))
		}

		return map[string]interface{}{
			constants.STAR:            star,
			constants.STARS:           SetToList(stars),
			constants.TOTAL_STARS:     movieGraph.TotalStars(),
			constants.MOVIES:          movies,
			constants.TOTAL_MOVIES:    movieGraph.TotalMovies(),
			constants.DIRECTORS:       SetToList(directors),
			constants.TOTAL_DIRECTORS: movieGraph.TotalDirectors(),
		}, nil
	},
}

func KevinBacon(start string, depth int) ([]string, []string, []string) {
	stars := make(map[string]bool)
	movieTitles := make(map[string]bool)
	directors := make(map[string]bool)
	KevinBaconInOut(start, stars, movieTitles, directors, depth)
	return SetToList(stars), SetToList(movieTitles), SetToList(directors)
}

func KevinBaconInOut(star string, stars, movieTitles, directors map[string]bool, depth int) {
	movieGraph.BFS(star, stars, movieTitles, directors, depth)
}

func buildToSearch(p graphql.ResolveParams, star string, movieTitle string, director string) []string {
	var toSearch []string
	if star != "" {
		toSearch = append(toSearch, star)
	}
	if director != "" {
		toSearch = append(toSearch, movieGraph.GetDirectedActors(director)...)
	}
	if movieTitle != "" {
		movie, err := movieGraph.GetMovieFromTitle(movieTitle)
		if err != nil {
			return toSearch
		}
		toSearch = append(toSearch, movie.Cast...)
	}
	return toSearch
}
