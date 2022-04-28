package lib

import "github.com/google/uuid"

type DirView struct {
	UUID uuid.UUID `json:"uuid" yaml:"uuid"`
	Root string    `json:"root" yaml:"root"`
}
