package homosubst

import (
	"bufio"
	"fmt"
	"homophone/filehelper"
	"homophone/oshelper"
	"os"
	"path"
	"strings"
	"unicode"
)

// Encrypt encrypts the file from the creation call with the built homophone substitution.
func (s *Substitutor) Encrypt(noOther bool) (string, error) {
	inFileName := s.fileName

	inFile, err := os.Open(inFileName)
	if err != nil {
		return ``, makeFileError(`open`, `in`, inFileName, err)
	}
	defer filehelper.CloseFile(inFile)

	outFileName := buildOutFileName(inFileName)
	var outFile *os.File
	outFile, err = os.OpenFile(outFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return ``, makeFileError(`open`, `out`, outFileName, err)
	}
	defer filehelper.CloseFile(outFile)

	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range text {
			r = unicode.ToUpper(r)
			if r >= 'A' && r <= 'Z' {
				_, err = writer.WriteRune(s.SubstituteRune(r))
			} else {
				if !noOther {
					_, err = writer.WriteRune(r)
				}
			}

			if err != nil {
				return ``, makeFileError(`write to`, `out`, outFileName, err)
			}
		}

		if !noOther {
			_, err = writer.WriteString(oshelper.NewLine)
			if err != nil {
				return ``, makeFileError(`write to`, `out`, outFileName, err)
			}
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

// buildOutFileName builds the file name of the output file.
func buildOutFileName(fileName string) string {
	dir := path.Dir(fileName)
	base := path.Base(fileName)
	ext := path.Ext(fileName)
	base = strings.TrimSuffix(base, ext)
	return path.Join(dir, base+"_homophone"+ext)
}

// makeFileError builds an error for a file error.
func makeFileError(operation string, direction string, fileName string, err error) error {
	return fmt.Errorf(`could not %s %sput file '%s': %w`, operation, direction, fileName, err)
}
