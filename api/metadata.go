package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"net/http"
)

func ListAllBlobsEndpointFactory(store interfaces.MetadataStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		buckets, err := store.GetAllBuckets()
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allBlobs := make([]*models.Blob, 0)

		for _, b := range buckets {
			ids, err := store.GetBlobIdsMatchingBucket(b)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, v := range ids {
				cur, err := store.RetrieveBlobById(v)
				if err != nil {
					fmt.Fprintf(w, "Error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				allBlobs = append(allBlobs, cur)
			}
		}

		w.Header().Add("Content-Type", "application/json")
		jsonMarshaller := json.NewEncoder(w)
		err = jsonMarshaller.Encode(allBlobs)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

}

func GetBlobDescriptionByIdEndpointFactory(store interfaces.MetadataStore) http.Handler {

	return nil

}
