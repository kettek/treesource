package xdgicons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Theme struct {
	Root string
}

type Icon struct {
	Ext   string `json:"Ext"`
	Bytes []byte `json:"Bytes"`
}

// FIXME: This doesn't actually do xdg.
func (t *Theme) GetIcon(icon string, size int, scale int) (Icon, error) {
	parts := strings.SplitN(icon, "/", 2)
	sizer := ""
	if scale != 0 && scale != 1 {
		sizer = fmt.Sprintf("%d@%d", size, scale)
	} else {
		sizer = fmt.Sprintf("%d", size)
	}
	p := filepath.Join(t.Root, parts[0], sizer, parts[1]+".svg")
	b, err := os.ReadFile(p)
	if err != nil {
		return Icon{}, err
	}
	return Icon{
		Ext:   ".svg",
		Bytes: b,
	}, nil
}
