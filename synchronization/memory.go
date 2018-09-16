package synchronization

import (
	"time"
	"sync"
	"log"
)

const tidyDuration = 5 * time.Second

type syncRecord struct {
	lastAcquired time.Time
	mutex *sync.Mutex
}

// MemorySynchronizationStore is a single-machine way of making sure
// that multiple people don't append to the same file at once.
type MemorySynchronizationStore struct {
	lockMap map[int64]*syncRecord
	lock sync.Mutex
}

// CreateMemorySynchronizationStore initializes a store.
func CreateMemorySynchronizationStore() (*MemorySynchronizationStore, error) {
	return &MemorySynchronizationStore{
		make(map[int64]*syncRecord),
		sync.Mutex{},
	}, nil
}

func (m *MemorySynchronizationStore) createRecordIfNotExists(id int64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.lockMap[id]; !ok {
		m.lockMap[id] = &syncRecord{time.Now(), &sync.Mutex{}}
	}

	return nil
}

// Lock blocks the calling thread until the identified Blob is available.
func (m *MemorySynchronizationStore) Lock(id int64) error {

	if _, ok := m.lockMap[id]; !ok {
		err := m.createRecordIfNotExists(id)
		if err != nil {
			return err
		}
	}

	// Update this field just to show that the lock is fresh.
	m.lockMap[id].lastAcquired = time.Now()

	// Acquire the lock or block trying
	m.lockMap[id].mutex.Lock()
	m.lockMap[id].lastAcquired = time.Now()

	return nil
}

func (m *MemorySynchronizationStore) tidy() {
	m.lock.Lock()
	defer m.lock.Unlock()

	cleanupMap := make(map[int64]struct{})

	currentTime := time.Now()
	for v := range m.lockMap {
		r := m.lockMap[v]
		if currentTime.Sub(r.lastAcquired) > tidyDuration {
			cleanupMap[v] = struct{}{}
		}
	}

	for v := range cleanupMap {
		delete(m.lockMap, v)
	}
}

// Unlock releases the lock on the underlying resource
func (m *MemorySynchronizationStore) Unlock(id int64) error {

	if _, ok := m.lockMap[id]; !ok {
		log.Printf("MemorySychronizationStore: integrity error for %d", id)
		return nil
	}

	m.lockMap[id].mutex.Unlock()
	m.lockMap[id].lastAcquired = time.Now()

	// Schedule a tidy operation
	defer m.tidy()

	return nil
}