package bot

import (
	"context"
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
	client, err := reddit.NewClient(credentials)

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

func (h *Handler) Handle(ctx context.Context) error {
	log.Println("fetching new post")
	topPosts, resp, err := h.client.Subreddit.NewPosts(ctx, os.Getenv("SUBREDDIT"), &reddit.ListOptions{Limit: 10})
	if err != nil {
		log.Printf("%+v", resp)
		log.Println(err)
		return err
	}

	for _, post := range topPosts {
		log.Println(post.Title)
		if !h.process(ctx, post) {
			log.Println("stop processing")
			break
		}
	}

	log.Println("done")
	return nil
}

func (h *Handler) process(ctx context.Context, post *reddit.Post) bool {
	if post.Created.Add(10 * time.Minute).Before(time.Now()) {
		log.Println("post is old, not commenting")
		return false
	}

	if !shouldComment(ctx, post) {
		log.Println("not allampapir related, not commenting")
		return true
	}

	_, _, err := h.client.Comment.Submit(ctx, post.FullID, "Szia!\nÚgy látom, az állampapírok összehasonlításában kérsz segítséget.\nHa még nem tetted meg, látogass el a https://allampapirkalkulator.hu/ oldalra, ahol ki tudod számolni a hozamokat és rengeteg más hasznos infót is találsz.\nÜdv")

	if err != nil {
		log.Println("error while adding comment")
		log.Println(err)
		return true
	}

	log.Println("comment added")

	return true
}

func shouldComment(ctx context.Context, post *reddit.Post) bool {
	if strings.Contains(post.Body, "allampapirkalkulator") ||
		strings.Contains(post.Body, "állampapírkalkulátor") ||
		strings.Contains(post.Body, "állampapír kalkulátor") ||
		strings.Contains(post.Body, "allampapir kalkulator") {
		return false
	}

	count := 0
	counts := map[string]int{}

	if strings.Contains(post.Title, "PMAP") ||
		strings.Contains(post.Title, "PMÁP") {
		count++
		counts["PMAP"]++
	}

	if strings.Contains(post.Title, "BMAP") ||
		strings.Contains(post.Title, "BMÁP") {
		count++
		counts["BMAP"]++
	}

	if strings.Contains(post.Title, "EMAP") ||
		strings.Contains(post.Title, "EMÁP") {
		count++
		counts["EMAP"]++
	}

	if strings.Contains(post.Title, "DKJ") {
		count++
		counts["DKJ"]++
	}

	if count >= 2 {
		return true
	}

	if strings.Contains(post.Body, "PMAP") ||
		strings.Contains(post.Body, "PMÁP") {
		count++
		counts["PMAP"]++
	}

	if strings.Contains(post.Body, "BMAP") ||
		strings.Contains(post.Body, "BMÁP") {
		count++
		counts["BMAP"]++
	}

	if strings.Contains(post.Body, "EMAP") ||
		strings.Contains(post.Body, "EMÁP") {
		count++
		counts["EMAP"]++
	}

	if strings.Contains(post.Title, "DKJ") {
		count++
		counts["DKJ"]++
	}

	if strings.Contains(post.Body, " hozam") ||
		strings.Contains(post.Title, "hozam") {
		count++
		counts["hozam"]++
	}

	if len(counts) >= 2 {
		return true
	}

	return false
}
