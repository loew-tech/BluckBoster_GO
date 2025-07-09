package graphsearch

import "blockbuster/api/data"

type MovieGraph struct {
	directedMovies map[string][]data.Movie
	starredWith    map[string]map[string]bool
	starredIn      map[string][]data.Movie
	NumDirectors   int
	NumStars       int
	NumMovies      int
	movieIDToMovie map[string]data.Movie
}

// BFS traverses the graph starting from an actor and collects related stars, movies, and directors.
func (g *MovieGraph) BFS(
	startStar string,
	stars map[string]bool,
	movieIDs map[string]bool,
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

				if movieIDs[movie.ID] {
					continue
				}
				movieIDs[movie.ID] = true

				for _, coStar := range movie.Cast {
					if !stars[coStar] {
						nextSearch = append(nextSearch, coStar)
					}
				}
			}
		}
		toSearch = nextSearch
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

func (g *MovieGraph) GetMoviesByID(movieIDs map[string]bool) []data.Movie {
	movies := make([]data.Movie, 0, len(movieIDs))
	for id := range movieIDs {
		if movie, ok := g.movieIDToMovie[id]; ok {
			movies = append(movies, movie)
		}
	}
	return movies
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
