package lib

import "fmt"

type SyncError struct {
	errors []error
}

func (e *SyncError) Error() string {
	var s string
	for _, e := range e.errors {
		s += fmt.Sprintf("%v\n", e.Error())
	}
	return s
}

type MissingDirectoryError struct {
	dir string
}

func (e *MissingDirectoryError) Error() string {
	return fmt.Sprintf("directory '%s' is missing", e.dir)
}
