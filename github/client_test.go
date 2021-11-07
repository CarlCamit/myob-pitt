package github

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCheckStatus(t *testing.T) {
	testCases := []struct {
		name     string
		status   int
		expected int
	}{
		{
			name:     "returns a 200 if the server responds with a 200",
			status:   http.StatusOK,
			expected: http.StatusOK,
		},
		{
			name:     "returns a 503 if the server responds with a 503",
			status:   http.StatusServiceUnavailable,
			expected: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer s.Close()

			c := NewClient(s.URL)

			status, err := c.CheckStatus()
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			if status != tc.expected {
				t.Errorf("expected status to be %v, got %v", tc.expected, status)
			}
		})
	}
}

func TestGetDescription(t *testing.T) {
	testCases := []struct {
		name         string
		responseFile string
		status       int
		expectedErr  error
		expected     string
	}{
		{
			name:         "returns an error if a 200 is not received",
			responseFile: "./testdata/empty.json",
			status:       http.StatusNotFound,
			expectedErr:  fmt.Errorf("cannot get description received a"),
			expected:     "",
		},
		{
			name:         "returns the description from the server response",
			responseFile: "./testdata/description.json",
			status:       http.StatusOK,
			expectedErr:  nil,
			expected:     "the description of the repository",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.responseFile)
			if err != nil {
				t.Fatalf("failed to open test data file, %v", err)
			}

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				io.Copy(w, f)
			}))
			defer s.Close()

			c := NewClient(s.URL)

			description, err := c.GetDescription()
			if err != nil {
				if tc.expectedErr == nil {
					t.Fatalf("expected err to be nil got %v", err)
				}

				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected err to contain %s, got %s", tc.expectedErr, err)
				}
			}

			if description != tc.expected {
				t.Errorf("expected description to be %s, got %s", tc.expected, description)
			}
		})
	}
}

func TestGetLatestSHA(t *testing.T) {
	testCases := []struct {
		name         string
		responseFile string
		status       int
		expectedErr  error
		expected     string
	}{
		{
			name:         "returns an error if a 200 is not received",
			responseFile: "./testdata/empty.json",
			status:       http.StatusNotFound,
			expectedErr:  fmt.Errorf("cannot get latest sha received a"),
			expected:     "",
		},
		{
			name:         "returns the sha from the server response",
			responseFile: "./testdata/sha.json",
			status:       http.StatusOK,
			expectedErr:  nil,
			expected:     "b10149508a0e0f67d23e38c6293cd7a8cdbe42e5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.responseFile)
			if err != nil {
				t.Fatalf("failed to open test data file, %v", err)
			}

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				io.Copy(w, f)
			}))
			defer s.Close()

			c := NewClient(s.URL)

			sha, err := c.GetLatestSHA()
			if err != nil {
				if tc.expectedErr == nil {
					t.Fatalf("expected err to be nil got %v", err)
				}

				if !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Errorf("expected err to contain %s, got %s", tc.expectedErr, err)
				}
			}

			if sha != tc.expected {
				t.Errorf("expected sha to be %v, got %v", tc.expected, sha)
			}
		})
	}
}
