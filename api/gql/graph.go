package gql

import (
	"context"
	"log"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/endpoints"
	"blockbuster/api/repos"
)

type MovieGraph struct {
	cacheMu      sync.RWMutex
	directed     map[string][]data.Movie
	starredWith  map[string]map[string]bool
	starredIn    map[string][]data.Movie
	NumDirectors int
	NumStars     int
	NumMovies    int
}

func NewMovieGraph() *MovieGraph {
	m := &MovieGraph{
		directed:     make(map[string][]data.Movie),
		starredWith:  make(map[string]map[string]bool),
		starredIn:    make(map[string][]data.Movie),
		NumDirectors: 0,
		NumStars:     0,
		NumMovies:    0,
	}
	m.populateCaches()
	return m
}

func (g *MovieGraph) populateCaches() {
	g.cacheMu.Lock()
	movieRepo := repos.NewMovieRepo(endpoints.GetDynamoClient())
	ctx := context.TODO()
	defer g.cacheMu.Unlock()
	if len(g.directed) != 0 && len(g.starredIn) != 0 && len(g.starredWith) != 0 {
		return
	}

	for _, page := range constants.PAGES {
		movies, err := movieRepo.GetMoviesByPage(ctx, string(page))
		if err != nil {
			log.Printf("Err fetching movies for page %s. Err: %s\n", page, err)
		}
		for _, movie := range movies {
			g.NumMovies++
			g.directed[movie.Director] = append(g.directed[movie.Director], movie)
			for _, actor := range movie.Cast {
				g.starredIn[actor] = append(g.starredIn[actor], movie)
				for _, coStar := range movie.Cast {
					if actor == coStar {
						continue
					}
					if _, found := g.starredWith[actor]; !found {
						g.starredWith[actor] = make(map[string]bool)
					}
					g.starredWith[actor][coStar] = true
				}
			}
		}
		g.NumDirectors = len(g.directed)
		g.NumStars = len(g.starredIn)
	}
}

func (g *MovieGraph) BFS(startStar string, stars map[string]bool, movies map[string]bool, directors map[string]bool, depth int) {
	toSearch := []string{startStar}
	exploredDepth := 0
	for len(toSearch) > 0 && exploredDepth < depth {
		exploredDepth++
		nextSearch := make([]string, 0)
		for _, star := range toSearch {
			stars[star] = true
			for _, movie := range g.starredIn[star] {
				if _, found := movies[movie.ID]; !found {
					movies[movie.ID] = true
					for _, coStar := range movie.Cast {
						if _, found := stars[coStar]; !found {
							stars[coStar] = true
							nextSearch = append(nextSearch, coStar)
						}
					}
				}
				directors[movie.Director] = true
			}
		}
		toSearch = nextSearch
	}
}

func (g *MovieGraph) GetDirectedMovies(director string) []data.Movie {
	return g.directed[director]
}

func (g *MovieGraph) GetDirectedActors(director string) []string {
	actors := make(map[string]bool)
	actorList := make([]string, 0, len(actors))
	for _, movie := range g.directed[director] {
		for _, actor := range movie.Cast {
			if _, found := actors[actor]; found {
				continue
			}
			actors[actor] = true
			actorList = append(actorList, actor)
		}
	}
	return actorList
}

func (g *MovieGraph) GetStarredIn(star string) []data.Movie {
	return g.starredIn[star]
}

func (g *MovieGraph) GetStarredWith(star string) []string {
	coStars := make([]string, 0)
	for coStar := range g.starredWith[star] {
		coStars = append(coStars, coStar)
	}
	return coStars
}
