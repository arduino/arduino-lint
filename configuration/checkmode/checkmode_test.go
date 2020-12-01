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

package checkmode

import (
	"reflect"
	"testing"

	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/stretchr/testify/assert"
)

func TestMode(t *testing.T) {
	defaultCheckModes := map[projecttype.Type]map[Type]bool{
		projecttype.Sketch: {
			LibraryManagerSubmission: false,
			LibraryManagerIndexed:    false,
			Official:                 false,
			All:                      true,
		},
		projecttype.Library: {
			LibraryManagerSubmission: true,
			LibraryManagerIndexed:    false,
			Official:                 false,
			All:                      true,
		},
	}

	customCheckModes := make(map[Type]bool)

	testProjectType := projecttype.Library

	mergedCheckModes := Modes(defaultCheckModes, customCheckModes, testProjectType)

	assert.True(t, reflect.DeepEqual(defaultCheckModes[testProjectType], mergedCheckModes), "Default configuration should be used when no custom configuration was set.")

	testCheckMode := Official
	customCheckModes[testCheckMode] = !defaultCheckModes[testProjectType][testCheckMode]
	mergedCheckModes = Modes(defaultCheckModes, customCheckModes, testProjectType)
	assert.Equal(t, customCheckModes[testCheckMode], mergedCheckModes[testCheckMode], "Should be set to custom value")
}
