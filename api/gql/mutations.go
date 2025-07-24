// File: mutations.go
package gql

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		ids := extractIDList(p.Args[constants.MOVIE_IDS])
		if len(ids) == 0 {
			return nil, getFormattedError("non-empty 'movieIds' argument is required for returnRentals mutation", http.StatusBadRequest)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, err
		}
		messages, _, err := memberService.Return(ctx, username, ids)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		movieID, err := getStringArg(p, constants.MOVIE_ID, "updateCart")
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		shouldRemoveFromCart, _ := p.Args[constants.REMOVE_FROM_CART].(bool)
		ctx, err := getContext(p)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		var inserted bool
		if shouldRemoveFromCart {
			inserted, err = memberService.RemoveFromCart(ctx, username, movieID)
		} else {
			inserted, err = memberService.AddToCart(ctx, username, movieID)
		}
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
		} else if !inserted {
			return fmt.Sprintf("did not modify cart for %s", username), nil
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
			return nil, getFormattedError("non-empty 'movieIds' argument is required for checkout mutation", http.StatusBadRequest)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		messages, _, err := memberService.Checkout(ctx, username, ids)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
			return nil, getFormattedError("'apiChoice' argument is required for setAPIChoice mutation", http.StatusBadRequest)
		}
		if apiChoice != constants.REST_API && apiChoice != constants.GRAPHQL_API {
			msg := fmt.Sprintf("apiChoice must be either %s or %s", constants.REST_API, constants.GRAPHQL_API)
			return nil, getFormattedError(msg, http.StatusBadRequest)
		}
		ctx, err := getContext(p)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusBadRequest)
		}
		err = memberService.SetAPIChoice(ctx, username, apiChoice)
		if err != nil {
			return nil, getFormattedError(err.Error(), http.StatusInternalServerError)
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
