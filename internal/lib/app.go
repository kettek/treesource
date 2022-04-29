package lib

import (
	"context"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// App struct
type App struct {
	ctx     context.Context
	views   []*DirView
	Project *Project
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Context() context.Context {
	return a.ctx
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// UnsavedError represents an error reporting if a project is unsaved.
type UnsavedError struct {
}

// Error returns error.
func (e *UnsavedError) Error() string {
	return "file is unsaved"
}

// NoProjectError represents an error when no project is loaded but a call is made interact with a project.
type NoProjectError struct {
}

// Error returns error.
func (e *NoProjectError) Error() string {
	return "no project is loaded"
}

// HasProject returns if a project is loaded.
func (a *App) HasProject() bool {
	return a.Project != nil
}

// NewProject creates a new treesource project file at the given path and adds the passed directory as its first directory.
func (a *App) NewProject(name string, dir string, ignoreDot bool) error {
	err := a.CloseProjectFile(false)
	if err != nil {
		if _, ok := err.(*NoProjectError); !ok {
			return err
		}
	}

	p := &Project{
		Title: strings.TrimSuffix(path.Base(name), path.Ext(name)),
		Path:  name,
	}
	if dir != "" {
		if err := p.AddDirectory(dir, ignoreDot); err != nil {
			return err
		}
	}

	a.Project = p

	// And save it.
	err = a.SaveProject(true)
	if err != nil {
		return err
	}

	return nil
}

// SaveProject saves the current project.
func (a *App) SaveProject(force bool) error {
	if a.Project == nil {
		return &NoProjectError{}
	}
	if a.Project.Changed() || force {
		b, err := yaml.Marshal(a.Project)
		if err != nil {
			return err
		}
		err = os.WriteFile(a.Project.Path, b, 0755)
		if err != nil {
			return err
		}
		a.Project.changed = false
	}

	return nil
}

// LoadFile loads a treesource project file. If the project is unsaved and force is not true, then an UnsavedError is returned.
func (a *App) LoadProjectFile(name string, force bool) error {
	err := a.CloseProjectFile(force)
	if err != nil {
		if _, ok := err.(*NoProjectError); !ok {
			return err
		}
	}

	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &a.Project)
	a.Project.Path = name

	return err
}

// CloseProjectFile closes the current project if one exists. If the project is unsaved and force is not true, then an UnsavedError is returned. If no project is open, then NoProjectError is returned.
func (a *App) CloseProjectFile(force bool) error {
	if a.Project == nil {
		return &NoProjectError{}
	}
	if a.Project.Changed() && !force {
		return &UnsavedError{}
	}
	a.Project = nil

	return nil
}

type DirEntry struct {
	Name  string   `json:"Name"`
	Info  FileInfo `json:"Info"`
	Error error    `json:"Error"`
}

type FileInfo struct {
	Size    int64       `json:"Size"`
	IsDir   bool        `json:"IsDir"`
	Mode    fs.FileMode `json:"FileMode"`
	ModTime time.Time   `json:"ModTime"`
}

// ListFiles
func (a *App) ListFiles(loc string) (entries []DirEntry, err error) {
	fsentries, err := os.ReadDir(loc)
	if err != nil {
		return nil, err
	}
	for _, e := range fsentries {
		de := DirEntry{
			Name: e.Name(),
		}
		info, err := e.Info()
		if err != nil {
			de.Error = err
		} else {
			de.Info = FileInfo{
				Size:    info.Size(),
				IsDir:   info.IsDir(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
			}
		}

		entries = append(entries, de)
	}
	return
}

func (a *App) NewDirView() uuid.UUID {
	dv := &DirView{
		UUID: uuid.New(),
	}
	a.views = append(a.views, dv)
	return dv.UUID
}

func (a *App) RemoveDirView(u uuid.UUID) {
	for i, dv := range a.views {
		match := true
		for i, v := range u {
			if dv.UUID[i] != v {
				match = false
			}
		}
		if !match {
			continue
		}
		a.views = append(a.views[:i], a.views[i+1:]...)
		break
	}
}
