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

// Package filehelper contains file utilities missing from the Go base libraries.
package filehelper

import (
	"fmt"
	"os"
)

// CloseFile closes a file and reports an error, if it occurs.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error closing file '%s': %v\n", file.Name(), err)
	}
}
