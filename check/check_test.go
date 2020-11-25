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

package check

import (
	"testing"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/util/test"
	"github.com/stretchr/testify/assert"
)

func Test_shouldRun(t *testing.T) {
	testTables := []struct {
		testName              string
		checkProjectType      projecttype.Type
		projectType           projecttype.Type
		disableModes          []checkmode.Type
		enableModes           []checkmode.Type
		libraryManagerSetting string
		complianceSetting     string
		shouldRunAssertion    assert.BoolAssertionFunc
		errorAssertion        assert.ErrorAssertionFunc
	}{
		{"Project type mismatch", projecttype.Library, projecttype.Sketch, []checkmode.Type{}, []checkmode.Type{}, "false", "specification", assert.False, assert.NoError},
		{"Disable mode match", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{}, "submit", "specification", assert.False, assert.NoError},
		{"Enable mode match", projecttype.Library, projecttype.Library, []checkmode.Type{}, []checkmode.Type{checkmode.LibraryManagerSubmission}, "submit", "specification", assert.True, assert.NoError},
		{"Disable mode default", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.Default}, []checkmode.Type{checkmode.LibraryManagerSubmission}, "update", "specification", assert.False, assert.NoError},
		{"Disable mode default override", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.Default}, []checkmode.Type{checkmode.LibraryManagerSubmission}, "submit", "specification", assert.True, assert.NoError},
		{"Enable mode default", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{checkmode.Default}, "update", "specification", assert.True, assert.NoError},
		{"Enable mode default override", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{checkmode.Default}, "submit", "specification", assert.False, assert.NoError},
		{"Unable to resolve", projecttype.Library, projecttype.Library, []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{checkmode.LibraryManagerIndexed}, "false", "specification", assert.False, assert.Error},
	}

	flags := test.ConfigurationFlags()

	for _, testTable := range testTables {
		flags.Set("library-manager", testTable.libraryManagerSetting)
		flags.Set("compliance", testTable.complianceSetting)

		configuration.Initialize(flags, []string{"/foo"})

		checkConfiguration := checkconfigurations.Type{
			ProjectType:  testTable.checkProjectType,
			DisableModes: testTable.disableModes,
			EnableModes:  testTable.enableModes,
		}

		project := project.Type{
			ProjectType: testTable.projectType,
		}
		run, err := shouldRun(checkConfiguration, project)
		testTable.errorAssertion(t, err, testTable.testName)
		if err == nil {
			testTable.shouldRunAssertion(t, run, testTable.testName)
		}
	}
}
