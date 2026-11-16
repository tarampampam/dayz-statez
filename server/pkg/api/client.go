package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrNoContent is returned on a "204: No Content" response, meaning the server has no data to serve. The data is
	// produced by the game mod, so an empty response usually points to a misconfigured server (the wrong data directory
	// is being watched), a game mod that failed to load, or a game server that is not running.
	ErrNoContent = errors.New("no content")

	// ErrBadAuth is returned on a "401: Unauthorized" response, meaning the provided auth token is missing or invalid.
	ErrBadAuth = errors.New("bad auth token")
)

// httpClient is the subset of [http.Client] used by [Client]. It is defined as an interface so the transport can be
// swapped out in tests (see [WithHTTPClient]).
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientOption customizes a [Client] during construction in [NewClient].
type ClientOption func(*Client)

// WithHTTPClient overrides the underlying HTTP client used to perform requests.
func WithHTTPClient(c httpClient) ClientOption { //nolint:iface
	return func(client *Client) { client.c = c }
}

// WithInsecureTLS disables TLS certificate verification on the default HTTP client. It has no effect when a custom
// client is supplied via [WithHTTPClient] or when the transport is not an [http.Transport].
//
// Use this only against servers with self-signed certificates and never in production.
func WithInsecureTLS() ClientOption {
	return func(client *Client) {
		if hc, ok1 := client.c.(*http.Client); ok1 {
			if tr, ok2 := hc.Transport.(*http.Transport); ok2 {
				if tr.TLSClientConfig == nil {
					tr.TLSClientConfig = &tls.Config{}
				}

				tr.TLSClientConfig.InsecureSkipVerify = true
			}
		}
	}
}

// Client is an HTTP client for the DayZ StateZ server API. Create it with [NewClient].
type Client struct {
	c         httpClient
	baseUrl   string
	authToken string
}

// NewClient creates a [Client] targeting the given base URL and authenticating with the given token. A trailing slash
// in the base URL is trimmed. The behavior can be further tuned with [ClientOption] values.
func NewClient(baseUrl, authToken string, opts ...ClientOption) *Client {
	client := &Client{
		c: &http.Client{
			Timeout: 30 * time.Second, //nolint:mnd
			Transport: &http.Transport{
				ForceAttemptHTTP2:     true,
				ResponseHeaderTimeout: 10 * time.Second, //nolint:mnd
				IdleConnTimeout:       90 * time.Second, //nolint:mnd
			},
		},
		baseUrl:   strings.TrimSpace(strings.TrimRight(baseUrl, "/")),
		authToken: authToken,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Animals fetches the currently tracked animals from the "/animals" endpoint.
func (c *Client) Animals(ctx context.Context) (*Animals, error) {
	resp, respErr := c.do(ctx, "/animals", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Animals](resp.Body)
}

// Areas fetches the defined map areas from the "/areas" endpoint.
func (c *Client) Areas(ctx context.Context) (*Areas, error) {
	resp, respErr := c.do(ctx, "/areas", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Areas](resp.Body)
}

// Events fetches the recorded game events from the "/events" endpoint.
func (c *Client) Events(ctx context.Context) (*Events, error) {
	resp, respErr := c.do(ctx, "/events", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Events](resp.Body)
}

// Hits fetches the recorded hits from the "/hits" endpoint.
func (c *Client) Hits(ctx context.Context) (*Hits, error) {
	resp, respErr := c.do(ctx, "/hits", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Hits](resp.Body)
}

// Players fetches the currently tracked players from the "/players" endpoint.
func (c *Client) Players(ctx context.Context) (*Players, error) {
	resp, respErr := c.do(ctx, "/players", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Players](resp.Body)
}

// Territories fetches the defined territories from the "/territories" endpoint.
func (c *Client) Territories(ctx context.Context) (*Territories, error) {
	resp, respErr := c.do(ctx, "/territories", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Territories](resp.Body)
}

// Vehicles fetches the currently tracked vehicles from the "/vehicles" endpoint.
func (c *Client) Vehicles(ctx context.Context) (*Vehicles, error) {
	resp, respErr := c.do(ctx, "/vehicles", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Vehicles](resp.Body)
}

// Weapons fetches the weapon usage statistics from the "/weapons" endpoint.
func (c *Client) Weapons(ctx context.Context) (*Weapons, error) {
	resp, respErr := c.do(ctx, "/weapons", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Weapons](resp.Body)
}

// World fetches the current world state from the "/world" endpoint.
func (c *Client) World(ctx context.Context) (*World, error) {
	resp, respErr := c.do(ctx, "/world", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[World](resp.Body)
}

// Zombies fetches the currently tracked zombies from the "/zombies" endpoint.
func (c *Client) Zombies(ctx context.Context) (*Zombies, error) {
	resp, respErr := c.do(ctx, "/zombies", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Zombies](resp.Body)
}

// Status fetches the public server/world status summary from the "/status" endpoint.
func (c *Client) Status(ctx context.Context) (*Status, error) {
	resp, respErr := c.do(ctx, "/status", http.NoBody)
	if respErr != nil {
		return nil, respErr
	}
	defer func() { _ = resp.Body.Close() }()

	if err := clientCheckResponse(resp); err != nil {
		return nil, err
	}

	return clientDecodeBody[Status](resp.Body)
}

// do builds and performs a GET request against the configured base URL, attaching the bearer auth token when set.
func (c *Client) do(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+path, body)
	if reqErr != nil {
		return nil, reqErr
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, respErr := c.c.Do(req)
	if respErr != nil {
		return nil, respErr
	}

	return resp, nil
}

// clientCheckResponse maps an HTTP response status to an error: [ErrNoContent] for 204, [ErrBadAuth] for 401, nil for
// 200, and a generic error for anything else.
func clientCheckResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusNoContent:
		return ErrNoContent
	case http.StatusUnauthorized:
		return ErrBadAuth
	case http.StatusOK:
		break // continue to decode the response body
	default:
		return errors.New("unexpected status code: " + resp.Status)
	}

	return nil
}

// clientDecodeBody decodes a JSON-encoded reader into a newly allocated value of type T.
func clientDecodeBody[T any](b io.Reader) (*T, error) {
	var data T
	if decodeErr := json.NewDecoder(b).Decode(&data); decodeErr != nil {
		return nil, decodeErr
	}

	return &data, nil
}
