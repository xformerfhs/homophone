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
	"homophone/oshelper"
	"os"
	"path"
	"strings"
)

// ******** Public type functions ********

// Decrypt decrypts given file with the loaded homophone substitution.
func (s *Substitutor) Decrypt(inFileName string) (string, error) {
	inFile, err := os.Open(inFileName)
	if err != nil {
		return ``, makeFileError(`open`, `in`, inFileName, err)
	}
	defer filehelper.CloseFile(inFile)

	outFileName := buildDecryptOutFileName(inFileName)
	var outFile *os.File
	outFile, err = os.OpenFile(outFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return ``, makeFileError(`open`, `out`, outFileName, err)
	}
	defer filehelper.CloseFile(outFile)

	decryptionMap := buildDecryptionMap(s.substitutions)

	err = s.decryptFile(inFile, outFile, decryptionMap, outFileName)
	if err != nil {
		return ``, err
	}

	return outFileName, nil
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
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		err = s.decryptOneLine(text, writer, decryptionMap, outFileName)
		if err != nil {
			return err
		}

		_, err = writer.WriteString(oshelper.NewLine)
		if err != nil {
			return makeFileError(`write to`, `out`, outFileName, err)
		}
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

func (s *Substitutor) decryptOneLine(text string, writer *bufio.Writer, decryptionMap map[rune]rune, outFileName string) error {
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

func buildDecryptionMap(substitutions [][]rune) map[rune]rune {
	result := make(map[rune]rune)
	destinationRune := 'A'
	for _, list := range substitutions {
		for _, substitution := range list {
			result[substitution] = destinationRune
		}
		destinationRune++
	}
	return result
}

// buildDecryptOutFileName builds the file name of the output file.
func buildDecryptOutFileName(fileName string) string {
	pos := strings.LastIndex(fileName, `_homophone`)
	if pos >= 0 {
		return fileName[:pos] + `_decrypted` + fileName[pos+10:]
	} else {
		dir := path.Dir(fileName)
		base := path.Base(fileName)
		ext := path.Ext(fileName)
		base = strings.TrimSuffix(base, ext)
		return path.Join(dir, base+"_decrypted"+ext)
	}
}
