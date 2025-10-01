package repos

import (
	"blockbuster/api/data"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoClientInterface interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
}

type MovieReadRepo interface {
	GetMoviesByPage(ctx context.Context, page string, purpose string) ([]data.Movie, error)
	GetMovieByID(ctx context.Context, movieID string, forCart bool) (data.Movie, error)
	GetMoviesByID(ctx context.Context, movieIDs []string, forCart bool) ([]data.Movie, error)
	GetTrivia(ctx context.Context, movieID string) (data.MovieTrivia, error)
	GetMovieMetrics(ctx context.Context, movieID string) (data.MovieMetrics, error)
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
	IterateRecommendationVoting(ctx context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, []string, error)
	UpdateMood(ctx context.Context, currentMood data.MovieMetrics, iteration int, movieIDs []string) (data.MovieMetrics, error)
}
