package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/Sentimentron/repositron/utils"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type UIFile struct {
	models.Blob
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

func createDownloadLink(id int64) string {
	return fmt.Sprintf("/v1/blobs/byId/%d/content", id)
}

func createDeleteLink(id int64) string {
	return fmt.Sprintf("delete/%d", id)
}

func formatDate(t time.Time) string {
	return t.Format("15:04:05 Jan 2 2006 MST")
}

func formatJSON(d map[string]interface{}) string {
	s, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(s)
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

			// Retrieve all the files inside this bucket
			allIds, err := store.GetBlobIdsMatchingBucket(b)
			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			allBlobs := make([]UIFile, 0)
			for _, v := range allIds {
				blob, err := store.RetrieveBlobById(v)
				if err != nil {
					log.Printf("Error retrieving blob with id %d: %v", v, err)
				} else {
					allBlobs = append(allBlobs, UIFile{*blob})
				}
			}

			cur.Contents = allBlobs
			displayData = append(displayData, cur)
		}

		// Render out a template.
		fmap := template.FuncMap{
			"formatBucketAsLink": formatBucketAsLink,
			"createDownloadLink": createDownloadLink,
			"createDeleteLink":   createDeleteLink,
			"formatDate":         formatDate,
			"formatJSON":         formatJSON,
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

func createActualDeleteLink(id int64) string {
	return fmt.Sprintf("/del/%d", id)
}

func DeleteConfirmEndpointFactory(store interfaces.MetadataStore, uiDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Convert URL parameter
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Retrieve the blob
		blob, err := store.RetrieveBlobById(id)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmap := template.FuncMap{
			"createActualDeleteLink": createActualDeleteLink,
		}

		t := template.Must(template.New("delete.html").Funcs(fmap).ParseFiles(path.Join(uiDir, "delete.html")))
		err = t.Execute(w, UIFile{*blob})
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func DeleteEndpointFactory(store interfaces.MetadataStore, contentStore interfaces.ContentStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Convert URL parameter
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		blob, err := store.RetrieveBlobById(id)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = store.DeleteBlobById(id)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = contentStore.DeleteBlobContent(blob)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 301)
	})
}

func ProcessUploadEndpointFactory(store interfaces.MetadataStore, contentStore interfaces.ContentStore) http.Handler {
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

		// Write the blob's content
		written, err := contentStore.AppendBlobContent(newBlob, &buf)
		if written.Size != newBlob.Size {
			fmt.Fprintf(w, "Did not write enough: %d out of %d byte(s), error: %v", written.Size, newBlob.Size, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if err != nil {
			fmt.Fprintf(w, "Write error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Read the blob's content back
		checksumBuffer := new(bytes.Buffer)
		read, err := contentStore.RetrieveBlobContent(newBlob, checksumBuffer)
		if read != newBlob.Size {
			fmt.Fprintf(w, "Did not read enough: %d out of %d byte(s), error: %v",
				read, newBlob.Size, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Compute the checksum of the item
		newBlob.Checksum = utils.ComputeSHA256Checksum(checksumBuffer)

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
