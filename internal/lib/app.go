package lib

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"math"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"treesource/internal/do"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
	"gopkg.in/yaml.v3"

	"image/png"

	// image decoders
	_ "image/gif"
	_ "image/jpeg"

	// extended image decoders
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// App struct
type App struct {
	ctx     context.Context
	Project *Project
	Session *Session
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
	p.history = do.History[*Project]{
		Target: p,
	}
	if dir != "" {
		if err := p.AddDirectory(dir, ignoreDot); err != nil {
			return err
		}
	}

	// Reset the history so the user cannot undo the initial directory.
	p.history.Reset()

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

// RemoveProjectDirectory removes a directory by its UUID.
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
	if a.Unsaved() || force {
		b, err := yaml.Marshal(a.Project)
		if err != nil {
			return err
		}
		err = os.WriteFile(a.Project.Path, b, 0755)
		if err != nil {
			return err
		}
		a.Project.history.SavedPos = a.Project.history.Pos
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
	a.Project.history = do.History[*Project]{
		Target: a.Project,
	}

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

func (a *App) Undo() {
	if a.Project == nil {
		return
	}
	a.Project.history.Undo()
}

func (a *App) Redo() {
	if a.Project == nil {
		return
	}
	a.Project.history.Redo()
}

func (a *App) Undoable() bool {
	if a.Project == nil {
		return false
	}
	return a.Project.history.Undoable()
}

func (a *App) Redoable() bool {
	if a.Project == nil {
		return false
	}
	return a.Project.history.Redoable()
}

func (a *App) Unsaved() bool {
	if a.Project == nil {
		return false
	}
	return a.Project.history.SavedPos != a.Project.history.Pos
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

// QueryFile queries a given file, returning stats for it if it exists.
func (a *App) QueryFile(root string, path string) (FileInfo, error) {
	p := filepath.Join(root, path)
	info, err := os.Stat(p)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Name:        info.Name(),
		Path:        p,
		Size:        info.Size(),
		Mode:        uint32(info.Mode()),
		Permissions: info.Mode().Perm().String(),
		Type:        info.Mode().Type().String(),
		Special:     !info.Mode().Perm().IsRegular(),
		ModTime:     info.ModTime(),
		Mimetype:    mime.TypeByExtension(filepath.Ext(path)),
	}, err
}

func (a *App) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (a *App) GenerateIcon(paths []string, opts IconOptions) (Icon, error) {
	path := filepath.Join(paths...)
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return Icon{}, err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return Icon{}, err
	}

	var w, h int

	if img.Bounds().Dx() <= opts.MaxWidth && img.Bounds().Dy() <= opts.MaxHeight {
		w = img.Bounds().Dx()
		h = img.Bounds().Dy()
	} else {
		ratio := math.Min(float64(opts.MaxWidth)/float64(img.Bounds().Dx()), float64(opts.MaxHeight)/float64(img.Bounds().Dy()))
		w = int(math.Round(float64(img.Bounds().Dx()) * ratio))
		h = int(math.Round(float64(img.Bounds().Dy()) * ratio))
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if opts.Method == "CatmullRom" {
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	} else if opts.Method == "NearestNeighbor" {
		draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		//} else if opts.Method == "ApproxBiLinear" {
	} else {
		draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	}

	var b bytes.Buffer
	o := bufio.NewWriter(&b)
	if err := png.Encode(o, dst); err != nil {
		return Icon{}, err
	}
	return Icon{
		Bytes:  b.Bytes(),
		Format: format,
	}, nil
}
