package lib

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Directory represents a source directory to pull files from.
type Directory struct {
	Emitter
	UUID       uuid.UUID         `json:"UUID" yaml:"UUID"`
	Path       string            `json:"Path" yaml:"Path"`
	Entries    []*DirectoryEntry `json:"Entries" yaml:"Entries"`
	IgnoreDot  bool              `json:"IgnoreDot" yaml:"IgnoreDot"`
	SyncOnLoad bool              `json:"SyncOnLoad" yaml:"SyncOnLoad"`
}

// Entry retrieves an entry matching the given name.
func (d *Directory) Entry(name string) *DirectoryEntry {
	for _, e := range d.Entries {
		if e.Path == name {
			return e
		}
	}
	return nil
}

// SyncEntries synchronizes the directory's entries with the on-disk file structure. Emits: sync, synced, add, found, missing
func (d *Directory) SyncEntries() error {
	d.Emit("sync", &DirectorySyncEvent{
		UUID: d.UUID,
		Name: d.Path,
	})
	unmatchedEntries := make([]*DirectoryEntry, len(d.Entries))
	copy(unmatchedEntries, d.Entries)
	var errors []error
	var walk func(local, name string)
	walk = func(local, name string) {
		entries, err := os.ReadDir(name)
		if err != nil {
			errors = append(errors, err)
			return
		}
		for _, e := range entries {
			fullpath := filepath.Join(name, e.Name())
			localpath := filepath.Join(local, e.Name())
			if d.IgnoreDot && e.Name()[0] == '.' {
				continue
			}
			if e.IsDir() {
				walk(localpath, fullpath)
			} else {
				if d.Entry(localpath) == nil {
					d.Entries = append(d.Entries, &DirectoryEntry{
						Path: localpath,
					})
					d.Emit("add", &DirectoryFileAddEvent{
						Name: d.Path,
						File: d.Entries[len(d.Entries)-1],
					})
				} else {
					for i, e := range unmatchedEntries {
						if e.Path == localpath {
							//unmatchedEntries = slices.Delete(unmatchedEntries, i, i+1)
							unmatchedEntries[i] = unmatchedEntries[len(unmatchedEntries)-1]
							unmatchedEntries = unmatchedEntries[:len(unmatchedEntries)-1]
							// Mark found entries as not missing if they were marked as such.
							if e.Missing {
								e.Missing = false
								d.Emit("found", &DirectoryFileFoundEvent{
									Name: d.Path,
									File: e,
								})
							}
							break
						}
					}
				}
			}
		}

		// Mark any unmatched entries as missing.
		for _, e := range unmatchedEntries {
			for i, e2 := range d.Entries {
				if e.Path == e2.Path {
					d.Entries[i].Missing = true
					d.Emit("missing", &DirectoryFileMissingEvent{
						Name: d.Path,
						File: e,
					})
					break
				}
			}
		}
	}

	walk("", d.Path)

	var err error
	if len(errors) > 0 {
		err = &SyncError{errors}
	}

	d.Emit("synced", &DirectorySyncedEvent{
		UUID:  d.UUID,
		Name:  d.Path,
		Error: err,
	})
	return err
}

// DirectoryEntry represents an entry in a treesource directory.
type DirectoryEntry struct {
	// Path is relative to the owning Directory's path.
	Path string   `json:"Path" yaml:"Path"`
	Tags []string `json:"Tags,omitempty" yaml:"Tags,omitempty"`
	// Missing represents if the entry is referring to a file that no longer exists.
	Missing bool `json:"Missing,omitempty" yaml:"Missing,omitempty"`
}
