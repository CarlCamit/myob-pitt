package github

import (
	"net/http"
)

// Client specifies the settings for
// communicating with the API
type Client struct {
	hostURL string
	// Re-use the same client so TCP connections can be cached
	// http.Client is safe for re-use across goroutines
	client *http.Client
}

// NewClient returns a Client for the
// provided hostURL
func NewClient(hostURL string) *Client {
	return &Client{
		hostURL: hostURL,
		client:  &http.Client{Transport: http.DefaultTransport},
	}
}
