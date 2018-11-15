package interfaces

import (
	"errors"
	"github.com/Sentimentron/repositron/models"
	"github.com/gorilla/mux"
	"io"
)

var BlobContentNotFoundError = errors.New("blob content missing")
var BlobMetadataError = errors.New("blob metadata issue")
var BlobContentConfigError = errors.New("bad store configuration")

type ContentStore interface {
	// Removes the content associated with a blob
	DeleteBlobContent(*models.Blob) error
	// Writes content associated with a blob
	// Updates metadata at the end.
	WriteBlobContent(*models.Blob, io.Reader) (int64, error)
	// Writes content to the end of a blob
	AppendBlobContent(*models.Blob, io.Reader) (int64, error)
	// Adds content at an arbitrary position within the file
	InsertBlobContent(*models.Blob, int64, io.Reader) (int64, error)
	// Retrieves a URL to access the blob's content
	RetrieveURLForBlobContent(*models.Blob, *mux.Router) (string, error)
	// Retrieves a blob's content
	RetrieveBlobContent(*models.Blob, io.Writer) (int64, error)
}
