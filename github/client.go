package github

import (
	"net/http"
)

// Github API encourages making requests
// with this Accept header
//
// https://docs.github.com/en/rest/overview/media-types#request-specific-version
var HeaderAccept = "application/vnd.github.v3+json"

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

// CheckStatus checks the availability of the
// Github API so that the health endpoint
// can respond appropriately
func (c *Client) CheckStatus() (int, error) {
	req, err := http.NewRequest(http.MethodGet, c.hostURL, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", HeaderAccept)

	res, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}

	return res.StatusCode, nil
}
