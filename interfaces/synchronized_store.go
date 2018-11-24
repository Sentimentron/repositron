package interfaces

import (
	"github.com/Sentimentron/repositron/models"
	"io"
	"github.com/gorilla/mux"
	"fmt"
	"bytes"
	"github.com/Sentimentron/repositron/utils"
)

// CombinedStore takes a content store and metadata store and updates one after another.
// It complies the with the BlobStore interface
type CombinedStore struct {
	m MetadataStore
	c ContentStore
}


// CreateCombinedStore combines a metadataStore and a contentStore together.
func CreateCombinedStore(metadataStore MetadataStore, contentStore ContentStore) *CombinedStore {
	return &CombinedStore{metadataStore, contentStore}
}

// DeleteBlobContent removes a file's content and metadata from Repositron
func (c *CombinedStore) DeleteBlobContent(blob *models.Blob) error {

	// Retrieve the metadata record
	info, err := c.m.RetrieveBlobById(blob.Id)
	if err != nil {
		return err
	}

	// Delete the metadata record first
	err = c.m.DeleteBlobById(blob.Id)
	if err != nil {
		return err
	}

	// Delete the disk content
	err = c.c.DeleteBlobContent(info)
	if err != nil {
		return err
	}

	return nil
}

func (c *CombinedStore) retrieveOrStoreMetadataIfNeeded(b *models.Blob) (*models.Blob, error) {
	// Retrieve the metadata record, if it exists
	info, err := c.m.RetrieveBlobById(b.Id)
	if err != nil {
		if err == NoMatchingBlobsError {
			// Haven't written the blob into the metadata store yet...
			info, err = c.m.StoreBlobRecord(b)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return info, nil
}

func (c *CombinedStore) WriteBlobContent(b *models.Blob, in io.Reader) (*models.Blob, int64, error) {
	// Create a meta-data record, if needed
	info, err := c.retrieveOrStoreMetadataIfNeeded(b)
	if err != nil {
		return nil, -1, err
	}

	// Write the blob's content
	written, err := c.c.WriteBlobContent(info, in)
	if err != nil {
		return nil, -1, err
	} else if written != b.Size {
		return nil, written, fmt.Errorf("write: wrong amount written: %d vs %d", written, b.Size)
	}

	info.Size = written
	return c.computeAndStoreChecksum(info)
}

func (c *CombinedStore) computeAndStoreChecksum(b *models.Blob) (*models.Blob, int64, error) {

	// Retrieve the blob's content again
	var buffer bytes.Buffer
	read, err := c.c.RetrieveBlobContent(b, &buffer)
	if read != b.Size {
		return nil, -1, fmt.Errorf("checksum: not enough read: %d vs %d", read, b.Size)
	}

	// Compute the item's checksum
	reader := bytes.NewReader(buffer.Bytes())
	b.Checksum = utils.ComputeSHA256Checksum(reader)

	// Finalize the checksum
	ret, err := c.m.FinalizeBlobRecord(b)
	if err != nil {
		return nil, -1, err
	}

	return ret, read, nil

}

func (c *CombinedStore) AppendBlobContent(b *models.Blob, reader io.Reader) (*models.Blob, int64, error) {

	// Create a blank metadata record, if necessary
	info, err := c.retrieveOrStoreMetadataIfNeeded(b)
	if err != nil {
		return nil, -1, err
	}

	// Append the content
	written, err := c.c.AppendBlobContent(info, reader)
	if err != nil {
		return nil, -1, err
	}

	info.Size += written
	return c.computeAndStoreChecksum(info)
}

func (c *CombinedStore) InsertBlobContent(b *models.Blob, offset int64, buf io.Reader) (*models.Blob, int64, error) {

	// Create a blank metadata record, if necessary
	info, err := c.retrieveOrStoreMetadataIfNeeded(b)
	if err != nil {
		return nil, -1, err
	}

	// Insert the content into storage
	written, err := c.c.InsertBlobContent(info, offset, buf)
	if err != nil {
		return nil, -1, err
	}

	newMaxOffset := offset + written
	if newMaxOffset > info.Size {
		info.Size = newMaxOffset
	}

	return c.computeAndStoreChecksum(info)
}

func (c *CombinedStore) RetrieveURLForBlobContent(b *models.Blob, r *mux.Router) (string, error) {
	return c.c.RetrieveURLForBlobContent(b, r)
}

func (c *CombinedStore) RetrieveBlobContent(b *models.Blob, w io.Writer) (int64, error) {
	return c.c.RetrieveBlobContent(b, w)
}