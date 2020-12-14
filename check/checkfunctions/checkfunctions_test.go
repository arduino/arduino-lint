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

package checkfunctions

import (
	"os"
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/check/checkdata"
	"github.com/arduino/arduino-lint/check/checkresult"
	"github.com/arduino/arduino-lint/project"
	"github.com/arduino/arduino-lint/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	testDataPath = paths.New(workingDirectory, "testdata", "general")
}

type checkFunctionTestTable struct {
	testName            string
	projectFolderName   string
	projectType         projecttype.Type
	superProjectType    projecttype.Type
	expectedCheckResult checkresult.Type
	expectedOutputQuery string
}

func checkCheckFunction(checkFunction Type, testTables []checkFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             testDataPath.Join(testTable.projectFolderName),
			ProjectType:      testTable.projectType,
			SuperprojectType: testTable.superProjectType,
		}

		checkdata.Initialize(testProject)

		result, output := checkFunction()
		assert.Equal(t, testTable.expectedCheckResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestMissingReadme(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Subproject", "readme", projecttype.Sketch, projecttype.Library, checkresult.Skip, ""},
		{"Readme", "readme", projecttype.Sketch, projecttype.Sketch, checkresult.Pass, ""},
		{"No readme", "no-readme", projecttype.Sketch, projecttype.Sketch, checkresult.Fail, ""},
	}

	checkCheckFunction(MissingReadme, testTables, t)
}

func TestMissingLicenseFile(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Subproject", "no-license-file", projecttype.Sketch, projecttype.Library, checkresult.Skip, ""},
		{"Has license", "license-file", projecttype.Sketch, projecttype.Sketch, checkresult.Pass, ""},
		{"Has license in subfolder", "license-file-in-subfolder", projecttype.Sketch, projecttype.Sketch, checkresult.Fail, ""},
		{"No license", "no-license-file", projecttype.Sketch, projecttype.Sketch, checkresult.Fail, ""},
	}

	checkCheckFunction(MissingLicenseFile, testTables, t)
}

func TestIncorrectArduinoDotHFileNameCase(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Incorrect, angle brackets", "arduino.h-angle", projecttype.Sketch, projecttype.Sketch, checkresult.Fail, ""},
		{"Incorrect, quotes", "arduino.h-quote", projecttype.Sketch, projecttype.Sketch, checkresult.Fail, ""},
		{"Correct case", "Arduino.h", projecttype.Sketch, projecttype.Sketch, checkresult.Pass, ""},
	}

	checkCheckFunction(IncorrectArduinoDotHFileNameCase, testTables, t)
}
