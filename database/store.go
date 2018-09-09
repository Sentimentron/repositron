package database

import (
	"errors"
	"github.com/Sentimentron/repositron/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

var NoMatchingBlobsError = errors.New("no matching blobs")

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
		INSERT INTO blobs (name, bucket, class, uploader, metadata, date) 
		VALUES (:name, :bucket, :class, :uploader, :metadata, :date)
`
	result, err := s.handle.NamedExec(sql, blob)
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	retBlob := *blob
	retBlob.Id = newId

	return &retBlob, err
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
	err := s.handle.Select(&ret, "SELECT * FROM blobs WHERE id = :id", id)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, NoMatchingBlobsError
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
		return nil, NoMatchingBlobsError
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
		return nil, NoMatchingBlobsError
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
		return nil, NoMatchingBlobsError
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
