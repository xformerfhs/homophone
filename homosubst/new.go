package homosubst

import (
	"bufio"
	"errors"
	"homophone/filehelper"
	"math"
	"math/rand"
	"os"
	"unicode"
)

// ======== Private constants =========

var substitutionAlphabet = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

// sourceAlphabetSize contains the size of the alphabet to map, i.e. A-Z.
var sourceAlphabetSize uint16 = 26

// ======== Public creation functions =========

// NewSubstitutor creates a new substitutor for the given file.
func NewSubstitutor(sourceFileName string) (*Substitutor, error) {
	runes := []rune(substitutionAlphabet)
	substitutionAlphabetSize := uint16(len(runes))

	result := &Substitutor{}

	result.substitutionAlphabetSize = substitutionAlphabetSize

	// 1. Get the charakter frequencies from the file.
	sourceFrequencies, totalCount, err := getFrequenciesFromFile(sourceFileName)
	if err != nil {
		return nil, err
	}

	// 2. Get the lengths of the substitutions of each character from the frequencies.
	var substitutionLengths []uint16
	substitutionLengths, err = getSubstitutionLengths(sourceFrequencies, totalCount, substitutionAlphabetSize)
	if err != nil {
		return nil, err
	}

	// 3. Build the substitution lists from the lengths.
	result.substitutions = generateSubstitutions(substitutionLengths, runes, substitutionAlphabetSize)

	result.substitutionIndex = make([]uint16, sourceAlphabetSize)

	result.fileName = sourceFileName

	return result, nil
}

// ======== Private functions =========

// getFrequenciesFromFile calculates the frequencies of each character in the file.
func getFrequenciesFromFile(fileName string) ([]int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(file)

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
		resultCount := uint16(0)
		for i, f := range sourceFrequencies {
			count := uint16(math.Max(math.Round(float64(f)/countPerSource), 1.0))
			result[i] = count
			resultCount += count
		}

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

		// 3. Halve the step size and bail out if there were too many tries.
		stepSize *= 0.5
		if stepSize < stepThreshold {
			return nil, errors.New(`unable to find a distribution`)
		}
	}

	return result, nil
}

// generateSubstitutions Generate the substitution characters from the lengths per character.
func generateSubstitutions(substitutionLengths []uint16, substitutionAlphabet []rune, substitutionAlphabetSize uint16) [][]rune {
	used := make([]bool, substitutionAlphabetSize)
	result := make([][]rune, sourceAlphabetSize)
	for i, l := range substitutionLengths {
		list := make([]rune, l)
		for j := range l {
			list[j] = substitutionAlphabet[getSubstitutionAlphabetIndex(used, substitutionAlphabetSize)]
		}
		result[i] = list
	}

	return result
}

// getSubstitutionAlphabetIndex gets the substitution index into the substitution alphabet.
func getSubstitutionAlphabetIndex(used []bool, usedSize uint16) int {
	for {
		i := rand.Intn(int(usedSize))

		if !used[i] {
			used[i] = true
			return i
		}
	}
}
