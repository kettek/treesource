package lib

import (
	"fmt"

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
		UUID:       uuid.New(),
		Path:       name,
		IgnoreDot:  ignoreDot,
		SyncOnLoad: true,
	}

	err := d.SyncEntries()
	if err != nil {
		return err
	}

	p.Directories = append(p.Directories, d)

	return nil
}
