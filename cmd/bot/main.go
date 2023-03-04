package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var client *reddit.Client

func handle(ctx context.Context) (string, error) {
	log.Println("fetchint new post")
	kiszamolo, resp, err := client.Subreddit.NewPosts(ctx, "kiszamolo", &reddit.ListOptions{Limit: 1})
	if err != nil {
		log.Println(resp)
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

	if os.Getenv("ENV") == "DEV" {
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
		handle(ctx)
		return
	}

	client = reddit.DefaultClient()

	lambda.Start(handle)
}
