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

package rulelevel

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/arduino-lint/internal/util/test"
	"github.com/stretchr/testify/assert"
)

func TestRuleLevel(t *testing.T) {
	testTables := []struct {
		testName              string
		infoModes             []rulemode.Type
		warningModes          []rulemode.Type
		errorModes            []rulemode.Type
		ruleResult            ruleresult.Type
		libraryManagerSetting string
		permissiveSetting     string
		expectedLevel         Type
		errorAssertion        assert.ErrorAssertionFunc
	}{
		{"Non-fail", []rulemode.Type{}, []rulemode.Type{}, []rulemode.Type{rulemode.LibraryManagerSubmission}, ruleresult.Skip, "submit", "false", Notice, assert.NoError},
		{"Error", []rulemode.Type{}, []rulemode.Type{}, []rulemode.Type{rulemode.LibraryManagerSubmission}, ruleresult.Fail, "submit", "false", Error, assert.NoError},
		{"Warning", []rulemode.Type{}, []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{}, ruleresult.Fail, "submit", "false", Warning, assert.NoError},
		{"Info", []rulemode.Type{rulemode.LibraryManagerSubmission}, []rulemode.Type{}, []rulemode.Type{}, ruleresult.Fail, "submit", "false", Info, assert.NoError},
		{"Default to Error", []rulemode.Type{}, []rulemode.Type{}, []rulemode.Type{rulemode.Default}, ruleresult.Fail, "submit", "false", Error, assert.NoError},
		{"Default to Warning", []rulemode.Type{}, []rulemode.Type{rulemode.Default}, []rulemode.Type{}, ruleresult.Fail, "submit", "false", Warning, assert.NoError},
		{"Default to Info", []rulemode.Type{rulemode.Default}, []rulemode.Type{}, []rulemode.Type{}, ruleresult.Fail, "submit", "false", Info, assert.NoError},
		{"Default override", []rulemode.Type{rulemode.Default}, []rulemode.Type{}, []rulemode.Type{rulemode.LibraryManagerSubmission}, ruleresult.Fail, "submit", "false", Error, assert.NoError},
		{"Unable to resolve", []rulemode.Type{}, []rulemode.Type{}, []rulemode.Type{}, ruleresult.Fail, "submit", "false", Info, assert.Error},
	}

	flags := test.ConfigurationFlags()

	for _, testTable := range testTables {
		flags.Set("library-manager", testTable.libraryManagerSetting)
		flags.Set("permissive", testTable.permissiveSetting)

		configuration.Initialize(flags, []string{"/foo"})

		ruleConfiguration := ruleconfiguration.Type{
			InfoModes:    testTable.infoModes,
			WarningModes: testTable.warningModes,
			ErrorModes:   testTable.errorModes,
		}

		lintedProject := project.Type{
			SuperprojectType: projecttype.Sketch,
		}

		level, err := RuleLevel(ruleConfiguration, testTable.ruleResult, lintedProject)
		testTable.errorAssertion(t, err, testTable.testName)
		if err == nil {
			assert.Equal(t, testTable.expectedLevel, level, testTable.testName)
		}
	}
}
