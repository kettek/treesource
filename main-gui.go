//go:build !tui

package main

import (
	"embed"
	"fmt"
	"treesource/internal/lib"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/dist
var assets embed.FS

var app *WApp

func main() {
	// Create an instance of the app structure
	app = &WApp{
		App:     *lib.NewApp(),
		started: false,
	}

	if err := lib.EnsureSession("default"); err != nil {
		panic(err)
	}

	session, err := lib.LoadSession("default")
	if err != nil {
		panic(err)
	}
	app.Session = session

	if err := app.SetupSession(); err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "treesource",
		Width:  1280,
		Height: 720,
		Assets: assets,
		Bind: []interface{}{
			app,
			&lib.Project{},
			&lib.Directory{},
			&lib.DirectoryEntry{},
			&Dialog{},
		},
		OnStartup: app.Startup,
	})

	if err != nil {
		println("Error:", err)
	}
}

// Expose some common stuff

type Dialog struct{}

func (d *Dialog) OpenDirectory(options runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(app.Context(), options)
}

func (d *Dialog) OpenFile(options runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(app.Context(), options)
}

func (d *Dialog) SaveFile(options runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(app.Context(), options)
}

func (d *Dialog) Message(options runtime.MessageDialogOptions) (string, error) {
	return runtime.MessageDialog(app.Context(), options)
}

type WApp struct {
	lib.App
	started bool
}

func (w *WApp) Ready() {
	if !w.started && w.Session.Project != "" {
		w.LoadProjectFile(w.Session.Project, true)
		w.started = true
	}
	if w.Project == nil {
		return
	}
	runtime.EventsEmit(w.Context(), "project-load", w.Project)
	w.RefreshTitle()
	// Also send the actual directory contents.

	for _, d := range w.Project.Directories {
		w.Project.Emit(lib.EventDirectory, lib.DirectoryAddEvent{
			UUID: d.UUID,
			Path: d.Path,
		})
	}
	for _, d := range w.Project.Directories {
		d.EmitAllEntries()
	}
}

func (w *WApp) SetupSession() error {
	if w.Session == nil {
		return &lib.MissingSessionError{}
	}
	// Setup session event handling.
	app.Session.On(lib.EventViewDirectoryAdd, func(e lib.Event) {
		runtime.EventsEmit(app.Context(), lib.EventViewDirectoryAdd, e)
	})
	app.Session.On(lib.EventViewDirectoryRemove, func(e lib.Event) {
		runtime.EventsEmit(app.Context(), lib.EventViewDirectoryRemove, e)
	})
	app.Session.On(lib.EventViewTagsAdd, func(e lib.Event) {
		runtime.EventsEmit(app.Context(), lib.EventViewTagsAdd, e)
	})
	app.Session.On(lib.EventViewTagsRemove, func(e lib.Event) {
		runtime.EventsEmit(app.Context(), lib.EventViewTagsRemove, e)
	})

	return nil
}

func (w *WApp) NewProject(name string, dir string, ignoreDot bool) error {
	err := w.App.NewProject(name, dir, ignoreDot)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-load", w.Project)
		w.RefreshTitle()
		w.InitProject()
		w.Session.Project = w.Project.Path
		w.Session.PendingSave()
	}
	return err
}

func (w *WApp) LoadProjectFile(name string, force bool) error {
	err := w.App.LoadProjectFile(name, force)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-load", w.Project)
		w.RefreshTitle()
		if w.Project.Changed() {
			runtime.EventsEmit(w.Context(), "project-changed")
			w.Project.Unchange()
		}
		w.InitProject()
		w.Session.Project = w.Project.Path
		w.Session.PendingSave()
	}
	return err
}

func (w *WApp) InitProject() error {
	w.Project.On("project-change", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventProjectChange, e)
	})
	w.Project.On("directory", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectory, e)
	})
	w.Project.On("directory-add", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryAdd, e)
	})
	w.Project.On("directory-remove", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryRemove, e)
	})
	w.Project.On("directory-sync", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectorySync, e)
	})
	w.Project.On("directory-synced", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectorySynced, e)
	})
	w.Project.On("directory-entry", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryEntry, e)
	})
	w.Project.On("directory-entry-add", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryEntryAdd, e)
	})
	w.Project.On("directory-entry-missing", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryEntryMissing, e)
	})
	w.Project.On("directory-entry-found", func(e lib.Event) {
		runtime.EventsEmit(w.Context(), lib.EventDirectoryEntryFound, e)
	})

	w.App.InitProject()

	return nil
}

func (w *WApp) CloseProjectFile(force bool) error {
	err := w.App.CloseProjectFile(force)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-unload", nil)
		w.RefreshTitle()
		w.Session.Project = ""
		w.Session.PendingSave()
	}
	return err
}

func (w *WApp) AddProjectDirectory(dir string, ignoreDot bool) error {
	return w.App.AddProjectDirectory(dir, ignoreDot)
}

func (w *WApp) GetProject() *lib.Project {
	return w.Project
}

func (w *WApp) RefreshTitle() {
	title := "treesource"
	if w.Project != nil {
		if w.Unsaved() {
			title = fmt.Sprintf("*%s - %s", w.Project.Title, title)
		} else {
			title = fmt.Sprintf("%s - %s", w.Project.Title, title)
		}
	}
	runtime.WindowSetTitle(w.Context(), title)
}
