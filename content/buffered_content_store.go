package content

import (
	"bytes"
	"errors"
	"github.com/Sentimentron/repositron/interfaces"
	"github.com/Sentimentron/repositron/models"
	"github.com/gorilla/mux"
	"io"
	"runtime"
	"sort"
	"sync"
	"time"
)

var NotCachedError = errors.New("not in cache")

type BufferedRequestType int

type cacheRecord struct {
	identifier   models.BlobIdentifier
	lastAccessed time.Time
	content      []byte
}

type cacheRecordSortedByAge []*cacheRecord

func (c cacheRecordSortedByAge) Len() int { return len(c) }
func (c cacheRecordSortedByAge) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c cacheRecordSortedByAge) Less(i, j int) bool {
	return c[i].lastAccessed.Sub(c[j].lastAccessed).Nanoseconds() < 0
}

// ReadHeavyBufferedContentSTore is a ContentStore that writes to disk immediately,
// but also maintains a least-recently-used cache of file contents.
type ReadHeavyBufferedContentStore struct {
	sync.RWMutex

	// Sets the maximum size of this cache. If the cache contains more items
	// than this, it will block new reads until the cache is back to the correct size.
	maximumSize int64

	// Sets the intended size of the cache. Beyond this point, the cache will
	// start dropping the least recently accessed items.
	stretchSize int64

	// The underlying ContentStore we're reading from.
	underlyingStore interfaces.BlobStore

	// A map of keys-bytes
	cache map[models.BlobIdentifier]*cacheRecord

	// Maintains the current size of the map
	storedSize int64

	// If true, indicates that there's a go-routine waiting to
	// clean-up any items.
	maintenanceScheduled bool
}

// DeleteBlobContent removes the selected blob from underlying storage and
// removes the blob from the cache too.
func (r *ReadHeavyBufferedContentStore) DeleteBlobContent(m *models.Blob) error {
	r.deleteItemFromCache(m)
	return r.underlyingStore.DeleteBlobContent(m)
}

// WriteBlobContent dumps any current record from the cache and dispatches the write
// to the store.
func (r *ReadHeavyBufferedContentStore) WriteBlobContent(m *models.Blob, ri io.Reader) (*models.Blob, int64, error) {
	r.deleteItemFromCache(m)
	return r.underlyingStore.WriteBlobContent(m, ri)
}

// AppendBlobContent dumps any current record from the cache and dispatches the write.
func (r *ReadHeavyBufferedContentStore) AppendBlobContent(m *models.Blob, ri io.Reader) (*models.Blob, int64, error) {
	r.deleteItemFromCache(m)
	return r.underlyingStore.AppendBlobContent(m, ri)
}

// InsertBlobContent dumps any current record from the cache and dispatches the write.
func (r *ReadHeavyBufferedContentStore) InsertBlobContent(m *models.Blob, pos int64, ri io.Reader) (*models.Blob, int64, error) {
	r.deleteItemFromCache(m)
	return r.underlyingStore.InsertBlobContent(m, pos, ri)
}

// RetrieveURLForBlobContent update the last accessed time and generates a link.
func (r *ReadHeavyBufferedContentStore) RetrieveURLForBlobContent(m *models.Blob, route *mux.Router) (string, error) {
	// TODO: fix me
	return "", nil
}

// RetrieveBlobContent updates the last accessed time and generates a link.
func (r *ReadHeavyBufferedContentStore) RetrieveBlobContent(m *models.Blob, w io.Writer) (int64, error) {
	checkCacheForContent := func() (int64, error) {
		r.RLock()
		defer r.RUnlock()
		if record, ok := r.cache[models.BlobIdentifier(m.Id)]; ok {
			record.lastAccessed = time.Now()
			return io.Copy(w, bytes.NewBuffer(record.content))
		}
		return -1, NotCachedError
	}

	insertIntoCache := func(c *cacheRecord) {
		r.Lock()
		defer r.Unlock()
		r.cache[c.identifier] = c
	}

	read, err := checkCacheForContent()
	if err == NotCachedError {
		newReader := &bytes.Buffer{}
		read, err = r.underlyingStore.RetrieveBlobContent(m, newReader)
		if err == nil {
			c := &cacheRecord{
				models.BlobIdentifier(m.Id),
				time.Now(),
				newReader.Bytes(),
			}
			insertIntoCache(c)
		}
	}

	return read, err
}

func (r *ReadHeavyBufferedContentStore) deleteItemFromCache(m *models.Blob) {
	r.Lock()
	defer r.Unlock()
	r.storedSize -= m.Size
	delete(r.cache, models.BlobIdentifier(m.Id))
}

// performMaintenanceIfNeeded schedules performMaintenance if it's needed.
func (r *ReadHeavyBufferedContentStore) performMaintenanceIfNeeded() {
	if r.maintenanceScheduled {
		return
	}
	if r.maintenanceScheduled {
		return
	}
	if r.storedSize > r.stretchSize {
		// If we're above the stretch size, perform maintenance asynchronously.
		r.Lock()
		go r.performMaintenance()
		r.maintenanceScheduled = true
		r.Unlock()
	} else if r.storedSize > r.maximumSize {
		// Perform maintenance synchronously
		r.Lock()
		r.maintenanceScheduled = true
		r.Unlock()
		r.performMaintenance()
	}
	r.maintenanceScheduled = true
}

// performMaintenance clears items from the cache until it's under the stretch size.
func (r *ReadHeavyBufferedContentStore) performMaintenance() {
	// Lock the cache for writing
	r.Lock()
	defer r.Unlock()

	defer func() {
		// Make sure maintenance completes at the end of this
		r.maintenanceScheduled = false
	}()

	// If we've completed maintenance already, do nothing
	if !r.maintenanceScheduled {
		return
	}

	// If we haven't yet reached the stretch size, there's no need
	// to do any maintenance
	if r.storedSize < r.stretchSize {
		return
	}

	// Sort the vector of last access times (earliest first)
	vec := make(cacheRecordSortedByAge, len(r.cache))
	cont := 0
	for key := range r.cache {
		vec[cont] = r.cache[key]
	}
	sort.Sort(vec)

	// Remove the items until we're below the stretch size
	for _, item := range vec {
		delete(r.cache, item.identifier)
		r.storedSize -= int64(len(item.content))
		if r.storedSize < r.stretchSize {
			break
		}
	}
	// Remove the vector of items we've examined
	vec = make(cacheRecordSortedByAge, 0)

	// Force a GC to clean up memory
	runtime.GC()
}
