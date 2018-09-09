package interfaces

import "github.com/Sentimentron/repositron/models"

type MetadataStore interface {
	// StoreBlobRecord commit WIP metadata to the database, returns a new Blob
	StoreBlobRecord(blob *models.Blob) (*models.Blob, error)
	// FinalizeBlobRecord stores the final file size,
	FinalizeBlobRecord(blob *models.Blob) (*models.Blob, error)

	DeleteBlobById(id int64) error
	RetrieveBlobById(id int64) (*models.Blob, error)

	GetBlobIdsMatchingChecksum(checksum string) ([]int64, error)
	GetBlobIdsMatchingName(name string) ([]int64, error)
	GetBlobIdsMatchingBucket(name string) ([]int64, error)

	// Retrieves each distinct bucket name.
	GetAllBuckets() ([]string, error)

	// Retrieves every blob

	Close() error
}
