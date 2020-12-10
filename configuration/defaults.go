// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package configuration

// The default configuration settings.

import (
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/sirupsen/logrus"
)

// Default check modes for each superproject type.
// Subprojects use the same check modes as the superproject.
var defaultCheckModes = map[projecttype.Type]map[checkmode.Type]bool{
	projecttype.Sketch: {
		checkmode.Strict:                   false,
		checkmode.Specification:            true,
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.Library: {
		checkmode.Strict:                   false,
		checkmode.Specification:            true,
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: true,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.Platform: {
		checkmode.Strict:                   false,
		checkmode.Specification:            true,
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
	projecttype.PackageIndex: {
		checkmode.Strict:                   false,
		checkmode.Specification:            true,
		checkmode.Permissive:               false,
		checkmode.LibraryManagerSubmission: false,
		checkmode.LibraryManagerIndexed:    false,
		checkmode.Official:                 false,
	},
}

var defaultLogLevel = logrus.FatalLevel
