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
	"fmt"
	"os"
)

// ******** Private constants ********

const (
	rcOK  = 0
	rcErr = 1
)

// ******** Private functions ********

func printErrorf(msgFormat string, args ...any) int {
	return printError(fmt.Sprintf(msgFormat, args...))
}

func printError(msg string) int {
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, msg)

	return rcErr
}

func printErrorfUsage(msgFormat string, args ...any) int {
	return printErrorUsage(fmt.Sprintf(msgFormat, args...))
}

func printErrorUsage(msg string) int {
	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, msg)

	return printUsage()
}

func printUsage() int {
	_, _ = fmt.Fprintln(os.Stderr, "\nUsage:\n   encrypt {file} [{noOther}]\n   decrypt {file} {subst file}")
	return rcErr
}
