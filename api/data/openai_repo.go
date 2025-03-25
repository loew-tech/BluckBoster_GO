package data

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	// "github.com/openai/openai-go"
	// "github.com/openai/openai-go/option"
)

func GetMovieTrivia(movie string) (string, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		// @TODO: handle err without failing app
		log.Fatalln("Err opening .env file")
	}

	// @TODO: remove hard coded values
	key, ok := envFile["OPENAI"]
	if !ok {
		log.Fatalln("Failed to retrieve api key from .env")
	}

	fmt.Println("**\nMovie=", movie, "\nkey=", key, "\n**")

	return "", nil
}
