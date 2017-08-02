package rrstorage

// Gerneral storage wrapper
type StorageWrapper interface {
	Save([]byte, string) error // do save binary
	Fetch(string) ([]byte, error)
}
