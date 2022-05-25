package lib

import (
	"fmt"
	"treesource/internal/do"

	"github.com/google/uuid"
)

// AddDirectoryAction adds the given directory at the provided index.
type AddDirectoryAction struct {
	Directory Directory
	Index     int
}

// Apply does the obvious.
func (a *AddDirectoryAction) Apply(p *Project) {
	fmt.Println("action: apply add dir")
	if len(p.Directories) == a.Index {
		p.Directories = append(p.Directories, *a.Directory.Clone())
	} else {
		p.Directories = append(p.Directories[:a.Index+1], p.Directories[a.Index:]...)
		p.Directories[a.Index] = *a.Directory.Clone()
	}
	p.Emit(EventDirectoryAdd, DirectoryAddEvent{
		UUID:       a.Directory.UUID,
		Path:       a.Directory.Path,
		IgnoreDot:  a.Directory.IgnoreDot,
		SyncOnLoad: a.Directory.SyncOnLoad,
		Separator:  a.Directory.Separator,
	})
	// Need to rehook garbage.
	p.InitDirectory(&p.Directories[a.Index])
	// Seems reasonable enough to emit all entries on load.
	p.Directories[a.Index].EmitAllEntries()
}

// Unapply does the obvious.
func (a *AddDirectoryAction) Unapply(p *Project) {
	fmt.Println("action: unapply add dir")
	p.Directories = append(p.Directories[:a.Index], p.Directories[a.Index+1:]...)
	p.Emit(EventDirectoryRemove, DirectoryRemoveEvent{
		UUID: a.Directory.UUID,
	})
}

// RemoveDirectoryAction removes the directory at the given index.
type RemoveDirectoryAction struct {
	Directory Directory
	Index     int
}

// Apply does the obvious.
func (a *RemoveDirectoryAction) Apply(p *Project) {
	fmt.Println("action: apply remove dir")
	for i, d := range p.Directories {
		if d.UUID.String() == a.Directory.UUID.String() {
			p.Directories = append(p.Directories[:i], p.Directories[i+1:]...)
			p.Emit(EventDirectoryRemove, DirectoryRemoveEvent{
				UUID: d.UUID,
			})
			break
		}
	}
}

// Unapply does the obvious.
func (a *RemoveDirectoryAction) Unapply(p *Project) {
	fmt.Println("action: unapply remove dir")
	if len(p.Directories) == a.Index {
		p.Directories = append(p.Directories, *a.Directory.Clone())
	} else {
		p.Directories = append(p.Directories[:a.Index+1], p.Directories[a.Index:]...)
		p.Directories[a.Index] = *a.Directory.Clone()
	}
	p.Emit(EventDirectoryAdd, DirectoryAddEvent{
		UUID:       a.Directory.UUID,
		Path:       a.Directory.Path,
		IgnoreDot:  a.Directory.IgnoreDot,
		SyncOnLoad: a.Directory.SyncOnLoad,
		Separator:  a.Directory.Separator,
	})
	// Need to rehook garbage.
	p.InitDirectory(&p.Directories[a.Index])
	// Seems reasonable enough to emit all entries on load.
	p.Directories[a.Index].EmitAllEntries()
}

type SyncDirectoryAction struct {
	// State of added or removed directory files.
}

func (a *SyncDirectoryAction) Apply(p *Project) {
	fmt.Println("action: apply sync dir")
}

func (a *SyncDirectoryAction) Unapply(p *Project) {
	fmt.Println("action: unapply sync dir")
}

type UpdateEntryAction struct {
	UUID     uuid.UUID
	Entry    DirectoryEntry
	path     string
	previous DirectoryEntry
}

func (a *UpdateEntryAction) Apply(p *Project) {
	dir, err := p.GetDirectoryByUUID(a.UUID)
	if err != nil {
		return
	}
	entry := dir.Entry(a.path)
	if entry == nil {
		return
	}
	a.previous = entry.Clone()
	entry.Subsume(a.Entry)
	dir.Emit(EventDirectoryEntryUpdate, DirectoryEntryUpdateEvent{
		UUID:  a.UUID,
		Entry: entry,
	})
}

func (a *UpdateEntryAction) Unapply(p *Project) {
	dir, err := p.GetDirectoryByUUID(a.UUID)
	if err != nil {
		return
	}
	entry := dir.Entry(a.Entry.Path)
	if entry == nil {
		return
	}
	entry.Subsume(a.previous)
	dir.Emit(EventDirectoryEntryUpdate, DirectoryEntryUpdateEvent{
		UUID:  a.UUID,
		Entry: entry,
	})
}

type RemoveEntryAction struct {
	UUID     uuid.UUID
	Path     string
	index    int
	previous *DirectoryEntry
}

func (a *RemoveEntryAction) Apply(p *Project) {
	dir, err := p.GetDirectoryByUUID(a.UUID)
	if err != nil {
		return
	}
	entry, index := dir.Remove(a.Path)
	if entry == nil {
		return
	}
	a.previous = entry
	a.index = index
	dir.Emit(EventDirectoryEntryRemove, DirectoryEntryRemoveEvent{
		UUID:  a.UUID,
		Entry: a.previous,
	})
}

func (a *RemoveEntryAction) Unapply(p *Project) {
	dir, err := p.GetDirectoryByUUID(a.UUID)
	if err != nil {
		return
	}
	if a.previous == nil {
		return
	}
	if len(dir.Entries) == a.index {
		dir.Entries = append(dir.Entries, a.previous)
	} else {
		dir.Entries = append(dir.Entries[:a.index+1], dir.Entries[a.index:]...)
		dir.Entries[a.index] = a.previous
	}
	dir.Emit(EventDirectoryEntryAdd, DirectoryEntryAddEvent{
		UUID:  a.UUID,
		Entry: a.previous,
	})
}

// GroupedAction represents a collection of actions.
type GroupedAction struct {
	Actions []do.Action[*Project]
}

// Apply applies the contained actions from the start to the end.
func (a *GroupedAction) Apply(p *Project) {
	for _, a2 := range a.Actions {
		a2.Apply(p)
	}
}

// Unapply unapplies the contains actions from the end to the start.
func (a *GroupedAction) Unapply(p *Project) {
	for i := len(a.Actions); i > 0; i-- {
		a.Actions[i].Unapply(p)
	}
}
