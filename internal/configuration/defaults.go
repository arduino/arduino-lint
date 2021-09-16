// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License, either
// version 3 of the License, or (at your option) any later version.
// This license covers the main part of Arduino Lint.
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
	"os"

	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
)

// Default rule modes for each superproject type.
// Subprojects use the same rule modes as the superproject.
var defaultRuleModes = map[projecttype.Type]map[rulemode.Type]bool{
	projecttype.Sketch: {
		rulemode.Strict:                   false,
		rulemode.Specification:            true,
		rulemode.Permissive:               false,
		rulemode.LibraryManagerSubmission: false,
		rulemode.LibraryManagerIndexed:    false,
		rulemode.LibraryManagerIndexing:   false,
		rulemode.Official:                 false,
	},
	projecttype.Library: {
		rulemode.Strict:                   false,
		rulemode.Specification:            true,
		rulemode.Permissive:               false,
		rulemode.LibraryManagerSubmission: true,
		rulemode.LibraryManagerIndexed:    false,
		rulemode.LibraryManagerIndexing:   false,
		rulemode.Official:                 false,
	},
	projecttype.Platform: {
		rulemode.Strict:                   false,
		rulemode.Specification:            true,
		rulemode.Permissive:               false,
		rulemode.LibraryManagerSubmission: false,
		rulemode.LibraryManagerIndexed:    false,
		rulemode.LibraryManagerIndexing:   false,
		rulemode.Official:                 false,
	},
	projecttype.PackageIndex: {
		rulemode.Strict:                   false,
		rulemode.Specification:            true,
		rulemode.Permissive:               false,
		rulemode.LibraryManagerSubmission: false,
		rulemode.LibraryManagerIndexed:    false,
		rulemode.LibraryManagerIndexing:   false,
		rulemode.Official:                 false,
	},
}

var defaultLogOutput = os.Stderr
