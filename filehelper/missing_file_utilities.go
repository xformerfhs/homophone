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
// Version: 1.1.0
//
// Change history:
//    2025-01-02: V1.0.0: Created.
//    2025-01-03: V1.1.0: Added "PathComponents".
//

// Package filehelper contains file utilities missing from the Go base libraries.
package filehelper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CloseWithName closes a NameCloser and prints an error message if closing failed.
func CloseWithName(c NameCloser) {
	err := c.Close()
	if err != nil {
		printCloseOperationError(c.Name(), err)
	}
}

// RealBaseName gets the base name of a file without the extension.
func RealBaseName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

// PathComponents returns the directory, base name and extension (with leading '.') of the supplied file path.
func PathComponents(filePath string) (string, string, string) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	base = strings.TrimSuffix(base, ext)
	return dir, base, ext
}

// ******** Private functions ********

// printCloseOperationError prints an error message for a file operation.
func printCloseOperationError(name string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, `Error closing '%s': %v`, name, err)
}
