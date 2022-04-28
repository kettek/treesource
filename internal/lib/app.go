package lib

import (
	"context"
	"io/fs"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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
