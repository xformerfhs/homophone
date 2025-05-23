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
// Version: 3.0.1
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-02: V1.1.0: Modularized.
//    2025-01-04: V2.0.0: New interface and structure.
//    2025-01-05: V2.1.0: Fix substitution distribution. Some small refactoring.
//    2025-01-05: V2.2.0: Fix substitution alphabet.
//    2025-01-05: V2.2.1: Correct handling of additional data that are not flags.
//    2025-01-06: V2.2.2: Correct error handling for substitution files.
//    2025-01-06: V2.2.3: Rename substitution file loader.
//    2025-02-08: V2.3.0: Faster substitution calculation.
//    2025-02-09: V2.4.0: Make substitution calculation more robust.
//    2025-02-09: V2.5.0: Simplified random adjustment, corrected distributor.
//    2025-02-10: V2.6.0: Show proportions when encrypting.
//    2025-02-11: V2.7.0: Read bytes, instead of runes.
//    2025-02-17: V2.8.0: Simplify decryption.
//    2025-02-17: V3.0.0: Work only with bytes, instead of runes.
//    2025-02-24: V3.0.1: Slightly improved efficiency.
//

package main

import (
	"os"
	"unicode"
	"unicode/utf8"
)

// myVersion contains the current version of this program.
const myVersion = `3.0.1`

// myCopyright contains the copyright of this program.
const myCopyright = `Copyright (c) 2024-2025 Frank Schwab`

// ******** Main function ********

// main is the entry point of the program.
func main() {
	os.Exit(realMain(os.Args[1:]))
}

// realMain is the real main function with a return code.
func realMain(args []string) int {
	defineCommandLineFlags()

	numArgs := len(args)
	if numArgs < 1 {
		return printUsageError(`Not enough arguments`)
	}

	var rc int

	r, _ := utf8.DecodeRuneInString(args[0])
	cmd := unicode.ToUpper(r)
	switch cmd {
	case 'D':
		rc = parseDecryption()
		if rc == rcOK {
			return doDecryption(inFileName, outFileName, substFileName)
		} else {
			return rc
		}

	case 'E':
		rc = parseEncryption()
		if rc == rcOK {
			return doEncryption(inFileName, outFileName, substFileName, keepOthers)
		} else {
			return rc
		}

	case 'H':
		return printUsageOnly()

	case 'V':
		return printVersion()

	default:
		return printUsageErrorf(`Unknown command: '%s'`, args[0])
	}
}
