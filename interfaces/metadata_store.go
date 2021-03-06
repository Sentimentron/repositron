package interfaces

import (
	"errors"
	"github.com/Sentimentron/repositron/models"
)

var NoMatchingBlobsError = errors.New("no matching blobs")

type MetadataStore interface {
	// StoreBlobRecord commit WIP metadata to the database, returns a new Blob
	StoreBlobRecord(blob *models.Blob) (*models.Blob, error)
	// FinalizeBlobRecord stores the final file size,
	FinalizeBlobRecord(blob *models.Blob) (*models.Blob, error)

	// EstimateSizeOfManagedContent returns an overall size estimate for the
	// amount of stuff stored in the database.
	EstimateSizeOfManagedContent() (int64, error)

	DeleteBlobById(id int64) error
	RetrieveBlobById(id int64) (*models.Blob, error)

	GetBlobIdsMatchingChecksum(checksum string) ([]int64, error)
	GetBlobIdsMatchingName(name string) ([]int64, error)
	GetBlobIdsMatchingBucket(name string) ([]int64, error)

	// Retrieves each distinct bucket name.
	GetAllBuckets() ([]string, error)

	Close() error
}
