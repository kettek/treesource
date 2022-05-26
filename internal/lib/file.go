package lib

import "time"

type FileInfo struct {
	Name        string
	Path        string
	Size        int64
	Mode        uint32
	Permissions string
	Type        string
	Special     bool
	ModTime     time.Time
	Mimetype    string
}

type Thumbnail struct {
	Format string `json:"Format"`
	Bytes  []byte `json:"Bytes"`
}

type ThumbnailOptions struct {
	MaxWidth  int    `json:"MaxWidth"`
	MaxHeight int    `json:"MaxHeight"`
	Method    string `json:"Method"` // should be NearestNeighbor, CatmullRom, ApproxBiLinear
}
