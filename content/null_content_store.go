package content

import (
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/gorilla/mux"
	"io"
)

// NullContentStore swallows all writes but does conform to the ContentStore interface.
// It's useful for when you're chaining things.
type NullContentStore struct {
}

func (n *NullContentStore) ContainsBlob(m *models.Blob) (bool, error) {
	return false, nil
}

func (n *NullContentStore) DeleteBlobContent(m *models.Blob) error {
	return nil
}

func (n *NullContentStore) WriteBlobContent(*models.Blob, io.Reader) (*models.Blob, error) {
	return nil, nil
}

func (n *NullContentStore) AppendBlobContent(*models.Blob, io.Reader) (*models.Blob, error) {
	return nil, nil
}

func (n *NullContentStore) InsertBlobContent(*models.Blob, int64, io.Reader) (*models.Blob, error) {
	return nil, nil
}

func (n *NullContentStore) RetrieveURLForBlobContent(*models.Blob, *mux.Router) (string, error) {
	return "", interfaces.MethodNotSupportedError
}

func (n *NullContentStore) RetrieveBlobContent(blob *models.Blob, w io.Writer) (int64, error) {
	return 0, interfaces.MethodNotSupportedError
}
