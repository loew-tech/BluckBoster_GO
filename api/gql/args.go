package gql

import "github.com/graphql-go/graphql"

var (
	usernameArg   = &graphql.ArgumentConfig{Type: graphql.ID}
	starArg       = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
	movieIDArg    = &graphql.ArgumentConfig{Type: graphql.ID, DefaultValue: ""}
	movieTitleArg = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
	movieIDsArg   = &graphql.ArgumentConfig{Type: graphql.NewList(graphql.ID), DefaultValue: []string{}}
	directorArg   = &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: ""}
)
