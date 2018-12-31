package interfaces

import (
	"github.com/Sentimentron/repositron/models"
	"github.com/gorilla/mux"
	"io"
)

type BlobStore interface {

	// If false, this blob store doesn't contain the given blob,
	ContainsBlob(*models.Blob) (bool, error)

	// Removes the content associated with a blob.
	// If the content cannot be deleted, return an error and leave the
	// blob's content intact. Under no other circumstances should we
	// return an error.
	DeleteBlobContent(*models.Blob) error
	// Writes content associated with a blob
	// Updates metadata at the end.
	WriteBlobContent(*models.Blob, io.Reader) (*models.Blob, int64, error)
	// Writes content to the end of a blob
	AppendBlobContent(*models.Blob, io.Reader) (*models.Blob, int64, error)
	// Adds content at an arbitrary position within the file
	InsertBlobContent(*models.Blob, int64, io.Reader) (*models.Blob, int64, error)
	// Retrieves a URL to access the blob's content
	RetrieveURLForBlobContent(*models.Blob, *mux.Router) (string, error)
	// Retrieves a blob's content
	RetrieveBlobContent(*models.Blob, io.Writer) (int64, error)
}
