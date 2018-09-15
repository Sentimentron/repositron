package api

import (
	"encoding/json"
	"github.com/Sentimentron/repositron/models"
	"net/http"
)

const APIVersion = "1.0"

func DescribeEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		err := encoder.Encode(models.APIDescription{APIVersion})
		if err != nil {
			panic(err)
		}
	})
}
