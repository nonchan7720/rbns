// +build !embed

package restserver

import (
	"flag"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/n-creativesystem/api-rbac/handler/restserver/internal/file"
)

var (
	root    string
	indexes bool
)

func init() {
	flag.StringVar(&root, "staticRoot", "./static/dist", "web ui static root")
	flag.BoolVar(&indexes, "staticIndexes", false, "web ui static allow indexes")
}

const INDEX = "index.html"

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func newFileSystem(root string, indexes bool) file.ServeFileSystem {
	return &localFileSystem{
		FileSystem: http.Dir(root),
		root:       root,
		indexes:    indexes,
	}
}
