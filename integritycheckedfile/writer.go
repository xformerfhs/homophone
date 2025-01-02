package integritycheckedfile

import (
	"crypto/hmac"
	"hash"
	"homophone/slicehelper"
	"os"
	"slices"
)

// ======== Public types ========

// Writer implements a writer for an integrity-checked file.
type Writer struct {
	file           *os.File
	hasher         hash.Hash
	additionalData []byte
}

// ======== Public creation functions ========

// NewWriter creates a new writer for an integrity-checked file.
func NewWriter(fileName string, hashFunc func() hash.Hash, key []byte, additionalData []byte) (*Writer, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return &Writer{
		file:           file,
		hasher:         hmac.New(hashFunc, key),
		additionalData: slices.Clone(additionalData),
	}, nil
}

// ======== Public functions ========

// Write writes the supplied data to the file.
func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.hasher.Write(p)
	if err != nil {
		return
	}

	return w.file.Write(p)
}

// WriteString writes a string to the file.
func (w *Writer) WriteString(s string) (n int, err error) {
	n, err = w.hasher.Write([]byte(s))
	if err != nil {
		return
	}

	return w.file.WriteString(s)
}

// Close closes the file.
func (w *Writer) Close() error {
	// 1. Hash additional data.
	hasher := w.hasher
	_, err := hasher.Write(w.additionalData)
	if err != nil {
		return err
	}

	// 2. Write checksum after data.
	file := w.file
	_, err = file.Write(hasher.Sum(nil))
	if err != nil {
		return err
	}

	// 3. Close the file.
	err = w.file.Close()
	if err != nil {
		return err
	}

	// 4. Destroy all data in the [Writer] struct.
	w.file = nil
	w.hasher = nil
	slicehelper.ClearNumber(w.additionalData)

	return nil
}
