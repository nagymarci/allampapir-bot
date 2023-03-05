package bot_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/nagymarci/allampapir-bot/internal/bot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()

	resp, err := ioutil.ReadFile("testdata/test.json")
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Log("request received")
		t.Log(r)
		if r.Method == http.MethodPost && r.URL.Path == "/api/comment" {
			t.Log("comment received")
			assert.True(t, true)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))

	defer srv.Close()

	os.Setenv("REDDIT_URL", srv.URL)
	os.Setenv("SUBREDDIT", "test")

	t.Log(srv.URL)

	bot.InitClient(srv.Client())
	require.NoError(t, bot.DefaultHandler.Handle(ctx))
}

const (
	inputDir  = "testdata/"
	subreddit = "test"
)

var createdAtRegexp = regexp.MustCompile(`\"created_utc\": [0-9]+\.[0-9]+,`)

func TestFlows(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		comment   bool
		created   int64
		commented bool
	}{
		"one_result_no_comment": {
			created: time.Now().Unix(),
		},
		"one_result_comment": {
			comment: true,
			created: time.Now().Unix(),
		},
		"one_result_too_old_to_comment": {
			comment: false,
			created: time.Now().Add(-11 * time.Minute).Unix(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp := readInput(t, name)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				t.Log("request received")
				t.Log(r)

				if r.Method == http.MethodPost && r.URL.Path == "/api/comment" {
					assert.True(t, test.comment, "should not comment")
					test.commented = true
					w.WriteHeader(http.StatusOK)
					return
				}

				if r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/r/%s/new", subreddit) {
					res := createdAtRegexp.ReplaceAll(
						resp,
						[]byte(fmt.Sprintf(`"created_utc": %d,`, test.created)),
					)

					w.WriteHeader(http.StatusOK)
					w.Write(res)
					return
				}

				t.Log(r.URL.Path)
				assert.Fail(t, "method and path not valid")
			}))

			defer srv.Close()

			os.Setenv("REDDIT_URL", srv.URL)
			os.Setenv("SUBREDDIT", subreddit)

			t.Log(srv.URL)

			bot.InitClient(srv.Client())
			require.NoError(t, bot.DefaultHandler.Handle(ctx))
			assert.Equal(t, test.comment, test.commented)
		})
	}
}

func readInput(t *testing.T, file string) []byte {
	response, err := os.ReadFile(filepath.Join(inputDir, fmt.Sprintf("%s.json", file)))
	require.NoError(t, err, "reading input")
	return response
}
