package lib

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Session represents the local session for treesource.
type Session struct {
	Emitter `json:"-" yaml:"-"`
	path    string
	Project string `json:"project" yaml:"project"`
	Views   struct {
		Directories []*DirectoryView
		Tags        []*TagsView
	}
	canceledSave chan struct{}
}

// Refresh causes all pertinent state to emit, so as to resync frontend.
func (s *Session) Refresh() {
	for _, d := range s.Views.Directories {
		s.Emit(EventViewDirectoryAdd, ViewDirectoryAddEvent{
			View: d,
		})
	}
	for _, t := range s.Views.Tags {
		s.Emit(EventViewTagsAdd, ViewTagsAddEvent{
			View: t,
		})
	}
}

// Save saves the session.
func (c *Session) Save() error {
	p, err := GetSessionPath(c.path)
	if err != nil {
		return err
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	// Make basic session file.
	err = os.WriteFile(p, b, 0755)
	return err
}

func (c *Session) PendingSave() {
	// Cancel any current pending saves.
	select {
	case c.canceledSave <- struct{}{}:
	default:
	}
	go func() {
		select {
		case <-time.After(500 * time.Millisecond):
			err := c.Save()
			if err != nil {
				panic(err)
			}
		case <-c.canceledSave:
		}
	}()
}

// LoadSession loads the session file.
func LoadSession(name string) (*Session, error) {
	c := &Session{
		Emitter:      *NewEmitter(),
		canceledSave: make(chan struct{}),
	}

	p, err := GetSessionPath(name)
	if err != nil {
		return nil, err
	}
	c.path = name

	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)

	return c, err
}

// GetSessionPath returns the path for the target session.
func GetSessionPath(name string) (string, error) {
	p, err := GetSessionDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, name+".yml"), err
}

// GetSessionDir returns `os.UserSession.Dir() + "treesource" + "sessions"`
func GetSessionDir() (string, error) {
	s, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	s = filepath.Join(s, "treesource", "sessions")

	return s, err
}

// EnsureSession ensures GetSessionDir's path is available and that there is an available session.yaml file within it.
func EnsureSession(name string) error {
	p, err := GetSessionDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(p, 0755)
	if err != nil {
		return err
	}

	p = filepath.Join(p, name+".yml")

	if _, err := os.Stat(p); err == nil {
		// exists
	} else if os.IsNotExist(err) {
		// does not exist
		b, err := yaml.Marshal(&Session{})
		if err != nil {
			return err
		}
		// Make basic session file.
		err = os.WriteFile(p, b, 0755)
		if err != nil {
			return err
		}
	} else {
		// other error
		return err
	}
	return nil
}

// Views
func (s *Session) AddDirectoryView(u uuid.UUID) error {
	s.Views.Directories = append(s.Views.Directories, &DirectoryView{
		UUID:      uuid.New(),
		Directory: u,
	})
	s.Emit(EventViewDirectoryAdd, ViewDirectoryAddEvent{
		View: s.Views.Directories[len(s.Views.Directories)-1],
	})
	s.PendingSave()
	return nil
}

func (s *Session) RemoveDirectoryView(u uuid.UUID) error {
	for i, d := range s.Views.Directories {
		if d.UUID.String() == u.String() {
			s.Views.Directories = append(s.Views.Directories[:i], s.Views.Directories[i+1:]...)
			s.Emit(EventViewDirectoryRemove, ViewDirectoryRemoveEvent{
				View: d,
			})
			s.PendingSave()
			return nil
		}
	}
	return &MissingDirectoryViewError{
		uuid: u,
	}
}

func (s *Session) AddTagsView(tags []string) error {
	s.Views.Tags = append(s.Views.Tags, &TagsView{
		UUID: uuid.New(),
		Tags: tags,
	})
	s.Emit(EventViewTagsAdd, ViewTagsAddEvent{
		View: s.Views.Tags[len(s.Views.Tags)-1],
	})
	s.PendingSave()
	return nil
}

func (s *Session) RemoveTagsView(u uuid.UUID) error {
	for i, t := range s.Views.Tags {
		if t.UUID.String() == u.String() {
			s.Views.Tags = append(s.Views.Tags[:i], s.Views.Tags[i+1:]...)
			s.Emit(EventViewTagsRemove, ViewTagsRemoveEvent{
				View: t,
			})
			s.PendingSave()
			return nil
		}
	}
	return &MissingTagsViewError{
		uuid: u,
	}
}
