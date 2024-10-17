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

var sketchesTestDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	sketchesTestDataPath = paths.New(workingDirectory, "testdata", "sketches")
}

type sketchRuleFunctionTestTable struct {
	testName            string
	sketchFolderName    string
	expectedRuleResult  ruleresult.Type
	expectedOutputQuery string
}

func checkSketchRuleFunction(ruleFunction Type, testTables []sketchRuleFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		t.Run(testTable.testName, func(t *testing.T) {
			expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

			testProject := project.Type{
				Path:             sketchesTestDataPath.Join(testTable.sketchFolderName),
				ProjectType:      projecttype.Sketch,
				SuperprojectType: projecttype.Sketch,
			}

			projectdata.Initialize(testProject)

			result, output := ruleFunction()
			assert.Equal(t, testTable.expectedRuleResult, result, testTable.testName)
			assert.True(t, expectedOutputRegexp.MatchString(output), fmt.Sprintf("%s (output: %s, assertion regex: %s)", testTable.testName, output, testTable.expectedOutputQuery))
		})
	}
}

func TestSketchNameMismatch(t *testing.T) {
	testTables := []sketchRuleFunctionTestTable{
		{"Valid", "Valid", ruleresult.Pass, ""},
		{"Name Mismatch", "NameMismatch", ruleresult.Fail, ""},
		{"Case Mismatch", "CaseMismatch", ruleresult.Fail, ""},
	}

	checkSketchRuleFunction(SketchNameMismatch, testTables, t)
}

func TestProhibitedCharactersInSketchFileName(t *testing.T) {
	testTables := []sketchRuleFunctionTestTable{
		{"Has prohibited characters", "ProhibitedCharactersInFileName", ruleresult.Fail, "^Prohibited CharactersInFileName.h$"},
		{"No prohibited characters", "AllowedCharactersInFilenames", ruleresult.Pass, ""},
	}

	checkSketchRuleFunction(ProhibitedCharactersInSketchFileName, testTables, t)
}

func TestSketchFileNameGTMaxLength(t *testing.T) {
	testTables := []sketchRuleFunctionTestTable{
		{"Has file name > max length", "FileNameTooLong", ruleresult.Fail, "^FileNameTooLong12345678901234567890123456789012345678901234567890.h$"},
		{"File names <= max length", "Valid", ruleresult.Pass, ""},
	}

	checkSketchRuleFunction(SketchFileNameGTMaxLength, testTables, t)
}

func TestPdeSketchExtension(t *testing.T) {
	testTables := []sketchRuleFunctionTestTable{
		{"Has .pde", "Pde", ruleresult.Fail, ""},
		{"No .pde", "Valid", ruleresult.Pass, ""},
	}

	checkSketchRuleFunction(PdeSketchExtension, testTables, t)
}

func TestIncorrectSketchSrcFolderNameCase(t *testing.T) {
	testTables := []sketchRuleFunctionTestTable{
		{"Incorrect case", "IncorrectSrcFolderNameCase", ruleresult.Fail, ""},
		{"Correct case", "Valid", ruleresult.Pass, ""},
	}

	checkSketchRuleFunction(IncorrectSketchSrcFolderNameCase, testTables, t)
}
