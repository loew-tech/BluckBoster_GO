package graphsearch

import "blockbuster/api/data"

type MovieGraphInterface interface {
	BFS(start string, stars, movieIDs, directors map[string]bool, depth int)
	GetDirectedActors(director string) []string
	GetDirectedMovies(director string) []data.Movie
	GetStarredIn(star string) []data.Movie
	GetStarredWith(star string) []string
	GetMovieFromTitle(title string) (data.Movie, error)
	TotalStars() int
	TotalMovies() int
	TotalDirectors() int
}
