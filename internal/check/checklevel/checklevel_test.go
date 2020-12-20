// This file is part of arduino-lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package checklevel

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/check/checkconfigurations"
	"github.com/arduino/arduino-lint/internal/check/checkresult"
	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/checkmode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/util/test"
	"github.com/stretchr/testify/assert"
)

func TestCheckLevel(t *testing.T) {
	testTables := []struct {
		testName              string
		infoModes             []checkmode.Type
		warningModes          []checkmode.Type
		errorModes            []checkmode.Type
		checkResult           checkresult.Type
		libraryManagerSetting string
		permissiveSetting     string
		expectedLevel         Type
		errorAssertion        assert.ErrorAssertionFunc
	}{
		{"Non-fail", []checkmode.Type{}, []checkmode.Type{}, []checkmode.Type{checkmode.LibraryManagerSubmission}, checkresult.Skip, "submit", "false", Notice, assert.NoError},
		{"Error", []checkmode.Type{}, []checkmode.Type{}, []checkmode.Type{checkmode.LibraryManagerSubmission}, checkresult.Fail, "submit", "false", Error, assert.NoError},
		{"Warning", []checkmode.Type{}, []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{}, checkresult.Fail, "submit", "false", Warning, assert.NoError},
		{"Info", []checkmode.Type{checkmode.LibraryManagerSubmission}, []checkmode.Type{}, []checkmode.Type{}, checkresult.Fail, "submit", "false", Info, assert.NoError},
		{"Default to Error", []checkmode.Type{}, []checkmode.Type{}, []checkmode.Type{checkmode.Default}, checkresult.Fail, "submit", "false", Error, assert.NoError},
		{"Default to Warning", []checkmode.Type{}, []checkmode.Type{checkmode.Default}, []checkmode.Type{}, checkresult.Fail, "submit", "false", Warning, assert.NoError},
		{"Default to Info", []checkmode.Type{checkmode.Default}, []checkmode.Type{}, []checkmode.Type{}, checkresult.Fail, "submit", "false", Info, assert.NoError},
		{"Default override", []checkmode.Type{checkmode.Default}, []checkmode.Type{}, []checkmode.Type{checkmode.LibraryManagerSubmission}, checkresult.Fail, "submit", "false", Error, assert.NoError},
		{"Unable to resolve", []checkmode.Type{}, []checkmode.Type{}, []checkmode.Type{}, checkresult.Fail, "submit", "false", Info, assert.Error},
	}

	flags := test.ConfigurationFlags()

	for _, testTable := range testTables {
		flags.Set("library-manager", testTable.libraryManagerSetting)
		flags.Set("permissive", testTable.permissiveSetting)

		configuration.Initialize(flags, []string{"/foo"})

		checkConfiguration := checkconfigurations.Type{
			InfoModes:    testTable.infoModes,
			WarningModes: testTable.warningModes,
			ErrorModes:   testTable.errorModes,
		}

		checkedProject := project.Type{
			SuperprojectType: projecttype.Sketch,
		}

		level, err := CheckLevel(checkConfiguration, testTable.checkResult, checkedProject)
		testTable.errorAssertion(t, err, testTable.testName)
		if err == nil {
			assert.Equal(t, testTable.expectedLevel, level, testTable.testName)
		}
	}
}
