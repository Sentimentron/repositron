package content

import (
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"io"
	"sync"
	"sync/atomic"
)

// AccountingContentStore wraps another ContentStore, but adds tracking so that
// you can determine how much is stored there.
type AccountingContentStore struct {
	store  interfaces.EstimatableContentStore
	lock   sync.Mutex
	stored int64
}

// CreateAccountingBlobStore returns a new AccountingBlobStore.
//
// To track the amount of content stored, it needs to know how much is already
// present in the store, so it uses a contentStores's EstimateSizeOfManagedContent
// method to work that out. If an error occurs whilst calling this method,
// the size estimate is set to zero and the error is returned, alongside a
// functioning AccountingContentStore.
func CreateAccountingContentStore(underlyingStore interfaces.EstimatableContentStore) (*AccountingContentStore, error) {

	// Estimate the size of the content using the metadata store.
	storedSizeEstimate, err := underlyingStore.EstimateSizeOfManagedContent()
	if err != nil {
		storedSizeEstimate = 0
	}

	return &AccountingContentStore{
		store:  underlyingStore,
		lock:   sync.Mutex{},
		stored: storedSizeEstimate,
	}, err
}

// ContainsBlob returns whether the wrapped store contains this item.
func (a *AccountingContentStore) ContainsBlob(b *models.Blob) (bool, error) {
	return a.store.ContainsBlob(b)
}

// DeleteBlobContent removes the content from the underlying store.
func (a *AccountingContentStore) DeleteBlobContent(b *models.Blob) error {
	err := a.store.DeleteBlobContent(b)
	if err != nil {
		// The underlying store was not able to release the content.
		return err
	}
	{
		// Atomically decrement the amount of stored content.
		atomic.AddInt64(&a.stored, -b.Size)
	}
	return nil
}

// WriteBlobContent replaces or overwites the content of a given blob.
func (a *AccountingContentStore) WriteBlobContent(b *models.Blob, r io.Reader) (*models.Blob, error) {
	written, err := a.store.WriteBlobContent(b, r)
	if err != nil {
		return written, err
	}
	// Subtract the previous size given in the old blob definition, add the new size
	delta := written.Size - b.Size
	atomic.AddInt64(&a.stored, delta)
	return written, nil
}

// AppendBlobContent appends content to a given Blob, if possible.
func (a *AccountingContentStore) AppendBlobContent(b *models.Blob, r io.Reader) (*models.Blob, error) {
	// Do the underlying store thing
	written, err := a.store.AppendBlobContent(b, r)
	if err != nil {
		return written, err
	}
	atomic.AddInt64(&a.stored, written.Size-b.Size)
	return written, nil
}

// InsertBlobContent inserts content at an arbitrary offset.
func (a *AccountingContentStore) InsertBlobContent(b *models.Blob, position int64, r io.Reader) (*models.Blob, error) {
	written, err := a.store.InsertBlobContent(b, position, r)
	if err != nil {
		return written, err
	}

	delta := written.Size - b.Size
	atomic.AddInt64(&a.stored, delta)
	return written, nil
}

// RetrieveURLForBlobContent retrieves a URL from the underlying store.
func (a *AccountingContentStore) RetrieveURLForBlobContent(b *models.Blob) (string, error) {
	return a.RetrieveURLForBlobContent(b)
}

// RetrieveBlobContent retrieves the content of a Blob.
func (a *AccountingContentStore) RetrieveBlobContent(m *models.Blob, w io.Writer) (int64, error) {
	return a.store.RetrieveBlobContent(m, w)
}

// EstimatedSizeOfManagedContent - return the estimate.
func (a *AccountingContentStore) EstimatedSizeOfManagedContent() (int64, error) {
	return a.stored, nil
}
