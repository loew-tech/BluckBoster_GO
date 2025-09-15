package graphsearch

import (
	"fmt"
	"math/rand"

	"blockbuster/api/data"
	"blockbuster/api/utils"
)

type MovieGraph struct {
	directedMovies     map[string][]data.Movie
	starredWith        map[string]map[string]bool
	starredIn          map[string][]data.Movie
	NumDirectors       int
	NumStars           int
	NumMovies          int
	movieTitleToMovie  map[string]data.Movie
	centroidToMovieIDs map[int][]string
}

// BFS traverses the graph starting from an actor and collects related stars, movies, and directors.
func (g *MovieGraph) BFS(
	startStar string,
	stars map[string]bool,
	movieTitles map[string]bool,
	directors map[string]bool,
	maxDepth int,
) {
	toSearch := []string{startStar}
	depth := 0

	for len(toSearch) > 0 && depth < maxDepth {
		depth++
		nextSearch := make([]string, 0, g.NumStars)

		for _, star := range toSearch {
			if stars[star] {
				continue
			}
			stars[star] = true

			for _, movie := range g.starredIn[star] {
				directors[movie.Director] = true

				if movieTitles[movie.Title] {
					continue
				}
				movieTitles[movie.Title] = true

				for _, coStar := range movie.Cast {
					if !stars[coStar] {
						nextSearch = append(nextSearch, coStar)
					}
				}
			}
		}
		toSearch = nextSearch
	}

	// Populate cache with stars seen before depth limit reached
	for _, star := range toSearch {
		stars[star] = true
	}
}

func (g *MovieGraph) GetDirectedMovies(director string) []data.Movie {
	return g.directedMovies[director]
}

func (g *MovieGraph) GetDirectedActors(director string) []string {
	seen := make(map[string]bool)
	var actors []string

	for _, movie := range g.directedMovies[director] {
		for _, actor := range movie.Cast {
			if !seen[actor] {
				seen[actor] = true
				actors = append(actors, actor)
			}
		}
	}
	return actors
}

func (g *MovieGraph) GetStarredIn(star string) []data.Movie {
	return g.starredIn[star]
}

func (g *MovieGraph) GetStarredWith(star string) []string {
	var coStars []string
	for coStar := range g.starredWith[star] {
		coStars = append(coStars, coStar)
	}
	return coStars
}

func (g *MovieGraph) GetIDFromTitle(title string) (string, error) {
	if movie, ok := g.movieTitleToMovie[title]; ok {
		return movie.ID, nil
	}
	return "", utils.LogError(fmt.Sprintf("failed to get movieID for title %s", title), nil)
}

func (g *MovieGraph) GetMovieFromTitle(title string) (data.Movie, error) {
	if movie, ok := g.movieTitleToMovie[title]; ok {
		return movie, nil
	}
	return data.TestMovies[0], utils.LogError(fmt.Sprintf("failed to retrieve movie with  %s", title), nil)
}

func (g *MovieGraph) GetRandomMovieFromCentroid(centroid int) (string, error) {
	if movieIDs, ok := g.centroidToMovieIDs[centroid]; ok {
		return movieIDs[rand.Intn(len(movieIDs))], nil
	}
	return "", utils.LogError(fmt.Sprintf("failed to retrieve random movie from centroid %d", centroid), nil)
}

func (g *MovieGraph) TotalStars() int {
	return g.NumStars
}

func (g *MovieGraph) TotalMovies() int {
	return g.NumMovies
}

func (g *MovieGraph) TotalDirectors() int {
	return g.NumDirectors
}
