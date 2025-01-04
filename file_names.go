//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
// Version: 1.0.0
//
// Change history:
//    2025-01-02: V1.0.0: Created.
//

package main

import (
	"homophone/filehelper"
	"path/filepath"
	"strings"
)

// ******** Private constants ********

// decryptedMarker is the last part of the default file name for an decrypted file.
const decryptedMarker = `_decrypted`

// homophoneMarker is the last part of the default file name for an encrypted file.
const homophoneMarker = `_homophone`

// ******** Public functions ********

// buildDecryptOutFilePath builds the file path of the decrypted output file.
func buildDecryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, decryptedMarker)
}

// buildEncryptOutFilePath builds the file path of the encrypted output file.
func buildEncryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, homophoneMarker)
}

// buildSubstFilePath builds the file path of the substitution file.
func buildSubstFilePath(filePath string) string {
	dir, base, ext := cleanPathComponents(filePath)
	
	if len(ext) > 0 {
		base = strings.TrimSuffix(base, ext) + `_` + ext[1:]
	}

	return filepath.Join(dir, base+`.subst`)
}

// ******** Private functions ********

// cleanPathComponents returns the directory, base name and extension
// with the program's markers removed from base name.
func cleanPathComponents(filePath string) (string, string, string) {
	dir, base, ext := filehelper.PathComponents(filePath)

	base = strings.TrimSuffix(base, homophoneMarker)
	base = strings.TrimSuffix(base, decryptedMarker)

	return dir, base, ext
}

// buildFilePathWithMarker builds a file path with a marker that is separated by '_' after the base name.
func buildFilePathWithMarker(filePath string, marker string) string {
	dir, base, ext := cleanPathComponents(filePath)

	return filepath.Join(dir, base+marker+ext)
}
