package github

import (
	"encoding/json"
	"net/http"
)

var (
	repoPath   = "repos/carlcamit/myob-pitt"
	commitPath = "/commits/main"

	// Github API encourages making requests
	// with this Accept header
	//
	// https://docs.github.com/en/rest/overview/media-types#request-specific-version
	HeaderAccept = "application/vnd.github.v3+json"
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

type Checker interface {
	CheckStatus() (int, error)
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

type Getter interface {
	GetDescription() (string, error)
	GetLatestSHA() (string, error)
}

type Repository struct {
	Description string `json:"description"`
}

func (c *Client) GetDescription() (string, error) {
	url := c.hostURL + repoPath
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", HeaderAccept)

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var repo Repository
	if err := json.NewDecoder(res.Body).Decode(&repo); err != nil {
		return "", err
	}

	return repo.Description, nil
}

type Commit struct {
	SHA string `json:"sha"`
}

func (c *Client) GetLatestSHA() (string, error) {
	url := c.hostURL + repoPath + commitPath
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", HeaderAccept)

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var commit Commit
	if err := json.NewDecoder(res.Body).Decode(&commit); err != nil {
		return "", err
	}

	return commit.SHA, nil
}
