package lib

import (
	"context"
	"os"
	"path"
	"strings"

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
		Title:   strings.TrimSuffix(path.Base(name), path.Ext(name)),
		Path:    name,
		Emitter: *NewEmitter(),
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

// AddProjectDirectory adds the given directory to the project.
func (a *App) AddProjectDirectory(dir string, ignoreDot bool) error {
	if a.Project == nil {
		return &NoProjectError{}
	}
	if err := a.Project.AddDirectory(dir, ignoreDot); err != nil {
		return err
	}
	return nil
}

func (a *App) RemoveProjectDirectory(uuid uuid.UUID) error {
	if a.Project == nil {
		return &NoProjectError{}
	}
	if err := a.Project.RemoveDirectoryByUUID(uuid); err != nil {
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
	a.Project.Emitter = *NewEmitter()
	a.Project.Path = name

	return err
}

func (a *App) InitProject() error {
	for _, d := range a.Project.Directories {
		d.Emitter = *NewEmitter()
		a.Project.Emit(EventDirectoryAdd, DirectoryAddEvent{
			UUID: d.UUID,
			Path: d.Path,
		})
		if err := a.Project.InitDirectory(&d); err != nil {
			panic(err)
		}
		d.EmitAllEntries()
	}
	return nil
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
