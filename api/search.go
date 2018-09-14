package api

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
)

type BlobSearch struct {
	Name *string `json:"name"`
	Checksum *string `json:"checksum"`
	Bucket *string `json:"bucket"`
}

func SearchBlobEndpointFactory(store interfaces.MetadataStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		// Parse the input
		defer r.Body.Close()
		var qry BlobSearch

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&qry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Specify that something must be matched
		if qry.Name == nil && qry.Checksum == nil && qry.Bucket == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "name, checksum, bucket, not specified")
			return
		}

		// This set will contain all matching results
		matchingSet := make(map[int64]struct{})
		// Search for matching names
		if qry.Name != nil {
			ids, err := store.GetBlobIdsMatchingName(*qry.Name)
			if err != nil {
				if err != interfaces.NoMatchingBlobsError {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error: %v", err)
					return
				}
			} else {
				for _, v := range ids {
					matchingSet[v] = struct{}{}
				}
			}
		}
		// Search for matching buckets
		if qry.Bucket != nil {
			ids, err := store.GetBlobIdsMatchingBucket(*qry.Bucket)
			if err != nil {
				if err != interfaces.NoMatchingBlobsError{
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error: %v", err)
					return
				}
			} else {
				for _, v := range ids {
					matchingSet[v] = struct{}{}
				}
			}
		}
		// Search for matching checksums
		if qry.Checksum != nil {
			ids, err := store.GetBlobIdsMatchingChecksum(*qry.Checksum)
			if err != nil {
				if err != interfaces.NoMatchingBlobsError{
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error: %v", err)
					return
				}
			} else {
				for _, v := range ids {
					matchingSet[v] = struct{}{}
				}
			}
		}

		// Format a list of results
		ret := make([]int64, 0)
		for v := range matchingSet {
			ret = append(ret, v)
		}

		// Encode the list to JSON
		encoder := json.NewEncoder(w)
		err = encoder.Encode(ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

	})


}