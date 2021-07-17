// +build embed

package restserver

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/n-creativesystem/rbnc/handler/restserver/internal/file"
)

//go:embed static/dist/*
var assets embed.FS

const INDEX = "index.html"

type embedFileSystem struct {
	http.FileSystem
	prefix string
}

func (e embedFileSystem) Exists(prefix, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func embedFolder(fsEmbed embed.FS, targetPath string) file.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func newFileSystem(root string, indexes bool) file.ServeFileSystem {
	return embedFolder(assets, "static/dist")
}
