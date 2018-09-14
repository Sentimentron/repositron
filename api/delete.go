package api

import (
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strconv"
)

func DeleteBlobByIdEndpointFactory(store interfaces.MetadataStore, staticDir string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse the file
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Remove the blob from the store
		err = store.DeleteBlobById(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Remove the blob from the filesystem
		uploadPath := path.Join(staticDir, fmt.Sprintf("%d", id))
		os.Remove(uploadPath)

		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

}
