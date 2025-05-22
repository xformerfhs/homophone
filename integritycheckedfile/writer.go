//
// SPDX-FileCopyrightText: Copyright 2024-2025 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.1.0
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-03: V1.1.0: Add Name function.
//

package integritycheckedfile

import (
	"crypto/hmac"
	"hash"
	"homophone/slicehelper"
	"os"
	"slices"
)

// ******** Public types ********

// Writer implements a writer for an integrity-checked file.
type Writer struct {
	file           *os.File
	hasher         hash.Hash
	additionalData []byte
}

// ******** Public creation functions ********

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

// ******** Public functions ********

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

// Name returns the name of the underlying file.
func (w *Writer) Name() string {
	return w.file.Name()
}
