package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"blockbuster/api/constants"
	"blockbuster/api/data"
	"blockbuster/api/repos"
	"blockbuster/api/utils"
)

type MoviesService struct {
	repo repos.MovieReadRepo
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

func NewMovieserviceWithRepo(repo repos.MovieReadRepo) *MoviesService {
	return &MoviesService{repo: repo}
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

func (s *MoviesService) GetMovies(c context.Context, movieIDs []string) ([]data.Movie, error) {
	movies, err := s.repo.GetMoviesByID(c, movieIDs, constants.CART)
	if err != nil {
		utils.LogError("err fetching movies", err)
		return nil, errors.New("failed to fetch movie")
	}
	return movies, nil
}

func (s *MoviesService) GetMovieMetrics(c context.Context, movieID string) (data.MovieMetrics, error) {
	metrics, err := s.repo.GetMovieMetrics(c, movieID)
	if err != nil {
		utils.LogError("failed to get movie metrics", err)
		return data.MovieMetrics{}, fmt.Errorf("failed to retrieve metrics for %s", movieID)
	}
	return metrics, nil
}

func (s *MoviesService) GetTrivia(c context.Context, movieID string) (data.MovieTrivia, error) {
	trivia, err := s.repo.GetTrivia(c, movieID)
	if err != nil {
		utils.LogError("failed to get trivia", err)
		return data.MovieTrivia{}, fmt.Errorf("failed to retrieve trivia for %s", movieID)
	}
	return trivia, nil
}
