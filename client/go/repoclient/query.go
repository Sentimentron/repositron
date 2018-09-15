package repoclient

import (
	"github.com/Sentimentron/repositron/models"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"strconv"
)

func (c *RepositronConnection) QueryById(blobId int64) (*models.Blob, error) {

	contentUrl := c.GetURL(fmt.Sprintf("v1/blobs/byId/%d", blobId))

	// Request the object
	resp, err := http.Get(contentUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code, expected %d, got %d", 200, resp.StatusCode)
	}

	// Decode the response
	var blobResponse models.Blob
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&blobResponse)
	if err != nil {
		return nil, err
	}

	return &blobResponse, nil
}

func (c *RepositronConnection) query(qry *models.BlobSearch) ([]int64, error) {

	contentUrl := c.GetURL("v1/blobs/search")

	// Form the request body
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(qry)
	if err != nil {
		return nil, err
	}

	// Post the query
	response, err := http.Post(contentUrl, "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decode the list of identifiers
	var ids []string
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&ids)

	ret := make([]int64, 0)
	for _, v := range ids {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}

	return ret, err
}

func (c *RepositronConnection) QueryByBucket(bucket string) ([]int64, error) {
	return c.query(&models.BlobSearch{Bucket: &bucket})
}

func (c *RepositronConnection) QueryByName(name string) ([]int64, error) {
	return c.query(&models.BlobSearch{Name: &name})
}

func (c *RepositronConnection) QueryByChecksum(checksum string) ([]int64, error) {
	return c.query(&models.BlobSearch{Checksum: &checksum})
}