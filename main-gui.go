//go:build !tui

package main

import (
	"embed"
	"treesource/internal/lib"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := lib.NewApp()

	dir := lib.DirView{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "treesource",
		Width:  1280,
		Height: 720,
		Assets: assets,
		Bind: []interface{}{
			app,
			&dir,
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
