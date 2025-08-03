package gql

import (
	"context"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"

	graphsearch "blockbuster/api/graph_search"
)

var (
	movieGraph         graphsearch.MovieGraphInterface
	initMovieGraphOnce sync.Once
	initMovieGraphErr  error
)

func GetGQLHandler() func(*gin.Context) {
	initMovieGraphOnce.Do(func() {
		movieGraph, initMovieGraphErr = graphsearch.GetMovieGraph()
		if initMovieGraphErr != nil {
			log.Printf("Failed to initialize MovieGraph: %v", initMovieGraphErr)
		}
	})

	schema := getSchema()
	gqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		corsHandler.Handler(gqlHandler).ServeHTTP(c.Writer, c.Request)
	}
}

type contextKeyGin struct{}

var ginContextKey = contextKeyGin{}
