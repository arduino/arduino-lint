package configuration

import (
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
)

func Initialize() {
	setDefaults()
	// TODO configuration according to command line input
	// TODO validate target path value, exit if not found
	targetPath = paths.New("e:/electronics/arduino/libraries/arduino-check-test-library")
	superprojectType = projecttype.Library
	customCheckModes[checkmode.Permissive] = false
}

var customCheckModes = make(map[checkmode.Type]bool)

func CheckModes(superprojectType projecttype.Type) map[checkmode.Type]bool {
	return checkmode.Modes(defaultCheckModes, customCheckModes, superprojectType)
}

var superprojectType projecttype.Type

func SuperprojectType() projecttype.Type {
	return superprojectType
}

var recursive bool

func Recursive() bool {
	return recursive
}

var targetPath *paths.Path

func TargetPath() *paths.Path {
	return targetPath
}
