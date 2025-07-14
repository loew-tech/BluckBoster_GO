package repos

import (
	"blockbuster/api/data"
	"context"
)

type MovieReadRepo interface {
	GetMoviesByPage(ctx context.Context, forGraph bool, page string) ([]data.Movie, error)
	GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error)
	GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error)
	GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error)
}

type MovieInventoryRepo interface {
	Rent(ctx context.Context, movie data.Movie) (bool, error)
	Return(ctx context.Context, movie data.Movie) (bool, error)
}

type ReadWriteMovieRepo interface {
	MovieReadRepo
	MovieInventoryRepo
}

type MemberRepoInterface interface {
	GetMemberByUsername(ctx context.Context, username string, cartOnly bool) (data.Member, error)
	GetCartMovies(ctx context.Context, username string) ([]data.Movie, error)
	GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error)
	ModifyCart(ctx context.Context, username, movieID, updateKey string, checkingOut bool) (bool, error)
	Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	SetMemberAPIChoice(ctx context.Context, username, apiChoice string) error
}
