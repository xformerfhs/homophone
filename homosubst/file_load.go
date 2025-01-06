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
// Version: 2.1.0
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-04: V2.0.0: Restructured.
//    2025-01-05: V2.0.1: Read substitution data in Go style.
//    2025-01-06: V2.1.0: Check file header before checking file integrity.
//

package homosubst

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"homophone/compressedinteger"
	"homophone/filehelper"
	"homophone/integritycheckedfile"
	"homophone/keygenerator"
	"homophone/randomlist"
	"io"
	"os"
)

// NewLoad creates a new [Substitutor] from a substitution file.
func NewLoad(substFileName string) (*Substitutor, error) {
	var err error

	err = checkHeader(substFileName)
	if err != nil {
		return nil, err
	}

	var r *integritycheckedfile.Reader
	r, err = integritycheckedfile.NewReader(
		substFileName,
		sha3.New256,
		keygenerator.GenerateKey(generator, salt),
		additionalData)
	if err != nil {
		return nil, err
	}
	defer filehelper.CloseWithName(r)

	// Check data length.
	if r.DataLen() != substitutionDataLength {
		return nil, errors.New(`wrong file size`)
	}

	// Skip the header that has been checked at the start of this function.
	headerLen := int64(len(fileMagic)) + 1
	_, err = r.Seek(headerLen, io.SeekStart)

	// Read the rest of the file.
	substitutionData := make([]byte, int(r.DataLen()-headerLen))
	var readBytes int
	readBytes, err = r.Read(substitutionData)
	if err != nil {
		return nil, err
	}
	if readBytes != len(substitutionData) {
		return nil, errors.New(`could not read all substitution data`)
	}

	// Load substitutions from read data.
	var substitutionAlphabetSize uint32
	var substitutions []*randomlist.RandomList[rune]
	substitutionAlphabetSize, substitutions, err = loadSubstitutionData(substitutionData)
	if err != nil {
		return nil, err
	}

	return &Substitutor{
		substitutions:            substitutions,
		substitutionAlphabetSize: uint16(substitutionAlphabetSize),
	}, nil
}

// loadSubstitutionData loads all substitution data.
func loadSubstitutionData(substitutionData []byte) (uint32, []*randomlist.RandomList[rune], error) {
	var err error
	var readBytes int

	// Check size of substitution alphabet.
	var substitutionAlphabetSize uint32
	substitutionAlphabetSize, readBytes, err = compressedinteger.ToUInt32(substitutionData)
	if err != nil {
		return 0, nil, err
	}

	if substitutionAlphabetSize != requiredSubstitutionAlphabetSize {
		return 0, nil, fmt.Errorf(`wrong substitution alphabet size: %d`, substitutionAlphabetSize)
	}

	var substitutions []*randomlist.RandomList[rune]
	substitutions, err = loadSubstitutionLists(substitutionData[readBytes:], substitutionAlphabetSize)
	if err != nil {
		return 0, nil, err
	}

	return substitutionAlphabetSize, substitutions, nil
}

// loadSubstitutionLists loads all substitution lists from the substitution data.
func loadSubstitutionLists(substitutionData []byte, substitutionAlphabetSize uint32) ([]*randomlist.RandomList[rune], error) {
	var err error
	var readBytes int

	// Read all substitution lists.
	substitutions := make([]*randomlist.RandomList[rune], sourceAlphabetSize)
	check := make(map[rune]bool)
	listCount := 0
	substitutionCount := 0
	for len(substitutionData) != 0 {
		// Get size of substitution list.
		var listSize uint32
		listSize, readBytes, err = compressedinteger.ToUInt32(substitutionData)
		if err != nil {
			return nil, err
		}

		substitutionData = substitutionData[readBytes:]
		listCount++
		substitutionCount += int(listSize)

		if listCount > int(sourceAlphabetSize) {
			return nil, errors.New(`too many substitution entries`)
		}

		if substitutionCount > int(substitutionAlphabetSize) {
			return nil, errors.New(`too many substitutions`)
		}

		// Get the substitution list.
		var list []rune
		list, substitutionData, err = loadOneSubstitutionList(listSize, substitutionData, check)
		if err != nil {
			return nil, err
		}

		substitutions[listCount-1] = randomlist.New(list)
	}

	// Check list size ...
	if listCount < int(sourceAlphabetSize) {
		return nil, errors.New(`not enough substitution entries`)
	}

	// ... and number of substitutions.
	if substitutionCount < int(substitutionAlphabetSize) {
		return nil, errors.New(`not enough substitutions`)
	}

	return substitutions, nil
}

// loadOneSubstitutionList loads one substitution list from the substitution data.
func loadOneSubstitutionList(listSize uint32, substitutionData []byte, check map[rune]bool) ([]rune, []byte, error) {
	var err error

	list := make([]rune, listSize)

	var readBytes int
	var entry uint32
	for i := range listSize {
		entry, readBytes, err = compressedinteger.ToUInt32(substitutionData)
		if err != nil {
			return nil, nil, err
		}

		substitutionData = substitutionData[readBytes:]
		entryRune := rune(entry)

		if check[entryRune] {
			return nil, nil, fmt.Errorf(`duplicate substitution entry: '%c'`, entryRune)
		}

		list[i] = entryRune
		check[entryRune] = true
	}

	return list, substitutionData, nil
}

// checkHeader checks the file header.
func checkHeader(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer filehelper.CloseWithName(f)

	// Check magic bytes.
	buffer := make([]byte, len(fileMagic))
	_, err = f.Read(buffer)
	if err != nil {
		return err
	}
	if !bytes.Equal(buffer, fileMagic) {
		return errors.New(`unknown file type`)
	}

	// Check version number.
	_, err = f.Read(buffer[:1])
	if err != nil {
		return err
	}
	if buffer[0] != actVersion {
		return fmt.Errorf(`unknown file version`)
	}

	return nil
}
