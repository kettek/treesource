package lib

import (
	"fmt"
	"treesource/internal/do"
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
		UUID: a.Directory.UUID,
		Path: a.Directory.Path,
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
		UUID: a.Directory.UUID,
		Path: a.Directory.Path,
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
