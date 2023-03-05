package bot

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	log.Println("creating client")
	client, err := reddit.NewReadonlyClient(
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
	topPosts, resp, err := h.client.Subreddit.NewPosts(ctx, os.Getenv("SUBREDDIT"), &reddit.ListOptions{Limit: 10})
	if err != nil {
		log.Printf("%+v", resp)
		bytes, err1 := io.ReadAll(resp.Body)
		if err1 != nil {
			return "", err
		}
		defer resp.Body.Close()

		log.Println(string(bytes))
		log.Println(err)
		return "", err
	}

	var res string
	for _, post := range topPosts {
		log.Println(post.Title)
		log.Println(post.Body)
		h.process(ctx, post)
	}

	log.Println("done")
	return res, nil
}

func (h *Handler) process(ctx context.Context, post *reddit.Post) {
	if post.Created.Add(10 * time.Minute).Before(time.Now()) {
		log.Println("post is old, not commenting")
		return
	}

	if !shouldComment(ctx, post) {
		log.Println("not allampapir related, not commenting")
		return
	}

	_, _, err := h.client.Comment.Submit(ctx, post.FullID, "this is a response")

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("comment added")
}

func shouldComment(ctx context.Context, post *reddit.Post) bool {
	return strings.Contains(post.Title, "PMAP")
}
