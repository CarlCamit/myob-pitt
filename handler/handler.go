package handler

import (
	"encoding/json"
	"net/http"

	"github.com/carlcamit/myob-pitt/github"
)

func Root() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world\r\n"))
	}
}

func Health(c github.Checker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := c.CheckStatus()
		if err == nil {
			w.WriteHeader(status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		response := make(map[string]string, 1)
		response["error"] = err.Error()

		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(&response)

	}
}

type MyApplication struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	SHA         string `json:"lastcommitsha"`
}

func Metadata(c github.Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")

		response := make(map[string]interface{}, 1)

		description, err := c.GetDescription()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			response["error"] = err.Error()
			enc.Encode(&response)
			return
		}

		sha, err := c.GetLatestSHA()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			response["error"] = err.Error()
			enc.Encode(&response)
			return
		}

		myApplication := []MyApplication{
			{
				Version:     "1.0.0",
				Description: description,
				SHA:         sha,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response["myapplication"] = myApplication
		enc.Encode(&response)
	}
}
