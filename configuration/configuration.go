package configuration

import (
	"github.com/arduino/arduino-check/projects/projecttype"
	"github.com/arduino/go-paths-helper"
)

type CheckMode int

const (
	Permissive CheckMode = iota
	LibraryManagerSubmission
	LibraryManagerIndexed
	Official
	Default
)

func Initialize() {
	setDefaults()
	// TODO configuration according to command line input
	// TODO validate target path value, exit if not found
	targetPath = paths.New("e:/electronics/arduino/libraries/arduino-check-test-library")
	superprojectType = projecttype.Library
}

var customCheckModes map[CheckMode]bool

func CheckModes(superprojectType projecttype.Type) map[CheckMode]bool {
	checkModes := make(map[CheckMode]bool)

	// Merge the default settings with any custom settings specified by the user
	for key, defaultValue := range defaultCheckModes[superprojectType] {
		customCheckModeValue, customCheckModeIsConfigured := customCheckModes[key]
		if customCheckModeIsConfigured {
			checkModes[key] = customCheckModeValue
		} else {
			checkModes[key] = defaultValue
		}
	}

	// This mode is always enabled
	checkModes[Default] = true

	return checkModes
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
