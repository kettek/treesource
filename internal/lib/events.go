package lib

import "github.com/google/uuid"

type Event interface{}

const EventDirectories string = "directories"

type DirectoriesEvent struct {
	Directories []DirectoryEvent
}

const EventDirectory string = "directory"

type DirectoryEvent struct {
	UUID uuid.UUID
	Name string
}

const EventDirectoryAdd string = "directory-add"

type DirectoryAddEvent struct {
	UUID uuid.UUID
	Name string
}

const EventDirectorySync string = "directory-sync"

type DirectorySyncEvent struct {
	UUID uuid.UUID
	Name string
}

const EventDirectorySynced string = "directory-synced"

type DirectorySyncedEvent struct {
	UUID  uuid.UUID
	Name  string
	Error error
}

const EventDirectoryEntry string = "directory-entry"

type DirectoryEntryEvent struct {
	Name  string
	Entry *DirectoryEntry
}

const EventDirectoryEntryAdd string = "directory-entry-add"

type DirectoryEntryAddEvent struct {
	Name  string
	Entry *DirectoryEntry
}

const EventDirectoryEntryMissing string = "directory-entry-missing"

type DirectoryEntryMissingEvent struct {
	Name  string
	Entry *DirectoryEntry
}

const EventDirectoryEntryFound string = "directory-entry-found"

type DirectoryEntryFoundEvent struct {
	Name  string
	Entry *DirectoryEntry
}
