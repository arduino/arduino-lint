// Package util contains general purpose utility code.
package util

import (
	"net/url"
	"path/filepath"

	"github.com/arduino/go-paths-helper"
)

// PathURI returns the URI representation of the path argument.
func PathURI(path *paths.Path) string {
	uriFriendlyPath := filepath.ToSlash(path.String())
	pathURI := url.URL{
		Scheme: "file",
		Path:   uriFriendlyPath,
	}

	return pathURI.String()
}
