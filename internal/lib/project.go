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

	err := d.SyncEntries()
	if err != nil {
		return err
	}

	d.On("sync", p.SyncDirectoryCallback)
	d.On("synced", p.SyncedDirectoryCallback)
	d.On("add", p.FileAddCallback)
	d.On("found", p.FileFoundCallback)
	d.On("missing", p.FileMissingCallback)

	p.Directories = append(p.Directories, d)

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
	p.Emit("directory-sync", e)
}

func (p *Project) SyncedDirectoryCallback(e Event) {
	p.Emit("directory-synced", e)
}

func (p *Project) FileAddCallback(e Event) {
	p.Emit("directory-file-add", e)
}

func (p *Project) FileMissingCallback(e Event) {
	p.Emit("directory-file-missing", e)
}

func (p *Project) FileFoundCallback(e Event) {
	p.Emit("directory-file-found", e)
}
