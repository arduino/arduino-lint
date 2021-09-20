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

package rule

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/util/test"
	"github.com/stretchr/testify/assert"
)

func Test_shouldRun(t *testing.T) {
	testTables := []struct {
		testName              string
		ruleProjectType       projecttype.Type
		ruleSuperprojectType  projecttype.Type
		projectType           projecttype.Type
		superprojectType      projecttype.Type
		disableModes          []rulemode.Type
		enableModes           []rulemode.Type
		libraryManagerSetting string
		complianceSetting     string
		shouldRunAssertion    assert.BoolAssertionFunc
		errorAssertion        assert.ErrorAssertionFunc
	}{
		{"Project type mismatch", projecttype.Library, projecttype.All, projecttype.Sketch, projecttype.Sketch, []rulemode.Type{}, []rulemode.Type{}, "false", "specification", assert.False, assert.NoError},
		{"Superproject type mismatch", projecttype.Sketch, projecttype.Library, projecttype.Sketch, projecttype.Sketch, []rulemode.Type{}, []rulemode.Type{}, "false", "specification", assert.False, assert.NoError},
		{"Disable mode match", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{}, "submit", "specification", assert.False, assert.NoError},
		{"Enable mode match", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{}, []rulemode.Type{rulemode.LibraryManagerSubmission}, "submit", "specification", assert.True, assert.NoError},
		{"Disable mode default", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.Default}, []rulemode.Type{rulemode.LibraryManagerSubmission}, "update", "specification", assert.False, assert.NoError},
		{"Disable mode default override", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.Default}, []rulemode.Type{rulemode.LibraryManagerSubmission}, "submit", "specification", assert.True, assert.NoError},
		{"Enable mode default", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{rulemode.Default}, "update", "specification", assert.True, assert.NoError},
		{"Enable mode default override", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{rulemode.Default}, "submit", "specification", assert.False, assert.NoError},
		{"Unable to resolve", projecttype.Library, projecttype.All, projecttype.Library, projecttype.Library, []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{rulemode.LibraryManagerIndexed}, "false", "specification", assert.False, assert.Error},
	}

	flags := test.ConfigurationFlags()

	for _, testTable := range testTables {
		flags.Set("library-manager", testTable.libraryManagerSetting)
		flags.Set("compliance", testTable.complianceSetting)

		configuration.Initialize(flags, []string{"/foo"})

		ruleConfiguration := ruleconfiguration.Type{
			ProjectType:      testTable.ruleProjectType,
			SuperprojectType: testTable.ruleSuperprojectType,
			DisableModes:     testTable.disableModes,
			EnableModes:      testTable.enableModes,
		}

		project := project.Type{
			ProjectType:      testTable.projectType,
			SuperprojectType: testTable.superprojectType,
		}
		run, err := shouldRun(ruleConfiguration, project)
		testTable.errorAssertion(t, err, testTable.testName)
		if err == nil {
			testTable.shouldRunAssertion(t, run, testTable.testName)
		}
	}
}
