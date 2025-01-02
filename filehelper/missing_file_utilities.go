package filehelper

import (
	"fmt"
	"os"
)

// CloseFile closes a file and reports an error, if it occurs.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error closing file '%s': %v\n", file.Name(), err)
	}
}
