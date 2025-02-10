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
// Version: 2.0.0
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-02-10: V2.0.0: Print proportions, if present.
//

package homosubst

import "fmt"

// ******** Public functions ********

// Print prints all substitutions.
func (s *Substitutor) Print() {
	substitutions := s.substitutions
	proportions := s.proportions
	for i, substitution := range substitutions {
		fmt.Printf(`   %c`, i+'A')
		if proportions != nil {
			printProportion(proportions[i])
		}
		fmt.Print(`: `)

		for _, r := range substitution.BaseList() {
			fmt.Printf(`%c`, r)
		}

		fmt.Println()
	}
}

// ******** Private functions ********

// printProportion prints a proportion.
func printProportion(proportion uint16) {
	fmt.Print(` (`)
	fixProportion := proportion / 100
	fmt.Printf(`%3d`, fixProportion)
	fmt.Print(`.`)
	fmt.Printf(`%02d`, proportion-(fixProportion*100))
	fmt.Print(`%)`)
}
