package main

import (
	"homophone/filehelper"
	"path/filepath"
	"strings"
)

// ******** Public functions ********

// buildDecryptOutFilePath builds the file path of the decrypted output file.
func buildDecryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, `decrypted`)
}

// buildEncryptOutFilePath builds the file path of the encrypted output file.
func buildEncryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, `homophone`)
}

// buildSubstFilePath builds the file path of the substitution file.
func buildSubstFilePath(filePath string) string {
	dir, base, ext := filehelper.PathComponents(filePath)
	if len(ext) > 0 {
		base = strings.TrimSuffix(base, ext) + `_` + ext[1:]
	}

	return filepath.Join(dir, base+`.subst`)
}

// ******** Private functions ********

// buildFilePathWithMarker builds a file path with a marker that is separated by '_' after the base name.
func buildFilePathWithMarker(filePath string, marker string) string {
	dir, base, ext := filehelper.PathComponents(filePath)
	return filepath.Join(dir, base+`_`+marker+ext)
}
