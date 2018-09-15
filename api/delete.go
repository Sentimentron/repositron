package api

import (
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DeleteBlobByIdEndpointFactory(metadataStore interfaces.MetadataStore, contentStore interfaces.ContentStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse the file
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Retrieve the blob
		blob, err := metadataStore.RetrieveBlobById(id)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Remove the blob from the filesystem
		err = contentStore.DeleteBlobContent(blob)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Remove the blob from the metadataStore
		err = metadataStore.DeleteBlobById(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)

	})

}
