package homosubst

// ======== Public types =========

// Substitutor contains the data needed for a homophone substitution cipher.
type Substitutor struct {
	fileName                 string
	substitutions            [][]rune
	substitutionIndex        []uint16
	substitutionAlphabetSize uint16
}
