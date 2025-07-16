// File: mutations.go
package gql

import (
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
	"blockbuster/api/utils"
)

var ReturnRentalsField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
	Args: graphql.FieldConfigArgument{
		constants.USERNAME:  usernameArg,
		constants.MOVIE_IDS: movieIDsArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "returnRentals")
		if err != nil {
			return nil, err
		}
		ids := extractIDList(p.Args[constants.MOVIE_IDS])
		if len(ids) == 0 {
			return nil, utils.LogError("movieIds argument is required for returnRentals mutation", nil)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		messages, _, err := memberRepo.Return(ctx, username, ids)
		if err != nil {
			return nil, utils.LogError(fmt.Sprintf("failed to return rentals for user %s", username), err)
		}
		return messages, nil
	},
}

var UpdateCartField = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		constants.USERNAME:         usernameArg,
		constants.MOVIE_ID:         movieIDArg,
		constants.REMOVE_FROM_CART: &graphql.ArgumentConfig{Type: graphql.Boolean},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "updateCart")
		if err != nil {
			return nil, err
		}
		movieID, err := getStringArg(p, constants.MOVIE_ID, "updateCart")
		if err != nil {
			return nil, err
		}
		shouldRemoveFromCart, _ := p.Args[constants.REMOVE_FROM_CART].(bool)
		action := constants.ADD
		if shouldRemoveFromCart {
			action = constants.DELETE
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		inserted, err := memberRepo.ModifyCart(ctx, username, movieID, action, false)
		if err != nil {
			return nil, utils.LogError(fmt.Sprintf("error modifying cart for user %s", username), err)
		} else if !inserted {
			return fmt.Sprintf("failed to modify cart for %s", username), nil
		}
		return constants.SUCCESS, nil
	},
}

var CheckoutField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
	Args: graphql.FieldConfigArgument{
		constants.USERNAME:  usernameArg,
		constants.MOVIE_IDS: movieIDsArg,
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, "checkout")
		if err != nil {
			return nil, err
		}
		ids := extractIDList(p.Args[constants.MOVIE_IDS])
		if len(ids) == 0 {
			msg := "movieIds argument is required for checkout mutation"
			return nil, errors.New(msg)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		messages, _, err := memberRepo.Checkout(ctx, username, ids)
		if err != nil {
			return nil, utils.LogError(fmt.Sprintf("failed to checkout for user %s:", username), err)
		}
		return messages, nil
	},
}

var SetAPIChoiceField = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		constants.USERNAME: usernameArg,
		constants.API_CHOICE: &graphql.ArgumentConfig{
			Type:         graphql.String,
			DefaultValue: constants.REST_API,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		username, err := getStringArg(p, constants.USERNAME, constants.SET_API_CHOICE)
		if err != nil {
			return nil, err
		}
		apiChoice, ok := p.Args[constants.API_CHOICE].(string)
		if !ok || apiChoice == "" {
			msg := "apiChoice argument is required for setAPIChoice mutation"
			return nil, errors.New(msg)
		}
		if apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API {
			msg := fmt.Sprintf("apiChoice must be either %s or %s", constants.REST_API, constants.GRAPHQL_API)
			return nil, errors.New(msg)
		}
		ctx, err := getContext(p)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = memberRepo.SetMemberAPIChoice(ctx, username, apiChoice)
		if err != nil {
			return nil, utils.LogError("", err)
		}
		return fmt.Sprintf("successfully set %s api choice to %s", username, apiChoice), nil
	},
}

// Helper to convert []interface{} to []string
func extractIDList(arg interface{}) []string {
	idsRaw, ok := arg.([]interface{})
	if !ok {
		return nil
	}
	ids := make([]string, len(idsRaw))
	for i, v := range idsRaw {
		ids[i], _ = v.(string)
	}
	return ids
}
