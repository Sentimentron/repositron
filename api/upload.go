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
		err = enc.Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
	})
}

func AppendContentEndpointFactory(store interfaces.MetadataStore, contentStore interfaces.ContentStore, synchronizationStore interfaces.SynchronizationStore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		defer r.Body.Close()
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w,"Error: No data provided")
			return
		}

		err = synchronizationStore.Lock(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		defer synchronizationStore.Unlock(id)

		// Retrieve the blob and lock it
		blob, err := store.RetrieveBlobById(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Write the content to the end of the blob
		bytesWritten, err := contentStore.AppendBlobContent(blob, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		if bytesWritten < r.ContentLength {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: did not write enough (expected %d, got %d)", r.ContentLength, bytesWritten)
			return
		}
		blob.Size += bytesWritten
		blob.Checksum = "<recalculating>"

		// Must commit at this stage, otherwise we may experience corruption
		// if we fail to update the checksum.
		blob, err = store.FinalizeBlobRecord(blob)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error (size record stage): %v, bytesWritten=%d", err, bytesWritten)
			return
		}

		// Update the checksum
		h := sha256.New()
		bytesRead, err := contentStore.RetrieveBlobContent(blob, h)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		if bytesRead == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: did not read enough: 0 size")
			return
		}
		if bytesRead != blob.Size {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: did not read enough: expected %d, got %d", blob.Size, bytesRead)
			return
		}

		blob.Checksum = fmt.Sprintf("%x", h.Sum(nil))
		// Finalize the append
		blob, err = store.FinalizeBlobRecord(blob)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error (checksum stage): %v", err)
			return
		}

		// TODO: cleanup these duplicate statuses
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err = enc.Encode(blob)
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
