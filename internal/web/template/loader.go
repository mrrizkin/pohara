package template

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/nikolalohinski/gonja/v2/loaders"
)

type httpFilesystemLoader struct {
	fs      http.FileSystem
	baseDir string
}

func newHttpFileSystemLoader(fs http.FileSystem, baseDir string) (loaders.Loader, error) {
	httpFs := &httpFilesystemLoader{
		fs:      fs,
		baseDir: baseDir,
	}
	if fs == nil {
		err := errors.New("httpfs cannot be nil")
		return nil, err
	}
	return httpFs, nil
}

func (h *httpFilesystemLoader) Resolve(name string) (string, error) {
	return name, nil
}

// Get returns an io.Reader where the template's content can be read from.
func (h *httpFilesystemLoader) Read(path string) (io.Reader, error) {
	fullPath := path
	if h.baseDir != "" {
		fullPath = fmt.Sprintf(
			"%s/%s",
			h.baseDir,
			fullPath,
		)
	}

	return h.fs.Open(fullPath)
}

func (h *httpFilesystemLoader) Inherit(from string) (loaders.Loader, error) {
	hfs := &httpFilesystemLoader{
		fs:      h.fs,
		baseDir: h.baseDir,
	}

	return hfs, nil
}
