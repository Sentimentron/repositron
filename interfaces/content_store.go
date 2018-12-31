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
var MethodNotSupportedError = errors.New("method not supported")

// ContentStore combines a separate MetadataStore and a BlobStore into
// something useful.
type ContentStore interface {
	ContainsBlob(*models.Blob) (bool, error)

	// Removes the content associated with a blob
	DeleteBlobContent(*models.Blob) error
	// WriteBlobContent replaces or creates content associated with a record.
	// i.e. if the content doesn't exist, it's created,
	// if it does exist, it's overwritten by the new content.
	// The Size field of the models.Blob argument should be ignored
	// The size returned should be the size of what was written using
	// the io.Reader.
	WriteBlobContent(*models.Blob, io.Reader) (*models.Blob, error)
	// Writes content to the end of a blob
	AppendBlobContent(*models.Blob, io.Reader) (*models.Blob, error)
	// Adds content at an arbitrary position within the file
	InsertBlobContent(*models.Blob, int64, io.Reader) (*models.Blob, error)
	// Retrieves a URL to access the blob's content
	RetrieveURLForBlobContent(*models.Blob, *mux.Router) (string, error)
	// Retrieves a blob's content
	RetrieveBlobContent(*models.Blob, io.Writer) (int64, error)
}

type EstimatableContentStore interface {
	ContentStore
	// EstimateSizeOfManagedContent returns a size estimate of the
	// amount of stuff managed in this contentStore.
	EstimateSizeOfManagedContent() (int64, error)
}

type EnumerableContentStore interface {
	ContentStore
	// RetrieveAllBlobs retrieves details of all Blobs in this ContentStore,
	// starts writing it into a channel. Blocks until all elements are wrtten
	// to the output channel.
	RetrieveAllBlobs(out chan *models.Blob) error
}
