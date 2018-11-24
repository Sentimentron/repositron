package interfaces

import (
	"github.com/Sentimentron/repositron/models"
	"io"
	"github.com/gorilla/mux"
)

type BlobStore interface {

	// Removes the content associated with a blob
	DeleteBlobContent(models.Blob) error
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
