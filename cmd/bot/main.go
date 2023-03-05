package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/nagymarci/allampapir-bot/internal/bot"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var client *reddit.Client

func main() {
	ctx := context.Background()

	if os.Getenv("ENV") == "AWS" {
		loadEnvironmentVariables()
	}

	bot.Init()

	if os.Getenv("ENV") == "DEV" {
		bot.DefaultHandler.Handle(ctx)
		return
	}

	lambda.Start(bot.DefaultHandler.Handle)
}

func loadEnvironmentVariables() {
	secretCache, err := secretcache.New()

	if err != nil {
		panic(err)
	}

	result, err := secretCache.GetSecretString(os.Getenv("SECRET_NAME"))

	secrets := map[string]string{}

	if json.Unmarshal([]byte(result), &secrets) != nil {
		panic(err)
	}

	log.Println("setting env variables")
	for key, value := range secrets {
		log.Println(key)
		os.Setenv(key, value)
	}
	log.Println("env variables set")
}
