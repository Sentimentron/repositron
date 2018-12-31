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
	url, err := r.Get("static").URL()
	return fmt.Sprintf("%s/%d", url, m.Id), err
}

func (s *FileSystemContentStore) WriteBlobContent(m *models.Blob, r io.Reader) (*models.Blob, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return nil, err
	}

	// Open the file on disk
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Write the content to disk
	written, err := io.Copy(f, r)

	// Update and return the blob reference
	if err != nil {
		return nil, err
	}
	ret := *m
	ret.Size = written
	return &ret, err
}

func (s *FileSystemContentStore) DeleteBlobContent(m *models.Blob) error {
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return err
	}

	return os.Remove(p)
}

func (s *FileSystemContentStore) AppendBlobContent(m *models.Blob, r io.Reader) (*models.Blob, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return nil, err
	}

	// Open for appending
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// Append the content
	appended, err := io.Copy(f, r)
	if err != nil {
		return nil, err
	}

	// Update the blob record and return
	ret := *m
	ret.Size += appended
	return &ret, err
}

func (s *FileSystemContentStore) ContainsBlob(m *models.Blob) (bool, error) {
	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return false, err
	}

	// Open for inserting
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *FileSystemContentStore) InsertBlobContent(m *models.Blob, offset int64, r io.Reader) (*models.Blob, error) {

	// Generate filesystem path
	p, err := s.getPathForId(m.Id)
	if err != nil {
		return nil, err
	}

	// Open for inserting
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Check how large the file is, it may require enlargement
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() < offset {
		f.Truncate(offset)
	}

	// Seek to the required offset and start writing
	newOffset, err := f.Seek(offset, 0)
	if err != nil {
		return nil, err
	} else if newOffset != offset {
		return nil, errors.New("offsets did not match")
	}

	// Copy the the area to the right offset
	inserted, err := io.Copy(f, r)
	if err != nil {
		return nil, err
	}

	finalOffset := newOffset + inserted

	// Update the return value
	newSize := finalOffset
	if info.Size() > finalOffset {
		newSize = info.Size()
	}
	ret := *m
	ret.Size = newSize
	return &ret, err

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
