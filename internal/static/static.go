package static

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var FS embed.FS

func DistFS(fsys fs.FS) (fs.FS, error) {
	fs, err := fs.Sub(fsys, "dist")
	if err != nil {
		return nil, err
	}

	return fs, nil
}
