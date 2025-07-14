package services

import (
	"context"

	"blockbuster/api/data"
)

type MembersServiceInterface interface {
	GetMember(ctx context.Context, username string) (data.Member, error)
	Login(ctx context.Context, username string) (data.Member, error)
	GetCartIDs(ctx context.Context, username string) ([]string, error)
	GetCartMovies(ctx context.Context, username string) ([]data.Movie, error)
	AddToCart(ctx context.Context, username, movieID string) (bool, error)
	RemoveFromCart(ctx context.Context, username, movieID string) (bool, error)
	Checkout(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	Return(ctx context.Context, username string, movieIDs []string) ([]string, int, error)
	GetCheckedOutMovies(ctx context.Context, username string) ([]data.Movie, error)
	SetAPIChoice(ctx context.Context, username, apiChoice string) error
}
