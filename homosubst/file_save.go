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
// Version: 2.0.1
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-04: V2.0.0: Restructured.
//    2025-01-05: V2.0.1: Use interface, instead of type.
//

package homosubst

import (
	"golang.org/x/crypto/sha3"
	"homophone/compressedinteger"
	"homophone/filehelper"
	"homophone/integritycheckedfile"
	"homophone/keygenerator"
	"homophone/randomlist"
	"io"
)

// Save saves substitution data to a substitution file.
func (s *Substitutor) Save(filePath string) error {
	w, err := integritycheckedfile.NewWriter(
		filePath,
		sha3.New256,
		keygenerator.GenerateKey(generator, salt),
		additionalData)
	if err != nil {
		return err
	}
	defer filehelper.CloseWithName(w)

	// Write magic bytes.
	_, err = w.Write(fileMagic)
	if err != nil {
		return err
	}

	// Write version.
	_, err = w.Write([]byte{actVersion})
	if err != nil {
		return err
	}

	// Write size of substitution alphabet.
	var size []byte
	size, err = compressedinteger.FromUInt32(uint32(s.substitutionAlphabetSize))
	if err != nil {
		return err
	}
	_, err = w.Write(size)
	if err != nil {
		return err
	}

	// Save substitution lists.
	err = saveSubstitutions(w, s.substitutions, size)

	return err
}

// saveSubstitutions saves the substitution lists.
func saveSubstitutions(w io.Writer, substitutions []*randomlist.RandomList[byte], size []byte) error {
	var err error

	// Save all substitution lists.
	for _, substitutionList := range substitutions {
		// Write length of substitution list.
		size, err = compressedinteger.FromUInt32(uint32(substitutionList.Len()))
		if err != nil {
			return err
		}

		_, err = w.Write(size)
		if err != nil {
			return err
		}

		// Write each substitution character.
		for _, r := range substitutionList.BaseList() {
			size, err = compressedinteger.FromUInt32(uint32(r))
			if err != nil {
				return err
			}

			_, err = w.Write(size)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
