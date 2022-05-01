package lib

import "github.com/google/uuid"

type Event interface{}

type DirectorySyncEvent struct {
	UUID uuid.UUID
	Name string
}

type DirectorySyncedEvent struct {
	UUID  uuid.UUID
	Name  string
	Error error
}

type DirectoryFileAddEvent struct {
	Name string
	File *DirectoryEntry
}

type DirectoryFileMissingEvent struct {
	Name string
	File *DirectoryEntry
}

type DirectoryFileFoundEvent struct {
	Name string
	File *DirectoryEntry
}
