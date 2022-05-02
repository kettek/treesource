package lib

import (
	"fmt"

	"github.com/google/uuid"
)

// Project represents a full treesource project.
type Project struct {
	Emitter
	Title       string      `json:"Title" yaml:"Title"`             // Title of the project.
	Path        string      `json:"Path" yaml:"Path"`               // Path from which the project file was read and should be saved to.
	Directories []Directory `json:"Directories" yaml:"Directories"` // Directories to pull from as sources.
	changed     bool
}

// Changed represents if the project has unsaved changes.
func (p *Project) Changed() bool {
	return p.changed
}

func (p *Project) Change() {
	p.changed = true
}

func (p *Project) Unchange() {
	p.changed = false
}

type DirectoryExistsError struct {
	dir string
}

func (e *DirectoryExistsError) Error() string {
	return fmt.Sprintf("directory '%s' already exists", e.dir)
}

func (p *Project) AddDirectory(name string, ignoreDot bool) error {
	// Do not add a directory if it already exists.
	for _, d := range p.Directories {
		if d.Path == name {
			return &DirectoryExistsError{name}
		}
	}

	d := Directory{
		UUID:       uuid.New(),
		Path:       name,
		IgnoreDot:  ignoreDot,
		SyncOnLoad: true,
		Emitter:    *NewEmitter(),
	}

	d.On("sync", p.SyncDirectoryCallback)
	d.On("synced", p.SyncedDirectoryCallback)
	d.On("entry", p.EntryCallback)
	d.On("add", p.EntryAddCallback)
	d.On("found", p.EntryFoundCallback)
	d.On("missing", p.EntryMissingCallback)

	err := d.SyncEntries()
	if err != nil {
		return err
	}

	p.Directories = append(p.Directories, d)

	p.Emit(EventDirectoryAdd, DirectoryAddEvent{
		UUID: d.UUID,
		Name: d.Path,
	})

	p.changed = true

	return nil
}

func (p *Project) SyncDirectory(name string) error {
	for _, d := range p.Directories {
		if d.Path == name {
			err := d.SyncEntries()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return &MissingDirectoryError{name}
}

//

func (p *Project) SyncDirectoryCallback(e Event) {
	fmt.Println(EventDirectorySync, e)
	p.Emit(EventDirectorySync, e)
}

func (p *Project) SyncedDirectoryCallback(e Event) {
	fmt.Println(EventDirectorySynced, e)
	p.Emit(EventDirectorySynced, e)
}

func (p *Project) EntryCallback(e Event) {
	fmt.Println(EventDirectoryEntry, e)
	p.Emit(EventDirectoryEntry, e)
}

func (p *Project) EntryAddCallback(e Event) {
	fmt.Println(EventDirectoryEntryAdd, e)
	p.Emit(EventDirectoryEntryAdd, e)
}

func (p *Project) EntryMissingCallback(e Event) {
	fmt.Println(EventDirectoryEntryMissing, e)
	p.Emit(EventDirectoryEntryMissing, e)
}

func (p *Project) EntryFoundCallback(e Event) {
	fmt.Println(EventDirectoryEntryFound, e)
	p.Emit(EventDirectoryEntryFound, e)
}
