package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Project represents a full treesource project.
type Project struct {
	Title       string      `json:"Title" yaml:"Title"`             // Title of the project.
	Path        string      `json:"Path" yaml:"Path"`               // Path from which the project file was read and should be saved to.
	Directories []Directory `json:"Directories" yaml:"Directories"` // Directories to pull from as sources.
	changed     bool
}

// Changed represents if the project has unsaved changes.
func (p *Project) Changed() bool {
	return p.changed
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
		UUID:      uuid.New(),
		Path:      name,
		ignoreDot: ignoreDot,
	}

	err := d.SyncEntries()
	if err != nil {
		return err
	}

	p.Directories = append(p.Directories, d)

	return nil
}

// Directory represents a source directory to pull files from.
type Directory struct {
	UUID      uuid.UUID         `json:"UUID" yaml:"UUID"`
	Path      string            `json:"Path" yaml:"Path"`
	Entries   []*DirectoryEntry `json:"Entries" yaml:"Entries"`
	ignoreDot bool
}

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

func (d *Directory) Entry(name string) *DirectoryEntry {
	for _, e := range d.Entries {
		if e.Path == name {
			return e
		}
	}
	return nil
}

func (d *Directory) SyncEntries() error {
	var errors []error
	var walk func(name string)
	walk = func(name string) {
		entries, err := os.ReadDir(name)
		if err != nil {
			errors = append(errors, err)
			return
		}
		for _, e := range entries {
			fullpath := filepath.Join(name, e.Name())
			if d.ignoreDot && e.Name()[0] == '.' {
				continue
			}
			if e.IsDir() {
				walk(fullpath)
			} else {
				if d.Entry(fullpath) == nil {
					d.Entries = append(d.Entries, &DirectoryEntry{
						Path: fullpath,
					})
				}
			}
		}
	}

	walk(d.Path)
	if len(errors) > 0 {
		return &SyncError{errors}
	}
	return nil
}

// DirectoryEntry represents an entry in a treesource directory.
type DirectoryEntry struct {
	// Path is relative to the owning Directory's path.
	Path string   `json:"Path" yaml:"Path"`
	Tags []string `json:"Tags" yaml:"Tags"`
}
