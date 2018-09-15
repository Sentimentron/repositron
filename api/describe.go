package api

import (
	"encoding/json"
	"net/http"
)

type APIDescription struct {
	APIVersion string `json:"api_version"`
}

const APIVersion = "1.0"

func DescribeEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		err := encoder.Encode(APIDescription{APIVersion})
		if err != nil {
			panic(err)
		}
	})
}
