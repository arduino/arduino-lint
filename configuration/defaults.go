package configuration

// The default configuration settings.

import (
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
)

func setDefaults() {
	superprojectType = projecttype.All
	recursive = true
	// TODO: targetPath defaults to current path
}

// Default check modes for each superproject type
// Subprojects use the same check modes as the superproject
var defaultCheckModes = map[projecttype.Type]map[checkmode.Type]bool{
	projecttype.Sketch: {
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.Library: {
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: true,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.Platform: {
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.PackageIndex: {
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
}
