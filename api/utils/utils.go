package utils

import (
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
