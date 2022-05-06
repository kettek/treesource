package lib

import "github.com/google/uuid"

type Event interface{}

const EventProjectChange string = "project-change"

type ProjectChangeEvent struct {
}

const EventDirectories string = "directories"

type DirectoriesEvent struct {
	Directories []DirectoryEvent
}

const EventDirectory string = "directory"

type DirectoryEvent struct {
	UUID uuid.UUID
}

const EventDirectoryAdd string = "directory-add"

type DirectoryAddEvent struct {
	UUID uuid.UUID
	Path string
}

const EventDirectoryRemove string = "directory-remove"

type DirectoryRemoveEvent struct {
	UUID uuid.UUID
}

const EventDirectorySync string = "directory-sync"

type DirectorySyncEvent struct {
	UUID uuid.UUID
}

const EventDirectorySynced string = "directory-synced"

type DirectorySyncedEvent struct {
	UUID  uuid.UUID
	Error error
}

const EventDirectoryEntry string = "directory-entry"

type DirectoryEntryEvent struct {
	UUID  uuid.UUID
	Entry *DirectoryEntry
}

const EventDirectoryEntryAdd string = "directory-entry-add"

type DirectoryEntryAddEvent struct {
	UUID  uuid.UUID
	Entry *DirectoryEntry
}

const EventDirectoryEntryMissing string = "directory-entry-missing"

type DirectoryEntryMissingEvent struct {
	UUID  uuid.UUID
	Entry *DirectoryEntry
}

const EventDirectoryEntryFound string = "directory-entry-found"

type DirectoryEntryFoundEvent struct {
	UUID  uuid.UUID
	Entry *DirectoryEntry
}

/*
Session -> View events
*/
const EventViewDirectoryAdd string = "view-directory-add"

type ViewDirectoryAddEvent struct {
	View *DirectoryView
}

const EventViewDirectoryRemove string = "view-directory-remove"

type ViewDirectoryRemoveEvent struct {
	View *DirectoryView
}

const EventViewDirectoryNavigate string = "view-directory-navigate"

type ViewDirectoryNavigateEvent struct {
	UUID uuid.UUID
	Path string
}

const EventViewTagsAdd string = "view-tags-add"

type ViewTagsAddEvent struct {
	View *TagsView
}

const EventViewTagsRemove string = "view-tags-remove"

type ViewTagsRemoveEvent struct {
	View *TagsView
}

const EventViewSelect string = "view-select"

type ViewSelectEvent struct {
	UUID uuid.UUID
}
