package homosubst

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"homophone/compressedinteger"
	"homophone/integritycheckedfile"
	"homophone/keygenerator"
	"path/filepath"
	"strings"
	"unicode"
)

var fileMagic = []byte(`HFDF`)

const actVersion byte = 0

var generator = []byte{
	0xfe, 0xb9, 0x66, 0x43,
	0x18, 0x5b, 0x51, 0xdf,
	0x86, 0x99, 0xe5, 0x09,
	0xa4, 0xdc, 0x0d, 0xad,
	0x82, 0xed, 0xc4, 0x30,
}

var integrityKey = []byte{
	0x74, 0xbc, 0x06, 0x3e,
	0x56, 0x17, 0xda, 0xd4,
	0xf2, 0xc7, 0x91, 0x37,
	0x2a, 0xe2, 0xbf, 0x32,
}

var additionalData = []byte(`HoTzpLoZ`)

func (s *Substitutor) Save() (string, error) {
	substFileName := makeSubstFileName(s.fileName)
	w, err := integritycheckedfile.NewWriter(
		substFileName,
		sha3.New256,
		keygenerator.GenerateKey(generator, integrityKey),
		additionalData)
	if err != nil {
		return ``, err
	}
	defer w.Close()

	_, err = w.Write(fileMagic)
	if err != nil {
		return ``, err
	}

	_, err = w.Write([]byte{actVersion})
	if err != nil {
		return ``, err
	}

	var size []byte

	size, err = compressedinteger.FromUInt32(uint32(s.substitutionAlphabetSize))
	if err != nil {
		return ``, err
	}
	_, err = w.Write(size)
	if err != nil {
		return ``, err
	}

	for _, l := range s.substitutions {
		size, err = compressedinteger.FromUInt32(uint32(len(l)))
		if err != nil {
			return ``, err
		}

		_, err = w.Write(size)
		if err != nil {
			return ``, err
		}

		for _, r := range l {
			size, err = compressedinteger.FromUInt32(uint32(r))
			if err != nil {
				return ``, err
			}

			_, err = w.Write(size)
			if err != nil {
				return ``, err
			}
		}
	}

	return substFileName, nil
}

func NewLoad(substFileName string) (*Substitutor, error) {
	r, err := integritycheckedfile.NewReader(
		substFileName,
		sha3.New256,
		keygenerator.GenerateKey(generator, integrityKey),
		additionalData)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if r.DataLen() != 136 {
		return nil, errors.New(`wrong file size`)
	}

	buffer := make([]byte, len(fileMagic))
	_, err = r.Read(buffer)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(buffer, fileMagic) {
		return nil, errors.New(`unknown file format`)
	}

	_, err = r.Read(buffer[:1])
	if err != nil {
		return nil, err
	}
	if buffer[0] != actVersion {
		return nil, fmt.Errorf(`unknown version`)
	}

	substitutionData := make([]byte, int(r.DataLen())-len(fileMagic)-1)
	var readBytes int
	readBytes, err = r.Read(substitutionData)
	if err != nil {
		return nil, err
	}

	actPos := 0

	var substitutionAlphabetSize uint32

	substitutionAlphabetSize, readBytes, err = compressedinteger.ToUInt32(substitutionData[actPos:])
	if err != nil {
		return nil, err
	}

	if substitutionAlphabetSize != 52 {
		return nil, fmt.Errorf(`wrong substituion alphabet size: %d`, substitutionAlphabetSize)
	}
	actPos += readBytes

	check := make(map[rune]bool)
	substitutions := make([][]rune, sourceAlphabetSize)
	listCount := 0
	substitutionCount := 0
	for actPos < len(substitutionData) {
		var listSize uint32
		listSize, readBytes, err = compressedinteger.ToUInt32(substitutionData[actPos:])
		if err != nil {
			return nil, err
		}
		actPos += readBytes
		listCount++
		substitutionCount += int(listSize)
		if listCount > int(sourceAlphabetSize) {
			return nil, errors.New(`too many substitution entries`)
		}
		if substitutionCount > int(substitutionAlphabetSize) {
			return nil, errors.New(`too many substitutions`)
		}
		list := make([]rune, listSize)
		var entry uint32
		for i := range listSize {
			entry, readBytes, err = compressedinteger.ToUInt32(substitutionData[actPos:])
			if err != nil {
				return nil, err
			}
			actPos += readBytes
			entryRune := rune(entry)
			if !unicode.IsLetter(entryRune) {
				return nil, errors.New(`invalid substitution entry`)
			}
			if check[entryRune] {
				return nil, fmt.Errorf(`duplicate substitution entry %c`, entryRune)
			}
			list[i] = entryRune
			check[entryRune] = true
		}

		substitutions[listCount-1] = list
	}

	if listCount < int(sourceAlphabetSize) {
		return nil, errors.New(`not enough substitution entries`)
	}

	if substitutionCount < int(substitutionAlphabetSize) {
		return nil, errors.New(`not enough substitutions`)
	}

	return &Substitutor{
		substitutions:            substitutions,
		substitutionAlphabetSize: uint16(substitutionAlphabetSize),
		substitutionIndex:        make([]uint16, sourceAlphabetSize)}, nil
}

func makeSubstFileName(fileName string) string {
	path := filepath.Dir(fileName)
	base := filepath.Base(fileName)
	ext := filepath.Ext(base)
	var realBase string
	if len(ext) > 0 {
		realBase = strings.TrimSuffix(base, ext) + `_` + ext[1:]
	}
	return filepath.Join(path, realBase+`.subst`)
}
