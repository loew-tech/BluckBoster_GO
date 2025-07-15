package graphsearch

import (
	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
	"context"
	"errors"
	"log"
	"sync"
)

var (
	movieGraph     *MovieGraph
	initGraphOnce  sync.Once
	initGraphError error
)

func GetMovieGraph() (*MovieGraph, error) {
	initGraphOnce.Do(func() {
		movieGraph, initGraphError = newMovieGraph(populateCaches)
		if initGraphError != nil {
			utils.LogError("error(s) occurred initializing MovieGraph", initGraphError)
		}
	})
	return movieGraph, initGraphError
}

// NewMovieGraph creates a MovieGraph and populates it using the provided function.
// This is useful for tests or alternate data loaders.
func newMovieGraph(populate func(*MovieGraph) error) (*MovieGraph, error) {
	graph := &MovieGraph{
		directedMovies:    make(map[string][]data.Movie),
		starredWith:       make(map[string]map[string]bool),
		starredIn:         make(map[string][]data.Movie),
		movieTitleToMovie: make(map[string]data.Movie),
	}
	if err := populate(graph); err != nil {
		return nil, err
	}
	return graph, nil
}

// populateCaches loads movie data from the repo and builds internal lookup maps.
// It continues on partial errors and returns a combined error if any occur.
func populateCaches(g *MovieGraph) error {
	var errs []error

	movieRepo := repos.NewMovieRepoWithDynamo()
	ctx := context.Background()

	for _, page := range constants.PAGES {
		movies, err := movieRepo.GetMoviesByPage(ctx, string(page), constants.FOR_GRAPH)
		if err != nil {
			log.Printf("Error fetching movies for page %v: %v\n", page, err)
			errs = append(errs, err)
			continue
		}

		for _, movie := range movies {
			g.NumMovies++
			g.movieTitleToMovie[movie.Title] = movie

			// Index by director
			g.directedMovies[movie.Director] = append(g.directedMovies[movie.Director], movie)

			for _, star := range movie.Cast {
				g.starredIn[star] = append(g.starredIn[star], movie)
				g.NumStars++

				if _, ok := g.starredWith[star]; !ok {
					g.starredWith[star] = make(map[string]bool)
				}
				for _, coStar := range movie.Cast {
					if star != coStar {
						g.starredWith[star][coStar] = true
					}
				}
			}
		}
	}

	g.NumDirectors = len(g.directedMovies)
	if len(errs) > 0 {
		return utils.LogError("graphsearch", errors.Join(errs...))
	}
	return nil
}
