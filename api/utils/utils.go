package utils

import (
	"blockbuster/api/data"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	once         sync.Once
)

func GetDynamoClient() *dynamodb.Client {
	once.Do(func() {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalln("FAILED TO INSTANTIATE MemberRepo", err)
		}
		dynamoClient = dynamodb.NewFromConfig(cfg)
	})
	return dynamoClient
}

// LogError wraps an error with a message, logs it, and returns it.
// If the original error is nil, creates a new one from the message.
func LogError(msg string, err error) error {
	if err == nil {
		err = errors.New(msg)
	} else {
		err = fmt.Errorf("%s: %w", msg, err)
	}
	log.Println(err)
	return err
}

// Contains returns true if the item is found in the list.
func Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// GetStringArg safely extracts a required string arg from the resolver params.
func GetStringArg(params gin.Params, argName string) (string, error) {
	val, ok := params.Get(argName)
	if !ok || val == "" {
		msg := fmt.Sprintf("%s argument is required", argName)
		log.Println(msg)
		return "", errors.New(msg)
	}
	return val, nil
}

func GetSliceFromMapKeys(m map[string]bool) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func AccumulateMovieMetricsWithWeight(a, b data.MovieMetrics, weight int) data.MovieMetrics {
	return data.MovieMetrics{
		Acting:         a.Acting + b.Acting*float64(weight),
		Action:         a.Action + b.Action*float64(weight),
		Cinematography: a.Cinematography + b.Cinematography*float64(weight),
		Comedy:         a.Comedy + b.Comedy*float64(weight),
		Directing:      a.Directing + b.Directing*float64(weight),
		Drama:          a.Drama + b.Drama*float64(weight),
		Fantasy:        a.Fantasy + b.Fantasy*float64(weight),
		Horror:         a.Horror + b.Horror*float64(weight),
		Romance:        a.Romance + b.Romance*float64(weight),
		StoryTelling:   a.StoryTelling + b.StoryTelling*float64(weight),
		Suspense:       a.Suspense + b.Suspense*float64(weight),
		Writing:        a.Writing + b.Writing*float64(weight),
	}
}

func AverageMetrics(m data.MovieMetrics, count int) data.MovieMetrics {
	if count <= 1 {
		return m
	}
	return data.MovieMetrics{
		Acting:         m.Acting / float64(count),
		Action:         m.Action / float64(count),
		Cinematography: m.Cinematography / float64(count),
		Comedy:         m.Comedy / float64(count),
		Directing:      m.Directing / float64(count),
		Drama:          m.Drama / float64(count),
		Fantasy:        m.Fantasy / float64(count),
		Horror:         m.Horror / float64(count),
		Romance:        m.Romance / float64(count),
		StoryTelling:   m.StoryTelling / float64(count),
		Suspense:       m.Suspense / float64(count),
		Writing:        m.Writing / float64(count),
	}
}

// Euclidean distance between two MovieMetrics
func MetricDistance(a, b data.MovieMetrics) float64 {
	sum := 0.0
	sum += (a.Acting - b.Acting) * (a.Acting - b.Acting)
	sum += (a.Action - b.Action) * (a.Action - b.Action)
	sum += (a.Cinematography - b.Cinematography) * (a.Cinematography - b.Cinematography)
	sum += (a.Comedy - b.Comedy) * (a.Comedy - b.Comedy)
	sum += (a.Directing - b.Directing) * (a.Directing - b.Directing)
	sum += (a.Drama - b.Drama) * (a.Drama - b.Drama)
	sum += (a.Fantasy - b.Fantasy) * (a.Fantasy - b.Fantasy)
	sum += (a.Horror - b.Horror) * (a.Horror - b.Horror)
	sum += (a.Romance - b.Romance) * (a.Romance - b.Romance)
	sum += (a.StoryTelling - b.StoryTelling) * (a.StoryTelling - b.StoryTelling)
	sum += (a.Suspense - b.Suspense) * (a.Suspense - b.Suspense)
	sum += (a.Writing - b.Writing) * (a.Writing - b.Writing)
	return sum
}
