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
//    2025-01-02: V1.1.0: Modularized.
//

package main

import (
	"os"
	"strings"
)

// ******** Main function ********

// main is the entry point of the program.
func main() {
	os.Exit(realMain(os.Args[1:]))
}

// realMain is the real main function with a return code.
func realMain(args []string) int {
	numArgs := len(args)
	if numArgs < 1 {
		return printErrorUsage(`Not enough arguments`)
	}

	cmd := strings.ToLower(args[0])
	switch cmd[0] {
	case 'e':
		return doEncryption(args[1], numArgs > 1)

	case 'd':
		if len(args) < 3 {
			return printErrorUsage(`Not enough arguments`)
		}
		return doDecryption(args[1], args[2])

	default:
		return printErrorfUsage(`Unknown command: '%s'`, cmd)
	}
}
