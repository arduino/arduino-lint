package sketch

import (
	"github.com/arduino/arduino-cli/arduino/globals"
	"github.com/arduino/go-paths-helper"
)

func HasMainFileValidExtension(filePath *paths.Path) bool {
	_, hasMainFileValidExtension := globals.MainFileValidExtensions[filePath.Ext()]
	if hasMainFileValidExtension {
		return true
	}
	return false
}
