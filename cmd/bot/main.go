package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var client *reddit.Client

//var secrets []string = []string{"CLIENT_ID", "CLIENT_SECRET", "USERNAME", "PASSWORD"}

func handle(ctx context.Context) (string, error) {
	log.Println("fetchint new post")
	kiszamolo, resp, err := client.Subreddit.NewPosts(ctx, "kiszamolo", &reddit.ListOptions{Limit: 1})
	if err != nil {
		log.Println(resp)
		log.Println(os.Getenv("CLIENT_ID"))
		return "", err
	}

	var res string
	for _, post := range kiszamolo {
		log.Println(post.Title)
		log.Println(post.Body)
		res = post.Title
	}

	log.Println("done")
	return res, nil
}

func main() {
	ctx := context.Background()

	if os.Getenv("ENV") == "AWS" {
		loadEnvironmentVariables()
	}

	credentials := reddit.Credentials{
		ID:       os.Getenv("CLIENT_ID"),
		Secret:   os.Getenv("CLIENT_SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	log.Println("creating client")
	c, err := reddit.NewClient(credentials)

	if err != nil {
		panic(err)
	}

	client = c

	if os.Getenv("ENV") == "DEV" {
		handle(ctx)
		return
	}

	lambda.Start(handle)
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

	for key, value := range secrets {
		os.Setenv(key, value)
	}

}
