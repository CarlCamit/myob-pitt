package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
		expectedErr    error
		expected       []byte
	}{
		{
			name:           "returns hello world",
			expectedStatus: http.StatusOK,
			expectedErr:    fmt.Errorf("cannot get description received a"),
			expected:       []byte("hello world\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(Root())
			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status code to be %v, got %v", tc.expectedStatus, res.StatusCode)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response body, %v", err)
			}

			if diff := bytes.Compare(b, tc.expected); diff != 0 {
				t.Errorf("expected response body to be %s, got %s", tc.expected, b)
			}
		})
	}
}

type TestChecker struct {
	status int
	err    error
}

func (c *TestChecker) CheckStatus() (int, error) {
	return c.status, c.err
}

func TestHealth(t *testing.T) {
	testCases := []struct {
		name        string
		status      int
		err         error
		expectErr   bool
		expectedErr map[string]string
	}{
		{
			name:        "returns any error encountered",
			status:      http.StatusInternalServerError,
			err:         fmt.Errorf("failed to perform request"),
			expectErr:   true,
			expectedErr: map[string]string{"error": "failed to perform request"},
		},
		{
			name:        "returns the status of the server",
			status:      http.StatusOK,
			err:         nil,
			expectErr:   false,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			c := &TestChecker{
				status: tc.status,
				err:    tc.err,
			}

			handler := http.HandlerFunc(Health(c))
			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.status {
				t.Errorf("expected status code to be %v, got %v", tc.status, res.StatusCode)
			}

			if tc.expectErr {
				var response map[string]string
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response body, %v", err)
				}

				if response["error"] != tc.expectedErr["error"] {
					t.Errorf("expected error response to be %v, got %v", tc.expectedErr["error"], response["error"])
				}
			}
		})
	}
}

type TestGetter struct {
	description    string
	descriptionErr error
	sha            string
	shaErr         error
}

func (c *TestGetter) GetDescription() (string, error) {
	return c.description, c.descriptionErr
}

func (c *TestGetter) GetLatestSHA() (string, error) {
	return c.sha, c.shaErr
}

func TestMetadata(t *testing.T) {
	testCases := []struct {
		name           string
		status         int
		description    string
		descriptionErr error
		sha            string
		shaErr         error
		expectErr      bool
		expectedErr    map[string]string
		expected       map[string][]MyApplication
	}{
		{
			name:           "returns errors getting repository description",
			status:         http.StatusInternalServerError,
			description:    "",
			descriptionErr: fmt.Errorf("failed to get description"),
			sha:            "",
			shaErr:         nil,
			expectErr:      true,
			expectedErr:    map[string]string{"error": "failed to get description"},
			expected:       nil,
		},
		{
			name:           "returns errors getting latest sha",
			status:         http.StatusInternalServerError,
			description:    "",
			descriptionErr: nil,
			sha:            "",
			shaErr:         fmt.Errorf("failed to get latest sha"),
			expectErr:      true,
			expectedErr:    map[string]string{"error": "failed to get latest sha"},
			expected:       nil,
		},
		{
			name:           "returns the metadata",
			status:         http.StatusOK,
			description:    "description of the repo",
			descriptionErr: nil,
			sha:            "b9c2d825df36cbf44f",
			shaErr:         nil,
			expectErr:      false,
			expectedErr:    nil,
			expected: map[string][]MyApplication{"myapplication": {{
				Version:     "1.0.0",
				Description: "description of the repo",
				SHA:         "b9c2d825df36cbf44f",
			}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/metadata", nil)
			w := httptest.NewRecorder()

			c := &TestGetter{
				description:    tc.description,
				descriptionErr: tc.descriptionErr,
				sha:            tc.sha,
				shaErr:         tc.shaErr,
			}

			handler := http.HandlerFunc(Metadata(c))
			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.status {
				t.Errorf("expected status code to be %v, got %v", tc.status, res.StatusCode)
			}

			if tc.expectErr {
				var response map[string]string
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response body, %v", err)
				}

				if response["error"] != tc.expectedErr["error"] {
					t.Errorf("expected error response to be %v, got %v", tc.expectedErr["error"], response["error"])
				}
			} else {
				var response map[string][]MyApplication
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response body, %v", err)
				}

				r := response["myapplication"][0]
				e := tc.expected["myapplication"][0]

				if r.Version != e.Version {
					t.Errorf("expected response version to be %v, got %v", e.Version, r.Version)
				}

				if r.Description != e.Description {
					t.Errorf("expected response description to be %v, got %v", e.Description, r.Description)
				}

				if r.SHA != e.SHA {
					t.Errorf("expected error response to be %v, got %v", e.SHA, r.SHA)
				}
			}
		})
	}
}
