package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/Sentimentron/repositron/utils"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type UIFile struct {
}

type UIBucket struct {
	Name     string
	Contents []UIFile
}

type UIData struct {
	Buckets []UIBucket
}

type UIUploadData struct {
	CurrentDate time.Time
}

func formatBucketAsLink(name string) string {
	return strings.Replace(name, " ", "_", -1)
}

func IndexEndpointFactory(store interfaces.MetadataStore, uiDir string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Retrieve a list of all buckets that are in this system.
		buckets, err := store.GetAllBuckets()
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Retrieve all the files inside these buckets and format them.
		displayData := make([]UIBucket, 0)
		for _, b := range buckets {
			cur := UIBucket{
				Name:     b,
				Contents: make([]UIFile, 0),
			}
			displayData = append(displayData, cur)
		}

		// Render out a template.
		fmap := template.FuncMap{
			"formatBucketAsLink": formatBucketAsLink,
		}
		t := template.Must(template.New("index.html").Funcs(fmap).ParseFiles(path.Join(uiDir, "index.html")))

		err = t.Execute(w, UIData{displayData})
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

}

func ProcessUploadEndpointFactory(store interfaces.MetadataStore, staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var currentBlob models.Blob

		// Process the upload metadata
		metadata := make(map[string]interface{})
		err := json.Unmarshal([]byte(r.FormValue("metadata")), &metadata)
		if err != nil {
			fmt.Fprintf(w, "Error interpreting metadata: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		currentBlob.Uploader = "Default User"
		currentBlob.Metadata = metadata
		currentBlob.Date = time.Now()
		currentBlob.Name = r.FormValue("name")
		currentBlob.Bucket = r.FormValue("bucket")
		currentBlob.Class = models.BlobType(r.FormValue("class"))

		var buf bytes.Buffer
		// Retrieve the form
		file, header, err := r.FormFile("upload")
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		io.Copy(&buf, file)

		currentBlob.Size = header.Size

		// Process validation
		validate := validator.New()
		err = validate.Struct(currentBlob)
		if err != nil {
			fmt.Fprintf(w, "Validation error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send the item to the store
		newBlob, err := store.StoreBlobRecord(&currentBlob)
		if err != nil {
			fmt.Fprintf(w, "Blob storage error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write the blob content to the output directory
		uploadPath := path.Join(staticDir, fmt.Sprintf("%d", newBlob.Id))
		f, err := os.Create(uploadPath)
		if err != nil {
			fmt.Fprintf(w, "Blob create error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bytesWritten, err := io.Copy(f, &buf)
		if int64(bytesWritten) != newBlob.Size {
			fmt.Fprintf(w, "Did not write enough: %d out of %d byte(s), error: %v", bytesWritten, newBlob.Size, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if err != nil {
			fmt.Fprintf(w, "Write error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Compute the checksum of the item
		newBlob.Checksum = utils.ComputeSHA256Checksum(&buf)

		// Finalize the store
		_, err = store.FinalizeBlobRecord(newBlob)
		if err != nil {
			fmt.Fprintf(w, "Write error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 301)

	})
}

func UploadEndpointFactory(uiDir string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmap := template.FuncMap{}
		t := template.Must(template.New("upload.html").Funcs(fmap).ParseFiles(path.Join(uiDir, "upload.html")))
		err := t.Execute(w, UIUploadData{time.Now()})
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}