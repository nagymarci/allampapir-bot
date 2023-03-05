package bot_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nagymarci/allampapir-bot/internal/bot"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header)
		t.Log("request received")
		t.Log(r)
		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	os.Setenv("REDDIT_URL", srv.URL)
	os.Setenv("SUBREDDIT", "test")

	t.Log(srv.URL)

	bot.InitClient(srv.Client())
	res, err := bot.DefaultHandler.Handle(ctx)
	require.NoError(t, err)

	t.Log(res)
}
