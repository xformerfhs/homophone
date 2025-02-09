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
//    2025-01-03: V1.1.0: Use "randomlist".
//    2025-01-03: V1.1.1: Refactor "getSubstitutionLengths". Fix randomAdjustment.
//    2025-01-05: V1.2.0: Correct substitution alphabet.
//    2025-02-08: V2.0.0: Use rune scanner, make substitution length calculation faster.
//

package homosubst

import (
	"bufio"
	"fmt"
	"homophone/distributor"
	"homophone/filehelper"
	"homophone/randomlist"
	"math/rand/v2"
	"os"
	"unicode"
)

// ******** Private constants ********

// substitutionAlphabet is the substitution alphabet.
// There must be no other characters than letters.
var substitutionAlphabet = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

// requiredSubstitutionAlphabetSize is the expected substitution alphabet size for substitution files.
var requiredSubstitutionAlphabetSize = uint32(len(substitutionAlphabet))

// sourceAlphabetSize contains the size of the alphabet to map, i.e. A-Z.
const sourceAlphabetSize uint16 = 26

// ******** Public creation functions ********

// NewSubstitutor creates a new substitutor for the given file.
func NewSubstitutor(sourceFileName string) (*Substitutor, error) {
	substitutionRunes := []rune(substitutionAlphabet)
	substitutionAlphabetSize := uint16(len(substitutionRunes))

	result := &Substitutor{}

	result.substitutionAlphabetSize = substitutionAlphabetSize

	// 1. Get the character frequencies from the file.
	sourceFrequencies, totalCount, err := getFrequenciesFromFile(sourceFileName)
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return nil, fmt.Errorf(`source file '%s' has no characters in the range A-Z`, sourceFileName)
	}

	// 2. Get the lengths of the substitutions of each character from the frequencies.
	var substitutionLengths []uint16
	substitutionLengths, err = getSubstitutionLengths(sourceFrequencies, totalCount, substitutionAlphabetSize)
	if err != nil {
		return nil, err
	}

	// 3. Build the substitution lists from the lengths.
	result.substitutions = generateSubstitutions(substitutionLengths, substitutionRunes, substitutionAlphabetSize)

	return result, nil
}

// ******** Private functions ********

// getFrequenciesFromFile calculates the frequencies of each character in the file.
func getFrequenciesFromFile(fileName string) ([]uint, uint, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseWithName(file)

	frequencies := make([]uint, sourceAlphabetSize)
	totalCount := uint(0)

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range text {
			r = unicode.ToUpper(r)
			if r >= 'A' && r <= 'Z' {
				frequencies[r-'A']++
				totalCount++
			}
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, 0, scanErr
	}

	return frequencies, totalCount, nil
}

// getSubstitutionLengths calculates the number of substitutions for each character
// from the frequencies.
func getSubstitutionLengths(sourceFrequencies []uint, totalCount uint, substitutionAlphabetSize uint16) ([]uint16, error) {
	result := make([]uint16, sourceAlphabetSize)

	calculateSubstitutionLengths(sourceFrequencies, totalCount, substitutionAlphabetSize, result)

	return result, nil
}

// calculateSubstitutionLengths calculates the substitution lengths.
// The result may not have the correct substitution count. This will be the case if there are
// a lot of characters with the same count.
func calculateSubstitutionLengths(
	sourceFrequencies []uint,
	totalCount uint,
	substitutionAlphabetSize uint16,
	substitutionLengths []uint16) {
	substitutionCount := initializeSubstitutionLengths(sourceFrequencies, substitutionLengths)

	// 2. Distribute the remaining substitution alphabet size among the characters.
	remainingCount := substitutionAlphabetSize - substitutionCount
	additionalLengths := distributor.SainteLagueDistribution(
		sourceFrequencies,
		totalCount,
		uint(remainingCount))

	// 3. Add the distributed lengths to the count of 1 that has already been set.
	for i := range substitutionLengths {
		substitutionLengths[i] += uint16(additionalLengths[i])
	}

	substitutionCount += remainingCount
}

// initializeSubstitutionLengths initializes the substitution lengths to have the value 1 for each
// character that appears in the source.
func initializeSubstitutionLengths(sourceFrequencies []uint, substitutionLengths []uint16) uint16 {
	substitutionCount := uint16(0)
	// 1. Each source character gets one substitution character.
	for i, f := range sourceFrequencies {
		// Only calculate a substitution length if the value occurs, at all.
		if f != 0 {
			substitutionLengths[i] = 1
			substitutionCount++
		}
	}
	return substitutionCount
}

// generateSubstitutions Generate the substitution characters from the lengths per character.
func generateSubstitutions(
	substitutionLengths []uint16,
	substitutionAlphabet []rune,
	substitutionAlphabetSize uint16) []*randomlist.RandomList[rune] {
	used := make([]bool, substitutionAlphabetSize)
	result := make([]*randomlist.RandomList[rune], sourceAlphabetSize)
	for i, substitutionLength := range substitutionLengths {
		list := make([]rune, substitutionLength)
		for j := range substitutionLength {
			list[j] = substitutionAlphabet[getSubstitutionAlphabetIndex(used, substitutionAlphabetSize)]
		}
		result[i] = randomlist.New(list)
	}

	return result
}

// getSubstitutionAlphabetIndex gets the substitution index into the substitution alphabet.
func getSubstitutionAlphabetIndex(used []bool, usedSize uint16) int {
	for {
		i := rand.IntN(int(usedSize))

		if !used[i] {
			used[i] = true
			return i
		}
	}
}
