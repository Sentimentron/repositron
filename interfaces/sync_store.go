package interfaces

type SynchronizationStore interface {

	// Holds the requesting resource until the identified resource is available.
	// Returns an error if it's not possible.
	Lock(int64) error

	// Releases the hold on the resource.
	Unlock(int64) error

}
