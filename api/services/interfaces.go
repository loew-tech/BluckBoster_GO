package services

import (
	"context"

	"blockbuster/api/data"
)

type MembersServiceInterface interface {
	GetMember(ctx context.Context, username string, forCart bool) (data.Member, error)
	Login(ctx context.Context, username string) (data.Member, error)
	GetCartIDs(ctx context.Context, username string) ([]string, error)
	GetCartMovies(ctx context.Context, username string) ([]data.Movie, error)
	AddToCart(ctx context.Context, username, movieID string) (bool, error)
	RemoveFromCart(ctx context.Context, username, movieID string) (bool, error)
	Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error)
	SetAPIChoice(ctx context.Context, username, apiChoice string) error
	GetIniitialVotingSlate(ctx context.Context) ([]string, error)
	IterateRecommendationVoting(ctx context.Context, currentMood data.MovieMetrics, iteration, numPrevSelected int, movieIDs []string) (data.MovieMetrics, []string, error)
	GetVotingFinalPicks(ctx context.Context, mood data.MovieMetrics) ([]string, error)
	UpdateMood(ctx context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, error)
}

type MoviesServiceInterface interface {
	GetMoviesByPage(ctx context.Context, page string) ([]data.Movie, error)
	GetMovie(ctx context.Context, movieID string) (data.Movie, error)
	GetMovies(ctx context.Context, movieIDs []string) ([]data.Movie, error)
	GetMovieMetrics(ctx context.Context, movieID string) (data.MovieMetrics, error)
	GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error)
}
