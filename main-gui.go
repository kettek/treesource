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
	app = &WApp{*lib.NewApp()}

	// Create application with options
	err := wails.Run(&options.App{
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
}

func (w *WApp) NewProject(name string, dir string, ignoreDot bool) error {
	err := w.App.NewProject(name, dir, ignoreDot)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-load", w.Project)
		runtime.WindowSetTitle(w.Context(), fmt.Sprintf("%s - %s", w.Project.Title, "treesource"))
	}
	return err
}

func (w *WApp) LoadProjectFile(name string, force bool) error {
	err := w.App.LoadProjectFile(name, force)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-load", w.Project)
		runtime.WindowSetTitle(w.Context(), fmt.Sprintf("%s - %s", w.Project.Title, "treesource"))
	}
	return err
}

func (w *WApp) CloseProjectFile(force bool) error {
	err := w.App.CloseProjectFile(force)
	if err == nil {
		runtime.EventsEmit(w.Context(), "project-unload", nil)
		runtime.WindowSetTitle(w.Context(), fmt.Sprintf("%s", "treesource"))
	}
	return err
}

func (w *WApp) GetProject() *lib.Project {
	return w.Project
}
