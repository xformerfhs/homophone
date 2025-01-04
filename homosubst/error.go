package homosubst

import "fmt"

// makeFileError builds an error for a file error.
func makeFileError(operation string, direction string, fileName string, err error) error {
	return fmt.Errorf(`could not %s %sput file '%s': %w`, operation, direction, fileName, err)
}
