package bot_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nagymarci/allampapir-bot/internal/bot"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()

	resp, err := ioutil.ReadFile("testdata/test.json")
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Log("request received")
		t.Log(r)
		// if r.Method == http.MethodPost {
		// 	w.WriteHeader(http.StatusOK)
		// 	fmt.Fprint(w, `{“access_token”:“60acf87776bda95357c7564e21e0b69b”,“refresh_token”:“authorizationCode”,“token_type”:“SSO”,“expires_in”:60000}`)
		// 	return
		// }

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
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
