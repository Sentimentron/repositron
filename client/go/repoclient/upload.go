package repoclient

import (
	"fmt"
	"github.com/Sentimentron/repositron/models"
	"io"
	"encoding/json"
	"bytes"
	"net/http"
	"log"
	"os"
	"io/ioutil"
)

type UploadReader struct {
	underlyingReader io.Reader
	lastReportedPercentage int
	uploadedSoFar int
	uploadTotal int64
}

func NewUploadReader(r io.Reader, size int64) *UploadReader {
	return &UploadReader{r, 0, 0, size}
}

func (u *UploadReader) Read(p []byte) (int, error) {

	// Read from the underlying reader
	r, err := u.underlyingReader.Read(p)
	if err != nil && err != io.EOF {
		log.Printf("Error reading: %v", err)
		return r, err
	}

	// Update progress and total
	u.uploadedSoFar += r
	percentage := (100.0 * float32(u.uploadedSoFar)) / float32(u.uploadTotal)
	fmt.Printf("%.2f%% uploaded...\r", percentage)
	return r, err
}

func (c *RepositronConnection) Upload(b *models.Blob, r io.Reader, verbose bool) error {
	client := &http.Client{}

	if verbose {
		r = NewUploadReader(r, b.Size)
	}

	metadataUrl := c.GetURL("v1/blobs")

	if verbose {
		log.Printf("Uploading to... %s", metadataUrl)
	}

	// Form the request body
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		return err
	}

	// Post the metadata to start the upload
	metadataRequest, err := http.NewRequest("PUT", metadataUrl, &buf)
	metadataRequest.ContentLength = int64(buf.Len())
	metadataResponse, err := client.Do(metadataRequest)
	if err != nil {
		return err
	}
	defer metadataResponse.Body.Close()

	// Decode the response
	var uploadResponse models.BlobUploadResponse
	var bodyReader io.Reader
	if verbose {
		bodyReader = io.TeeReader(metadataResponse.Body, os.Stdout)
	} else {
		bodyReader = metadataResponse.Body
	}
	dec := json.NewDecoder(bodyReader)
	err = dec.Decode(&uploadResponse)
	if err != nil {
		return err
	}

	// Upload the blob content
	contentURL := c.GetRawURL(uploadResponse.RedirectURL)
	if verbose {
		log.Printf("Uploading content to... %s", contentURL)
	}
	request, err := http.NewRequest("PUT", contentURL, r)
	request.ContentLength = b.Size
	response, err := client.Do(request)
	if err != nil {
		return err
	} else {
		defer response.Body.Close()
		if response.StatusCode != 202 {
			bytes, _ := ioutil.ReadAll(response.Body)
			return fmt.Errorf("bad status code: expected 202, got: %d (%s)", response.StatusCode, bytes)
		}
	}

	if verbose {
		fmt.Print("\n")
	}

	return nil
}
