//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
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
//    2024-09-17: V1.0.0: Created.
//

package homosubst

func (s *Substitutor) SubstituteByte(b byte) rune {
	if b >= 'A' && b <= 'Z' {
		return s.substituteByte(b)
	} else {
		return rune(b)
	}
}

func (s *Substitutor) SubstituteRune(r rune) rune {
	return s.substituteByte(byte(r))
}

func (s *Substitutor) substituteByte(b byte) rune {
	bi := b - 'A'
	index := s.substitutionIndex[bi]
	substitutionList := s.substitutions[bi]
	substitutionListSize := uint16(len(substitutionList))
	result := substitutionList[index]

	if substitutionListSize > 1 {
		index++
		if index >= substitutionListSize {
			index = 0
		}
		s.substitutionIndex[bi] = index
	}

	return result
}
