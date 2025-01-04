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
// Version: 2.0.0
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-04: V2.0.0: Restructured.
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
)

func NewLoad(substFileName string) (*Substitutor, error) {
	r, err := integritycheckedfile.NewReader(
		substFileName,
		sha3.New256,
		keygenerator.GenerateKey(generator, salt),
		additionalData)
	if err != nil {
		return nil, err
	}
	defer filehelper.CloseWithName(r)

	if r.DataLen() != substitutionDataLength {
		return nil, errors.New(`wrong file size`)
	}

	// Check magic bytes.
	buffer := make([]byte, len(fileMagic))
	_, err = r.Read(buffer)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(buffer, fileMagic) {
		return nil, errors.New(`unknown file format`)
	}

	// Check version number.
	_, err = r.Read(buffer[:1])
	if err != nil {
		return nil, err
	}
	if buffer[0] != actVersion {
		return nil, fmt.Errorf(`unknown version`)
	}

	// Read the rest of the file.
	substitutionData := make([]byte, int(r.DataLen())-len(fileMagic)-1)
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
	substitutions, err = loadSubstitutionLists(substitutionData, readBytes, substitutionAlphabetSize)
	if err != nil {
		return 0, nil, err
	}

	return substitutionAlphabetSize, substitutions, nil
}

// loadSubstitutionLists loads all substitution lists from the substitution data.
func loadSubstitutionLists(substitutionData []byte, actPos int, substitutionAlphabetSize uint32) ([]*randomlist.RandomList[rune], error) {
	var err error
	var readBytes int

	// Read all substitution lists.
	substitutions := make([]*randomlist.RandomList[rune], sourceAlphabetSize)
	check := make(map[rune]bool)
	listCount := 0
	substitutionCount := 0
	for actPos < len(substitutionData) {
		// Get size of substitution list.
		var listSize uint32
		listSize, readBytes, err = compressedinteger.ToUInt32(substitutionData[actPos:])
		if err != nil {
			return nil, err
		}

		actPos += readBytes
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
		list, actPos, err = loadOneSubstitutionList(listSize, substitutionData, check, actPos)
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
func loadOneSubstitutionList(listSize uint32, substitutionData []byte, check map[rune]bool, actPos int) ([]rune, int, error) {
	var err error

	list := make([]rune, listSize)

	var readBytes int
	var entry uint32
	for i := range listSize {
		entry, readBytes, err = compressedinteger.ToUInt32(substitutionData[actPos:])
		if err != nil {
			return nil, 0, err
		}

		actPos += readBytes
		entryRune := rune(entry)

		if check[entryRune] {
			return nil, 0, fmt.Errorf(`duplicate substitution entry: '%c'`, entryRune)
		}

		list[i] = entryRune
		check[entryRune] = true
	}

	return list, actPos, nil
}
