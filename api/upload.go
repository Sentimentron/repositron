package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func UploadDescriptionEndpointFactory(store interfaces.MetadataStore, router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse the upload content
		defer r.Body.Close()
		upload := &models.Blob{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(upload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Check the upload content
		err = upload.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Send the upload description to the store
		blob, err := store.StoreBlobRecord(upload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Generate a redirect for processing the actual upload
		redirectURL, err := router.Get("ContentUpload").URL("id", fmt.Sprintf("%d", blob.Id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		if len(redirectURL.String()) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", "cannot determine upload")
			return
		}

		// Write the response out to the client
		response := models.BlobUploadResponse{RedirectURL: redirectURL.String(), Blob: blob}
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
	})
}

func UploadContentEndpointFactory(metadataStore interfaces.MetadataStore, contentStore interfaces.ContentStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

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
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Create a teereader so we can stream the content out to disk and compute the checksum simulatenously
		h := sha256.New()
		tee := io.TeeReader(r.Body, h)
		bytesWritten, err := contentStore.WriteBlobContent(blob, tee)
		if bytesWritten != r.ContentLength {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", "didn't write enough")
			return
		}
		blob.Checksum = fmt.Sprintf("%x", h.Sum(nil))
		blob.Size = r.ContentLength

		// Finalize the upload
		_, err = metadataStore.FinalizeBlobRecord(blob)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)

	})

}
