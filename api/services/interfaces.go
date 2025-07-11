package services

import (
	"context"

	"blockbuster/api/data"
)

// @TODO: remove status return?
type MembersServiceInterface interface {
	GetMember(c context.Context, username string) (int, data.Member, error)
	Login(c context.Context, username string) (int, data.Member, error)
	GetCartIDs(c context.Context, username string) (int, []string, error)
	GetCartMovies(c context.Context, username string) (int, []data.Movie, error)
	AddToCart(c context.Context, username, movieID string) (int, error)
	RemoveFromCart(c context.Context, username, movieID string) (int, error)
	Checkout(c context.Context, username string, movieIDs []string) (int, []string, int, error)
	Return(c context.Context, username string, movieIDs []string) (int, []string, int, error)
	GetCheckedOutMovies(c context.Context, username string) (int, []data.Movie, error)
	SetAPIChoice(c context.Context, username, apiChoice string) (int, string, error)
}
