package integritycheckedfile

import (
	"crypto/hmac"
	"crypto/subtle"
	"errors"
	"hash"
	"io"
	"os"
)

// ======== Public types ========

// Reader implements a read for an integrity-checked file.
type Reader struct {
	file     *os.File
	dataLen  int64
	position int64
}

// ======== Public constants ========

// ErrFileCorrupt is returned when the checksum does not match the file data.
var ErrFileCorrupt = errors.New(`file is corrupt`)

// ======== Private constants ========

// readBufferSize is the buffer size used for the check of the file.
const readBufferSize = 1_024 << 2

// ======== Public creation functions ========

// NewReader creates a new integrity-checked file reader.
func NewReader(fileName string, hashFunc func() hash.Hash, key []byte, additionalData []byte) (*Reader, error) {
	// 1. Open file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	// Do not defer a close!

	// 2. Get file size.
	var fi os.FileInfo
	fi, err = file.Stat()
	if err != nil {
		return nil, err
	}

	// 3. Calculate length of data without the checksum.
	hasher := hmac.New(hashFunc, key)
	dataLen := fi.Size() - int64(hasher.Size())

	// 4. Check, if the checksum matches the file data.
	err = checkFileIntegrity(file, dataLen, hasher, additionalData)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	// 5. Reset file position.
	_, err = file.Seek(0, 0)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	// All done. Return reader.
	return &Reader{
		file:     file,
		dataLen:  dataLen,
		position: 0,
	}, nil
}

// ======== Public functions ========

// DataLen returns the length of the data in the file.
func (r *Reader) DataLen() int64 {
	return r.dataLen
}

// Read reads data from the file into the provided buffer.
func (r *Reader) Read(p []byte) (n int, err error) {
	// 1. If the position is at or beyond the data length, return EOF.
	if r.position >= r.dataLen {
		return 0, io.EOF
	}

	// 2. Read only as many data as are indicated by [dataLen].
	pLen := len(p)
	newPosition := r.position + int64(pLen)
	if newPosition > r.dataLen {
		pLen = int(r.dataLen - r.position)
	}

	n, err = r.file.Read(p[:pLen])
	if err != nil {
		return
	}

	r.position += int64(n)

	return
}

// Close closes the file.
func (r *Reader) Close() error {
	// 1. Close the file.
	err := r.file.Close()
	if err != nil {
		return err
	}

	// 2. Destroy all data in the reader struct.
	r.file = nil
	r.position = 0
	r.dataLen = 0

	return nil
}

// ======== Private functions ========

// checkFileIntegrity checks if the file has the correct checksum.
func checkFileIntegrity(file *os.File, dataLength int64, hasher hash.Hash, additionalData []byte) error {
	buffer := make([]byte, readBufferSize)

	var n int
	var err error

	toRead := dataLength
	for toRead > 0 {
		if toRead < readBufferSize {
			n, err = file.Read(buffer[:toRead])
		} else {
			n, err = file.Read(buffer)
		}
		if err != nil {
			return err
		}

		_, err = hasher.Write(buffer[:n])
		if err != nil {
			return err
		}

		toRead -= int64(n)
	}

	_, err = hasher.Write(additionalData)
	if err != nil {
		return err
	}

	hashSize := hasher.Size()
	n, err = file.Read(buffer[:hashSize])
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(hasher.Sum(nil), buffer[:hashSize]) == 0 {
		return ErrFileCorrupt
	}

	return nil
}
