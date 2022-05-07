package lib

import "github.com/google/uuid"

type DirectoryView struct {
	UUID      uuid.UUID `json:"uuid" yaml:"uuid"`
	Directory uuid.UUID `json:"directory" yaml:"directory"`
	WD        string    `json:"wd" yaml:"wd"`
	Selected  []string  `json:"selected" yaml:"selected"`
}

type TagsView struct {
	UUID     uuid.UUID `json:"uuid" yaml:"uuid"`
	Tags     []string  `json:"tags" yaml:"tags"`
	Selected []string  `json:"selected" yaml:"selected"`
}
