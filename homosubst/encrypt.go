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
	"fmt"
	"homophone/filehelper"
	"homophone/oshelper"
	"os"
	"path"
	"strings"
	"unicode"
)

// ******** Public type functions ********

// Encrypt encrypts the file named in the creation call with the built homophone substitution.
func (s *Substitutor) Encrypt(noOther bool) (string, error) {
	inFileName := s.fileName

	inFile, err := os.Open(inFileName)
	if err != nil {
		return ``, makeFileError(`open`, `in`, inFileName, err)
	}
	defer filehelper.CloseFile(inFile)

	outFileName := buildOutFileName(inFileName)
	var outFile *os.File
	outFile, err = os.OpenFile(outFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return ``, makeFileError(`open`, `out`, outFileName, err)
	}
	defer filehelper.CloseFile(outFile)

	err = s.encryptFile(inFile, outFile, noOther, outFileName)
	if err != nil {
		return ``, err
	}

	return outFileName, nil
}

// ******** Private type functions ********

// encryptFile encrypts a file.
func (s *Substitutor) encryptFile(
	inFile *os.File,
	outFile *os.File,
	noOther bool,
	outFileName string,
) error {
	var err error

	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		err = s.encryptOneLine(text, writer, noOther, outFileName)
		if err != nil {
			return err
		}

		if !noOther {
			_, err = writer.WriteString(oshelper.NewLine)
			if err != nil {
				return makeFileError(`write to`, `out`, outFileName, err)
			}
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

// encryptOneLine encrypts one line.
func (s *Substitutor) encryptOneLine(text string, writer *bufio.Writer, noOther bool, outFileName string) error {
	var err error

	for _, r := range text {
		r = unicode.ToUpper(r)
		if r >= 'A' && r <= 'Z' {
			_, err = writer.WriteRune(s.SubstituteRune(r))
		} else {
			if !noOther {
				_, err = writer.WriteRune(r)
			}
		}

		if err != nil {
			return makeFileError(`write to`, `out`, outFileName, err)
		}
	}
	return nil
}

// ******** Private functions ********

// buildOutFileName builds the file name of the output file.
func buildOutFileName(fileName string) string {
	dir := path.Dir(fileName)
	base := path.Base(fileName)
	ext := path.Ext(fileName)
	base = strings.TrimSuffix(base, ext)
	return path.Join(dir, base+"_homophone"+ext)
}

// makeFileError builds an error for a file error.
func makeFileError(operation string, direction string, fileName string, err error) error {
	return fmt.Errorf(`could not %s %sput file '%s': %w`, operation, direction, fileName, err)
}
