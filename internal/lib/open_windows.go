package lib

import (
	"os/exec"
	"path/filepath"
)

func (a *App) OpenFile(paths []string) error {
	// This seems bad.
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", filepath.Join(paths...))
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
