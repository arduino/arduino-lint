package platform

import (
	"github.com/arduino/go-paths-helper"
)

var empty struct{}

var configurationFilenames = map[string]struct{}{
	"boards.txt":         empty,
	"boards.local.txt":   empty,
	"platform.txt":       empty,
	"platform.local.txt": empty,
	"programmers.txt":    empty,
}

func IsConfigurationFile(filePath *paths.Path) bool {
	_, isConfigurationFile := configurationFilenames[filePath.Base()]
	if isConfigurationFile {
		return true
	}
	return false
}

var requiredConfigurationFilenames = map[string]struct{}{
	// Arduino platforms must always contain a boards.txt
	"boards.txt": empty,
}

func IsRequiredConfigurationFile(filePath *paths.Path) bool {
	_, isRequiredConfigurationFile := requiredConfigurationFilenames[filePath.Base()]
	if isRequiredConfigurationFile {
		return true
	}
	return false
}
