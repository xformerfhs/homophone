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
//    2025-01-03: V1.1.0: Use randomlist. Correct rune substitution handling.
//    2025-02-17: V2.0.0: Handle only bytes.
//

package homosubst

// ******** Public functions ********

// SubstituteByte substitutes a byte.
func (s *Substitutor) SubstituteByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return s.substitutions[b-'A'].RandomElement()
	} else {
		return b
	}
}
