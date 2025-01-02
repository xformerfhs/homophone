package homosubst

import (
	"bufio"
	"homophone/filehelper"
	"homophone/oshelper"
	"os"
	"path"
	"strings"
)

// Decrypt decrypts given file with the loaded homophone substitution.
func (s *Substitutor) Decrypt(inFileName string) (string, error) {
	inFile, err := os.Open(inFileName)
	if err != nil {
		return ``, makeFileError(`open`, `in`, inFileName, err)
	}
	defer filehelper.CloseFile(inFile)

	outFileName := buildDecryptOutFileName(inFileName)
	var outFile *os.File
	outFile, err = os.OpenFile(outFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return ``, makeFileError(`open`, `out`, outFileName, err)
	}
	defer filehelper.CloseFile(outFile)

	decryptionMap := buildDecryptionMap(s.substitutions)

	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range text {
			decrypted, found := decryptionMap[r]
			if !found {
				decrypted = r
			}

			_, err = writer.WriteRune(decrypted)
			if err != nil {
				return ``, makeFileError(`write to`, `out`, outFileName, err)
			}
		}

		_, err = writer.WriteString(oshelper.NewLine)
		if err != nil {
			return ``, makeFileError(`write to`, `out`, outFileName, err)
		}
	}

	if scanner.Err() != nil {
		return ``, scanner.Err()
	}

	err = writer.Flush()
	if err != nil {
		return ``, makeFileError(`flush`, `out`, outFileName, err)
	}

	return outFileName, nil
}

func buildDecryptionMap(substitutions [][]rune) map[rune]rune {
	result := make(map[rune]rune)
	destinationRune := 'A'
	for _, list := range substitutions {
		for _, substitution := range list {
			result[substitution] = destinationRune
		}
		destinationRune++
	}
	return result
}

// buildDecryptOutFileName builds the file name of the output file.
func buildDecryptOutFileName(fileName string) string {
	pos := strings.LastIndex(fileName, `_homophone`)
	if pos >= 0 {
		return fileName[:pos] + `_decrypted` + fileName[pos+10:]
	} else {
		dir := path.Dir(fileName)
		base := path.Base(fileName)
		ext := path.Ext(fileName)
		base = strings.TrimSuffix(base, ext)
		return path.Join(dir, base+"_decrypted"+ext)
	}
}
