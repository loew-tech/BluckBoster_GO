package gql

import (
	"context"
	"errors"
	"log"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/endpoints"
	"blockbuster/api/repos"
)

type MovieGraph struct {
	directed     map[string][]data.Movie
	starredWith  map[string]map[string]bool
	starredIn    map[string][]data.Movie
	NumDirectors int
	NumStars     int
	NumMovies    int
	aveEdgeNum   int
}

var (
	mg      *MovieGraph
	once    sync.Once
	initErr error
)

func NewMovieGraph() (*MovieGraph, error) {
	once.Do(func() {
		mg = &MovieGraph{
			directed:     make(map[string][]data.Movie),
			starredWith:  make(map[string]map[string]bool),
			starredIn:    make(map[string][]data.Movie),
			NumDirectors: 0,
			NumStars:     0,
			NumMovies:    0,
		}
		initErr = populateCaches(mg)
	})
	return mg, initErr
}

func populateCaches(g *MovieGraph) error {
	var errs []error
	movieRepo := repos.NewMovieRepo(endpoints.GetDynamoClient())
	ctx := context.Background()
	for _, page := range constants.PAGES {
		movies, err := movieRepo.GetMoviesByPage(ctx, string(page))
		if err != nil {
			log.Printf("Err fetching movies for page %v. Err: %v\n", page, err)
			errs = append(errs, err)
			continue
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
	}
	g.NumDirectors = len(g.directed)
	g.NumStars = len(g.starredIn)
	err := getAverageStarredWithSize(g)
	if err != nil {
		log.Printf("Error calculating average starredWith with size: %v\n", err)
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func getAverageStarredWithSize(g *MovieGraph) error {
	total := 0
	if g.NumStars == 0 {
		return errors.New("no stars found in the graph")
	}
	for _, coStars := range g.starredWith {
		total += len(coStars)
	}
	g.aveEdgeNum = total / g.NumStars
	return nil
}

func (g *MovieGraph) BFS(startStar string, stars map[string]bool, movies map[string]bool, directors map[string]bool, maxDepth int) {
	toSearch := []string{startStar}
	depth := 0
	for len(toSearch) > 0 && depth < maxDepth {
		depth++
		nextSearch := make([]string, 0, g.NumStars)
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
	var actorList []string
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
	var coStars []string
	for coStar := range g.starredWith[star] {
		coStars = append(coStars, coStar)
	}
	return coStars
}
