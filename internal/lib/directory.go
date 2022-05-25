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
	Separator  string            `json:"Separator" yaml:"-"`
	Entries    []*DirectoryEntry `json:"Entries" yaml:"Entries"`
	IgnoreDot  bool              `json:"IgnoreDot" yaml:"IgnoreDot"`
	SyncOnLoad bool              `json:"SyncOnLoad" yaml:"SyncOnLoad"`
}

func (d *Directory) Clone() *Directory {
	d2 := &Directory{}
	d2.UUID = d.UUID
	d2.Path = d.Path
	d2.Separator = string(os.PathSeparator)
	d2.IgnoreDot = d.IgnoreDot
	d2.SyncOnLoad = d.SyncOnLoad
	d2.Emitter = *NewEmitter()

	for _, e := range d.Entries {
		e2 := e.Clone()
		d2.Entries = append(d2.Entries, &e2)
	}

	return d2
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

// Remove removes an entry matching the given name and returns it and its index.
func (d *Directory) Remove(name string) (*DirectoryEntry, int) {
	for i, e := range d.Entries {
		if e.Path == name {
			d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
			return e, i
		}
	}
	return nil, -1
}

func (d *Directory) EmitAllEntries() {
	for _, e := range d.Entries {
		d.Emit("entry", &DirectoryEntryEvent{
			UUID:  d.UUID,
			Entry: e,
		})
	}
}

// SyncEntries synchronizes the directory's entries with the on-disk file structure. Emits: sync, synced, add, found, missing
func (d *Directory) SyncEntries() error {
	d.Emit("sync", &DirectorySyncEvent{
		UUID: d.UUID,
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
					d.Emit("add", &DirectoryEntryAddEvent{
						UUID:  d.UUID,
						Entry: d.Entries[len(d.Entries)-1],
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
								d.Emit("found", &DirectoryEntryFoundEvent{
									UUID:  d.UUID,
									Entry: e,
								})
							}
							break
						}
					}
				}
			}
		}
	}

	walk("", d.Path)

	// Mark any unmatched entries as missing.
	for _, e := range unmatchedEntries {
		for i, e2 := range d.Entries {
			if e.Path == e2.Path {
				d.Entries[i].Missing = true
				d.Emit("missing", &DirectoryEntryMissingEvent{
					UUID:  d.UUID,
					Entry: e,
				})
				break
			}
		}
	}

	var err error
	if len(errors) > 0 {
		err = &SyncError{errors}
	}

	d.Emit("synced", &DirectorySyncedEvent{
		UUID:  d.UUID,
		Error: err,
	})
	return err
}

// DirectoryEntry represents an entry in a treesource directory.
type DirectoryEntry struct {
	// Path is relative to the owning Directory's path.
	Path   string   `json:"Path" yaml:"Path"`
	Tags   []string `json:"Tags" yaml:"Tags,omitempty"`
	Rating float64  `json:"Rating" yaml:"Rating,omitempty"`
	// Missing represents if the entry is referring to a file that no longer exists.
	Missing bool `json:"Missing,omitempty" yaml:"Missing,omitempty"`
}

func (e *DirectoryEntry) Clone() (e2 DirectoryEntry) {
	e2.Path = e.Path
	copy(e2.Tags, e.Tags)
	e2.Rating = e.Rating
	e2.Missing = e.Missing
	return
}

func (e *DirectoryEntry) Subsume(o DirectoryEntry) {
	e.Path = o.Path
	e.Tags = o.Tags
	e.Rating = o.Rating
	e.Missing = o.Missing
}
