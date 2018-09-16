package database

import (
	"github.com/Sentimentron/repositron/models"
	"github.com/jmoiron/sqlx"

	"fmt"
	"github.com/Sentimentron/repositron/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	path   string
	handle *sqlx.DB
}

// CreateStore generates or opens a blob store.
func CreateStore(path string) (*Store, error) {

	// Create the store if it does not exist
	err := CreateDatabaseIfNotExists(path)
	if err != nil {
		return nil, err
	}

	// Check that it's in the right format.
	_, err = GetDatabaseSchemaVersion(path)
	if err != nil {
		return nil, err
	}

	// Open the store for real this time
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Store{path, db}, nil
}

// Close disposes of the store and any underlying resources.
func (s *Store) Close() error {
	return s.handle.Close()
}

// StoreBlobRecord inserts a WIP-blob into the database and allocates an id.
// Specifically, it stores the name, bucket, class, uploader, and metadata.
func (s *Store) StoreBlobRecord(blob *models.Blob) (*models.Blob, error) {

	sql := `
		INSERT INTO blobs (name, bucket, class, uploader, metadata, date, sha1, size) 
		VALUES (:name, :bucket, :class, :uploader, :metadata, :date, :sha1, :size)
`
	result, err := s.handle.NamedExec(sql, blob)
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.RetrieveBlobById(newId)
}

// FinalizeBlobRecord completes a blob and records all fields.
func (s *Store) FinalizeBlobRecord(blob *models.Blob) (*models.Blob, error) {

	sql := `
		UPDATE blobs SET 
			name = :name, 
			bucket = :bucket, 
			date = :date,
			class = :class,
			sha1 = :sha1, 
			uploader = :uploader, 
			metadata = :metadata,
			size = :size
		WHERE
			id = :id
	`

	if blob.Checksum == "" {
		return nil, fmt.Errorf("required finalization data missing: checksum")
	}
	if blob.Bucket == "" {
		return nil, fmt.Errorf("required finalization data missing: bucket")
	}
	if blob.Name == "" {
		return nil, fmt.Errorf("required finalization data missing: name")
	}
	if blob.Class == "" {
		return nil, fmt.Errorf("required finalization data missing: class")
	}
	if blob.Uploader == "" {
		return nil, fmt.Errorf("required finalization data missing: uploader")
	}
	if blob.Metadata == nil {
		return nil, fmt.Errorf("required finalization data missing: metadata")
	}
	if blob.Size == 0 {
		return nil, fmt.Errorf("required finalization data missing: size")
	}

	if blob.Checksum == "" || blob.Bucket == "" || blob.Name == "" || blob.Class == "" || blob.Uploader == "" || blob.Metadata == nil || blob.Size == 0 {
		return nil, fmt.Errorf("required finalization data missing")
	}

	// Process the update
	_, err := s.handle.NamedExec(sql, blob)
	if err != nil {
		return nil, err
	}

	// Retrieve the new blob
	return s.RetrieveBlobById(blob.Id)
}

// RetrieveBlobById returns a blob record from the database with a given ID.
func (s *Store) RetrieveBlobById(id int64) (*models.Blob, error) {
	ret := make([]models.Blob, 0)
	err := s.handle.Select(&ret, "SELECT id, name, bucket, date, class, sha1, uploader, metadata, size FROM blobs WHERE id = :id", id)
	if err != nil {
		return nil, fmt.Errorf("RetrieveBlobsById: %v", err)
	}
	if len(ret) == 0 {
		return nil, interfaces.NoMatchingBlobsError
	}
	return &ret[0], nil
}

// GetBlobIdsMatchingChecksum retrieves a list of blobs which match a given SHA1.
func (s *Store) GetBlobIdsMatchingChecksum(checksum string) ([]int64, error) {
	ret := make([]int64, 0)
	err := s.handle.Select(&ret, "SELECT id FROM blobs WHERE sha1 = $1", checksum)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, interfaces.NoMatchingBlobsError
	}
	return ret, nil
}

// GetBlobIdsMatchingName retrieves a list of blobs which match a name.
func (s *Store) GetBlobIdsMatchingName(name string) ([]int64, error) {
	ret := make([]int64, 0)
	err := s.handle.Select(&ret, "SELECT id FROM blobs WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, interfaces.NoMatchingBlobsError
	}
	return ret, nil
}

// GetBlobIdsMatchingBucket retrieves a list of blobs which match a bucket.
func (s *Store) GetBlobIdsMatchingBucket(bucket string) ([]int64, error) {
	ret := make([]int64, 0)
	err := s.handle.Select(&ret, "SELECT id FROM blobs WHERE bucket = $1", bucket)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, interfaces.NoMatchingBlobsError
	}
	return ret, nil
}

// DeleteBlobById deletes a record.
func (s *Store) DeleteBlobById(id int64) error {
	// Process the update
	_, err := s.handle.Exec(`DELETE FROM blobs WHERE id = $1`, id)
	return err
}

// GetAllBuckets retrieves a list of all the available buckets
func (s *Store) GetAllBuckets() ([]string, error) {
	ret := make([]string, 0)
	err := s.handle.Select(&ret, "SELECT DISTINCT bucket FROM blobs")
	return ret, err
}
