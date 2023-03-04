package main

import (
	"context"
	"log"
	"os"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func main() {
	ctx := context.Background()
	credentials := reddit.Credentials{
		ID:       os.Getenv("CLIENT_ID"),
		Secret:   os.Getenv("CLIENT_SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	log.Println("creating client")
	client, err := reddit.NewClient(credentials)

	if err != nil {
		panic(err)
	}

	log.Println("fetchint new post")
	kiszamolo, resp, err := client.Subreddit.NewPosts(ctx, "kiszamolo", &reddit.ListOptions{Limit: 1})
	if err != nil {
		log.Println(resp)
		panic(err)
	}

	for _, post := range kiszamolo {
		log.Println(post.Title)
		log.Println(post.Body)
	}

	log.Println("exiting")
}
