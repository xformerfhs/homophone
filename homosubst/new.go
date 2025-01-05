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
// Version: 1.1.1
//
// Change history:
//    2024-09-17: V1.0.0: Created.
//    2025-01-03: V1.1.0: Use "randomlist".
//    2025-01-03: V1.1.1: Refactor "getSubstitutionLengths". Fix randomAdjustment.
//

package homosubst

import (
	"bufio"
	"errors"
	"fmt"
	"homophone/filehelper"
	"homophone/randomlist"
	"homophone/slicehelper"
	"math"
	"math/rand/v2"
	"os"
	"unicode"
)

// ******** Private constants ********

var substitutionAlphabet = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
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
func getFrequenciesFromFile(fileName string) ([]int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseWithName(file)

	frequencies := make([]int, sourceAlphabetSize)
	totalCount := 0

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
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

	if scanner.Err() != nil {
		return nil, 0, scanner.Err()
	}

	return frequencies, totalCount, nil
}

// getSubstitutionLengths calculates the number of substitutions for each character
// from the frequencies.
func getSubstitutionLengths(sourceFrequencies []int, totalCount int, substitutionAlphabetSize uint16) ([]uint16, error) {
	result := make([]uint16, sourceAlphabetSize)
	countPerSource := float64(totalCount) / float64(substitutionAlphabetSize)
	stepSize := countPerSource * 0.5
	stepThreshold := 0.5 / float64(totalCount)

	// Vary the count per source until the sum of the substitution lengths
	// matches the number of characters in the substitution alphabet.
	for {
		// 1. Calculate the substitution lengths from the frequencies.
		resultCount := calculateSubstitutionLengths(sourceFrequencies, countPerSource, result)

		// 2. If the number is not correct, make a bisection.
		if resultCount > substitutionAlphabetSize {
			countPerSource += stepSize
		} else {
			if resultCount < substitutionAlphabetSize {
				countPerSource -= stepSize
			} else {
				// We got the right number. Done.
				break
			}
		}

		// 3. Halve the step size.
		stepSize *= 0.5
		if stepSize < stepThreshold {
			// There is no converging solution, but we are near.
			// Make random adjustments to match [substitutionAlphabetSize].
			if randomAdjustment(result, substitutionAlphabetSize, resultCount) == 0 {
				break
			} else {
				// Bail out, if it is still not possible to match [substitutionAlphabetSize].
				return nil, errors.New(`unable to find a distribution`)
			}
		}
	}

	return result, nil
}

// calculateSubstitutionLengths calculates the substitution lengths with the current [countPerSource]
func calculateSubstitutionLengths(sourceFrequencies []int, countPerSource float64, substitutionLengths []uint16) uint16 {
	substitutionCount := uint16(0)
	for i, f := range sourceFrequencies {
		// Only calculate a substitution length if the value occurs, at all.
		if f != 0 {
			substitutionLength := hillHuntingtonRound(f, countPerSource)
			substitutionLengths[i] = substitutionLength
			substitutionCount += substitutionLength
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

// hillHuntingRound rounds the quotient of count/countPerSource according to the
// Hill/Huntington method (https://en.wikipedia.org/wiki/Huntington%E2%80%93Hill_method).
func hillHuntingtonRound(count int, countPerSource float64) uint16 {
	proportion := float64(count) / countPerSource
	minProportion := math.Floor(proportion)
	maxProportion := minProportion + 1
	roundingBoundary := math.Sqrt(minProportion * maxProportion)
	if proportion < roundingBoundary {
		return uint16(minProportion)
	} else {
		return uint16(maxProportion)
	}
}

// randomAdjustment makes random adjustments to the substitution lengths until the total count matches the wanted count.
func randomAdjustment(substitutionCount []uint16, wantedCount uint16, currentCount uint16) int16 {
	diffCount := int16(currentCount) - int16(wantedCount)
	substitutionLength := len(substitutionCount)
	eligibleIndices := make([]int, 0)
	for i := 0; i < substitutionLength; i++ {
		if substitutionCount[i] > 1 {
			eligibleIndices = append(eligibleIndices, i)
		}
	}

	eligibleLength := len(eligibleIndices)
	for diffCount != 0 && eligibleLength != 0 {
		i := rand.IntN(eligibleLength)
		si := eligibleIndices[i]
		c := substitutionCount[si]
		if diffCount > 0 {
			c--
			diffCount--
		} else {
			c++
			diffCount++
		}

		substitutionCount[si] = c

		if diffCount != 0 {
			eligibleIndices = slicehelper.RemoveNoOrder(eligibleIndices, i)
			eligibleLength--
		}
	}

	return diffCount
}
