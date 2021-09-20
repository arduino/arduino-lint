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

package rulefunction

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	testDataPath = paths.New(workingDirectory, "testdata", "general")
}

type ruleFunctionTestTable struct {
	testName            string
	projectFolderName   string
	expectedRuleResult  ruleresult.Type
	expectedOutputQuery string
}

func checkRuleFunction(ruleFunction Type, testTables []ruleFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             testDataPath.Join(testTable.projectFolderName),
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.Library,
		}

		projectdata.Initialize(testProject)

		result, output := ruleFunction()
		assert.Equal(t, testTable.expectedRuleResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), fmt.Sprintf("%s (output: %s, assertion regex: %s)", testTable.testName, output, testTable.expectedOutputQuery))
	}
}

func TestMissingReadme(t *testing.T) {
	testTables := []ruleFunctionTestTable{
		{"Readme", "readme", ruleresult.Pass, ""},
		{"No readme", "no-readme", ruleresult.Fail, ""},
	}

	checkRuleFunction(MissingReadme, testTables, t)
}

func TestMissingLicenseFile(t *testing.T) {
	testTables := []ruleFunctionTestTable{
		{"Has license", "license-file", ruleresult.Pass, ""},
		{"Has license in subfolder", "license-file-in-subfolder", ruleresult.Fail, ""},
		{"No license", "no-license-file", ruleresult.Fail, ""},
	}

	checkRuleFunction(MissingLicenseFile, testTables, t)
}

func TestIncorrectArduinoDotHFileNameCase(t *testing.T) {
	testTables := []ruleFunctionTestTable{
		{"Incorrect, angle brackets", "arduino.h-angle", ruleresult.Fail, ""},
		{"Incorrect, quotes", "arduino.h-quote", ruleresult.Fail, ""},
		{"Correct case", "Arduino.h", ruleresult.Pass, ""},
	}

	checkRuleFunction(IncorrectArduinoDotHFileNameCase, testTables, t)
}
