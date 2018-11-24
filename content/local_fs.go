package content

import (
	"errors"
	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/Sentimentron/repositron/utils"
	"github.com/gorilla/mux"
	"io"
	"os"
	"path"
)

// FileSystemContentStore lives in a local directory on this machine.
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
	return fmt.Sprintf("%s/%d", url.String(), m.Id), nil
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

func (s *FileSystemContentStore) AppendBlobContent(m *models.Blob, r io.Reader) (int64, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return -1, err
	}

	// Open for appending
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return -1, err
	}

	defer f.Close()

	return io.Copy(f, r)
}

func (s *FileSystemContentStore) InsertBlobContent(m *models.Blob, offset int64, r io.Reader) (int64, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return -1, err
	}

	// Open for inserting
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	// Check how large the file is, it may require enlargment
	info, err := f.Stat()
	if err != nil {
		return -1, err
	}
	if info.Size() < offset {
		f.Truncate(offset)
	}

	// Seek to the required offset and start writing
	newOffset, err := f.Seek(offset, 0)
	if err != nil {
		return -1, err
	} else if newOffset != offset {
		return -1, errors.New("offsets did not match")
	}

	return io.Copy(f, r)
}

func (s *FileSystemContentStore) RetrieveBlobContent(m *models.Blob, w io.Writer) (int64, error) {
	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return -1, err
	}

	// Open for reading
	f, err := os.OpenFile(p, os.O_RDONLY, 0600)
	if err != nil {
		return -1, err
	}

	defer f.Close()
	return io.Copy(w, f)
}
