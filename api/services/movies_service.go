package services

import (
	"context"
	"fmt"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

type MoviesService struct {
	repo repos.ReadWriteMovieRepo
}

var (
	instantiateMovieServiceOnce sync.Once
	moviesService               *MoviesService
)

func GetMovieService() *MoviesService {
	instantiateMovieServiceOnce.Do(func() {
		moviesService = &MoviesService{repo: repos.NewMovieRepoWithDynamo()}
	})
	return moviesService
}

func (s *MoviesService) GetMoviesByPage(c context.Context, page string) ([]data.Movie, error) {
	movies, err := s.repo.GetMoviesByPage(c, page, constants.NOT_FOR_GRAPH)
	if err != nil {
		utils.LogError("", err)
		return nil, fmt.Errorf("failed to retrieve movies for page %s", page)
	}
	return movies, nil
}

func (s *MoviesService) GetMovie(c context.Context, movieID string) (data.Movie, error) {
	movie, err := s.repo.GetMovieByID(c, movieID, constants.NOT_CART)
	if err != nil {
		utils.LogError("err fetching movie", err)
		return data.Movie{}, fmt.Errorf("failed to fetch movie with id %s", movieID)
	}
	return movie, nil
}

func (s *MoviesService) GetTrivia(c context.Context, movieID string) (data.MovieTrivia, error) {
	trivia, err := s.repo.GetTrivia(c, movieID)
	if err != nil {
		return data.MovieTrivia{}, utils.LogError("failed to get trivia", err)
	}
	return trivia, nil
}
