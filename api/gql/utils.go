package gql

import (
	"blockbuster/api/constants"
	"blockbuster/api/services"
	"blockbuster/api/utils"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

var (
	movieService  = services.GetMovieService()
	memberService = services.GetMemberService()
)

// getStringArg safely extracts a required string arg from the resolver params.
func getStringArg(p graphql.ResolveParams, argName string, field string) (string, error) {
	val, ok := p.Args[argName].(string)
	if !ok || val == "" {
		msg := fmt.Sprintf("%s argument is required for %s", argName, field)
		log.Println(msg)
		return "", errors.New(msg)
	}
	return val, nil
}

func getContext(p graphql.ResolveParams) (*gin.Context, error) {
	ctx, ok := p.Context.Value(ginContextKey).(*gin.Context)
	if !ok {
		msg := "gin context not found in resolve params"
		utils.LogError(msg, nil)
		return nil, errors.New(msg)
	}
	return ctx, nil
}

func SetToList[T comparable](set map[T]bool) []T {
	list := make([]T, 0, len(set))
	for item := range set {
		list = append(list, item)
	}
	return list
}

func getFormattedError(msg string, status int) gqlerrors.FormattedError {
	return gqlerrors.FormattedError{
		Message: msg,
		Extensions: map[string]interface{}{
			constants.CODE: status,
		},
	}
}
