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
