package bot

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var DefaultHandler *Handler

func Init() {
	credentials := reddit.Credentials{
		ID:       os.Getenv("CLIENT_ID"),
		Secret:   os.Getenv("CLIENT_SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	log.Println("creating client")
	client, err := reddit.NewClient(credentials, reddit.WithBaseURL(os.Getenv("REDDIT_URL")))

	if err != nil {
		panic(err)
	}

	DefaultHandler = &Handler{
		client: client,
	}
}

func InitClient(c *http.Client) {
	credentials := reddit.Credentials{
		ID:       os.Getenv("CLIENT_ID"),
		Secret:   os.Getenv("CLIENT_SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	log.Println("creating client")
	client, err := reddit.NewClient(
		credentials,
		reddit.WithHTTPClient(c),
		reddit.WithTokenURL(os.Getenv("REDDIT_URL")),
		reddit.WithBaseURL(os.Getenv("REDDIT_URL")),
	)

	if err != nil {
		panic(err)
	}

	DefaultHandler = &Handler{
		client: client,
	}
}

type Handler struct {
	client *reddit.Client
}

func (h *Handler) Handle(ctx context.Context) (string, error) {
	log.Println("fetching new post")
	topPosts, resp, err := h.client.Subreddit.NewPosts(ctx, os.Getenv("SUBREDDIT"), &reddit.ListOptions{Limit: 1})
	if err != nil {
		log.Println(resp)
		return "", err
	}

	var res string
	for _, post := range topPosts {
		log.Println(post.Title)
		log.Println(post.Body)
		res = post.Title
	}

	log.Println("done")
	return res, nil
}
