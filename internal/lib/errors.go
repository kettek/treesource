package lib

import (
	"fmt"

	"github.com/google/uuid"
)

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
	dir  string
	uuid uuid.UUID
}

func (e *MissingDirectoryError) Error() string {
	return fmt.Sprintf("directory '%s' is missing", e.dir)
}

type MissingDirectoryViewError struct {
	uuid uuid.UUID
}

func (e *MissingDirectoryViewError) Error() string {
	return fmt.Sprintf("directory '%s' for view is missing", e.uuid)
}

type MissingTagsViewError struct {
	uuid uuid.UUID
}

func (e *MissingTagsViewError) Error() string {
	return fmt.Sprintf("tags view '%s' is missing", e.uuid)
}

type MissingSessionError struct {
}

func (e *MissingSessionError) Error() string {
	return fmt.Sprintf("missing session")
}
