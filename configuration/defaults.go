package configuration

import "github.com/arduino/arduino-check/projects/projecttype"

func setDefaults() {
	superprojectType = projecttype.All
	recursive = true
	// TODO: targetPath defaults to current path
}

// Default check modes for each superproject type
// Subprojects use the same check modes as the superproject
var defaultCheckModes = map[projecttype.Type]map[CheckMode]bool{
	projecttype.Sketch: {
		Permissive:               false,
		LibraryManagerSubmission: false,
		LibraryManagerIndexed:    false,
		Official:                 false,
	},
	projecttype.Library: {
		Permissive:               false,
		LibraryManagerSubmission: true,
		LibraryManagerIndexed:    false,
		Official:                 false,
	},
	projecttype.Platform: {
		Permissive:               false,
		LibraryManagerSubmission: false,
		LibraryManagerIndexed:    false,
		Official:                 false,
	},
	projecttype.PackageIndex: {
		Permissive:               false,
		LibraryManagerSubmission: false,
		LibraryManagerIndexed:    false,
		Official:                 false,
	},
}
