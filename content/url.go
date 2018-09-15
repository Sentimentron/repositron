package content

import (
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/Sentimentron/repositron/utils"
	"github.com/gorilla/mux"
	"io"
	"os"
	"path"
)

type FileSystemContentStore struct {
	PrefixPath string
}

func CreateStore(staticDir string) (*FileSystemContentStore, error) {
	if !utils.IsDirectory(staticDir) {
		return nil, interfaces.BlobContentConfigError
	}

	return &FileSystemContentStore{staticDir}, nil
}

func (s *FileSystemContentStore) getPathForId(id int64) (string, error) {
	if id <= 0 {
		return "", interfaces.BlobMetadataError
	}
	return path.Join(s.PrefixPath, fmt.Sprintf("%d", id)), nil
}

func (s *FileSystemContentStore) RetrieveURLForBlobContent(m *models.Blob, r *mux.Router) (string, error) {
	url := r.Get("static")
	return fmt.Sprintf("%s/%d", url, m.Id), nil
}

func (s *FileSystemContentStore) WriteBlobContent(m *models.Blob, r io.Reader) (int64, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return -1, err
	}

	// Open the file on disk
	f, err := os.Create(p)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	// Write the content to disk
	return io.Copy(f, r)
}

func (s *FileSystemContentStore) DeleteBlobContent(m *models.Blob) error {
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return err
	}

	return os.Remove(p)
}
