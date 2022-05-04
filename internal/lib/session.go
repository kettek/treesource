package lib

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Session represents the local session for treesource.
type Session struct {
	path    string
	Project string `json:"project" yaml:"project"`
	Views   struct {
		Directories []DirectoryView
		Tags        []TagsView
	}
	canceledSave chan struct{}
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
		canceledSave: make(chan struct{}),
	}

	p, err := GetSessionPath(name)
	if err != nil {
		return nil, err
	}
	c.path = p

	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)

	return nil, err
}

// GetSessionPath returns the path for the target session.
func GetSessionPath(name string) (string, error) {
	p, err := GetSessionDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, name+".yml"), err
}

// GetSessionDir returns `os.UserSession.Dir() + "treesource"`
func GetSessionDir() (string, error) {
	s, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	s = filepath.Join(s, "treesource")

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
