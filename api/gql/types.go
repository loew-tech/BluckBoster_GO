package gql

import (
	"github.com/graphql-go/graphql"

	"blockbuster/api/constants"
)

var MovieType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MOVIE_TYPE,
	Fields: graphql.Fields{
		constants.ID:        &graphql.Field{Type: graphql.String},
		constants.INVENTORY: &graphql.Field{Type: graphql.Int},
		constants.RATING:    &graphql.Field{Type: graphql.String},
		constants.REVIEW:    &graphql.Field{Type: graphql.String},
		constants.RENTED:    &graphql.Field{Type: graphql.Int},
		constants.SYNOPSIS:  &graphql.Field{Type: graphql.String},
		constants.TRIVIA:    &graphql.Field{Type: graphql.String},
		constants.YEAR:      &graphql.Field{Type: graphql.String},
		constants.CAST:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.DIRECTOR:  &graphql.Field{Type: graphql.String},
		constants.TITLE:     &graphql.Field{Type: graphql.String},
	},
})

var MemberType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.MEMBER_TYPE,
	Fields: graphql.Fields{
		constants.USERNAME:    &graphql.Field{Type: graphql.String},
		constants.FIRSTNAME:   &graphql.Field{Type: graphql.String},
		constants.LASTNAME:    &graphql.Field{Type: graphql.String},
		constants.CHECKED_OUT: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.CART_STRING: &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.RENTED:      &graphql.Field{Type: &graphql.List{OfType: graphql.String}},
		constants.TYPE:        &graphql.Field{Type: graphql.String},
	},
})

var KevingBaconType = graphql.NewObject(graphql.ObjectConfig{
	Name: constants.KEVING_BACON_TYPE,
	Fields: graphql.Fields{
		constants.STAR:            &graphql.Field{Type: graphql.String},
		constants.STARS:           &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_STARS:     &graphql.Field{Type: graphql.Int},
		constants.MOVIES:          &graphql.Field{Type: graphql.NewList(MovieType)},
		constants.TOTAL_MOVIES:    &graphql.Field{Type: graphql.Int},
		constants.DIRECTORS:       &graphql.Field{Type: graphql.NewList(graphql.String)},
		constants.TOTAL_DIRECTORS: &graphql.Field{Type: graphql.Int},
	},
})
