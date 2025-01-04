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
//    2025-01-02: V1.1.0: Refactored for less complexity.
//

package homosubst

import (
	"bufio"
	"homophone/filehelper"
	"homophone/randomlist"
	"os"
)

// ******** Public type functions ********

// Decrypt decrypts given file with the loaded homophone substitution.
func (s *Substitutor) Decrypt(encryptedFileName string, decryptedFileName string) error {
	encryptedFile, err := os.Open(encryptedFileName)
	if err != nil {
		return makeFileError(`open`, `in`, encryptedFileName, err)
	}
	defer filehelper.CloseWithName(encryptedFile)

	var decryptedFile *os.File
	decryptedFile, err = os.OpenFile(decryptedFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return makeFileError(`open`, `out`, decryptedFileName, err)
	}
	defer filehelper.CloseWithName(decryptedFile)

	decryptionMap := buildDecryptionMap(s.substitutions)

	err = s.decryptFile(encryptedFile, decryptedFile, decryptionMap, decryptedFileName)
	if err != nil {
		return err
	}

	return nil
}

// ******** Private type functions ********

func (s *Substitutor) decryptFile(
	inFile *os.File,
	outFile *os.File,
	decryptionMap map[rune]rune,
	outFileName string,
) error {
	var err error

	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		text := scanner.Text()
		err = s.decryptOneRune(text, writer, decryptionMap, outFileName)
		if err != nil {
			return err
		}

		/*
			_, err = writer.WriteString(oshelper.NewLine)
			if err != nil {
				return makeFileError(`write to`, `out`, outFileName, err)
			}
		*/
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	err = writer.Flush()
	if err != nil {
		return makeFileError(`flush`, `out`, outFileName, err)
	}

	return nil
}

func (s *Substitutor) decryptOneRune(text string, writer *bufio.Writer, decryptionMap map[rune]rune, outFileName string) error {
	var err error

	for _, r := range text {
		decrypted, found := decryptionMap[r]
		if !found {
			decrypted = r
		}

		_, err = writer.WriteRune(decrypted)
		if err != nil {
			return makeFileError(`write to`, `out`, outFileName, err)
		}
	}

	return nil
}

// ******** Private functions ********

func buildDecryptionMap(substitutions []*randomlist.RandomList[rune]) map[rune]rune {
	result := make(map[rune]rune)
	destinationRune := 'A'
	for _, list := range substitutions {
		for _, substitution := range list.BaseList() {
			result[substitution] = destinationRune
		}
		destinationRune++
	}
	return result
}
