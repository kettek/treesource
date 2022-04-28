package lib

import (
	"context"
	"io/fs"
	"os"
	"time"

	"github.com/google/uuid"
)

// App struct
type App struct {
	ctx   context.Context
	views []*DirView
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
