package content

import (
	"github.com/Sentimentron/repositron/models"
	"io"
	"github.com/gorilla/mux"
	"github.com/Sentimentron/repositron/interfaces"
)

// NullContentStore swallows all writes but does conform to the ContentStore interface.
// It's useful for when you're chaining things.
type NullContentStore struct {

}

func (n *NullContentStore) DeleteBlobContent(m *models.Blob) error {
	return nil
}

func (n *NullContentStore) WriteBlobContent(*models.Blob, io.Reader) (int64, error) {
	return 0, nil
}

func (n *NullContentStore) AppendBlobContent(*models.Blob, io.Reader) (int64, error) {
	return 0, nil
}

func (n *NullContentStore) InsertBlobContent(*models.Blob, int64, io.Reader) (int64, error) {
	return 0, nil
}

func (n *NullContentStore) RetrieveURLForBlobContent(*models.Blob, *mux.Router) (string, error) {
	return "", interfaces.MethodNotSupportedError
}

func (n *NullContentStore) RetrieveBlobContent(blob *models.Blob, io.Writer) (int64, error) {
	return 0, interfaces.MethodNotSupportedError
}