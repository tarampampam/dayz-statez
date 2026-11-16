package api_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"gh.tarampamp.am/dayz-statez/server/internal/testutil/assert"
	"gh.tarampamp.am/dayz-statez/server/pkg/api"
)

// mockHTTPClient returns the preconfigured response and error for every request, ignoring the request itself.
type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m mockHTTPClient) Do(*http.Request) (*http.Response, error) { return m.resp, m.err } //nolint:nilnil

// newResponse builds an [http.Response] with the given status code and body.
func newResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Status:     http.StatusText(code),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("endpoints", func(t *testing.T) {
		t.Parallel()

		for name, tt := range map[string]struct {
			call func(context.Context, *api.Client) (any, error)
			body string
		}{
			"animals":     {func(c context.Context, a *api.Client) (any, error) { return a.Animals(c) }, `{}`},
			"areas":       {func(c context.Context, a *api.Client) (any, error) { return a.Areas(c) }, `[]`},
			"events":      {func(c context.Context, a *api.Client) (any, error) { return a.Events(c) }, `{}`},
			"hits":        {func(c context.Context, a *api.Client) (any, error) { return a.Hits(c) }, `{}`},
			"players":     {func(c context.Context, a *api.Client) (any, error) { return a.Players(c) }, `[]`},
			"territories": {func(c context.Context, a *api.Client) (any, error) { return a.Territories(c) }, `[]`},
			"vehicles":    {func(c context.Context, a *api.Client) (any, error) { return a.Vehicles(c) }, `[]`},
			"weapons":     {func(c context.Context, a *api.Client) (any, error) { return a.Weapons(c) }, `{}`},
			"world":       {func(c context.Context, a *api.Client) (any, error) { return a.World(c) }, `{}`},
			"zombies":     {func(c context.Context, a *api.Client) (any, error) { return a.Zombies(c) }, `{}`},
			"status":      {func(c context.Context, a *api.Client) (any, error) { return a.Status(c) }, `{}`},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				client := api.NewClient(
					"",
					"",
					api.WithHTTPClient(mockHTTPClient{resp: newResponse(http.StatusOK, tt.body)}), //nolint:bodyclose
				)

				result, err := tt.call(context.Background(), client)

				assert.NoError(t, err)
				assert.NotNil(t, result)
			})
		}
	})

	t.Run("status codes", func(t *testing.T) {
		t.Parallel()

		for name, tt := range map[string]struct {
			giveCode int
			giveBody string
			wantErr  error // expected sentinel error (errors.Is); nil means no error expected
		}{
			"ok":           {http.StatusOK, `{"by_type":{"wolf":1}}`, nil},
			"no content":   {http.StatusNoContent, "", api.ErrNoContent},
			"unauthorized": {http.StatusUnauthorized, "", api.ErrBadAuth},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				client := api.NewClient(
					"",
					"",
					api.WithHTTPClient(mockHTTPClient{resp: newResponse(tt.giveCode, tt.giveBody)}), //nolint:bodyclose
				)

				result, err := client.Animals(context.Background())

				if tt.wantErr != nil {
					assert.ErrorIs(t, err, tt.wantErr)
					assert.Nil(t, result)

					return
				}

				assert.NoError(t, err)
				assert.NotNil(t, result)
			})
		}
	})

	t.Run("unexpected status", func(t *testing.T) {
		t.Parallel()

		client := api.NewClient(
			"",
			"",
			api.WithHTTPClient(mockHTTPClient{resp: newResponse(http.StatusInternalServerError, "")}), //nolint:bodyclose
		)

		result, err := client.Animals(context.Background())

		assert.ErrorContains(t, err, "unexpected status code")
		assert.Nil(t, result)
	})

	t.Run("decode error", func(t *testing.T) {
		t.Parallel()

		client := api.NewClient(
			"",
			"",
			api.WithHTTPClient(mockHTTPClient{resp: newResponse(http.StatusOK, `not a json`)}), //nolint:bodyclose
		)

		result, err := client.Animals(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("request error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("boom")

		client := api.NewClient("", "", api.WithHTTPClient(mockHTTPClient{err: wantErr}))

		result, err := client.Animals(context.Background())

		assert.ErrorIs(t, err, wantErr)
		assert.Nil(t, result)
	})
}
