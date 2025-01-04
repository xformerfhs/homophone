package filehelper

// NameCloser is the interface for closing and getting a name.
type NameCloser interface {
	Close() error
	Name() string
}
