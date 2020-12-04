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

package checkfunctions

import (
	"os"
	"regexp"
	"testing"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var sketchesTestDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	sketchesTestDataPath = paths.New(workingDirectory, "testdata", "sketches")
}

type sketchCheckFunctionTestTable struct {
	testName            string
	sketchFolderName    string
	expectedCheckResult checkresult.Type
	expectedOutputQuery string
}

func checkSketchCheckFunction(checkFunction Type, testTables []sketchCheckFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             sketchesTestDataPath.Join(testTable.sketchFolderName),
			ProjectType:      projecttype.Sketch,
			SuperprojectType: projecttype.Sketch,
		}

		checkdata.Initialize(testProject, schemasPath)

		result, output := checkFunction()
		assert.Equal(t, testTable.expectedCheckResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestIncorrectSketchSrcFolderNameCase(t *testing.T) {
	testTables := []sketchCheckFunctionTestTable{
		{"Incorrect case", "IncorrectSrcFolderNameCase", checkresult.Fail, ""},
		{"Correct case", "Valid", checkresult.Pass, ""},
	}

	checkSketchCheckFunction(IncorrectSketchSrcFolderNameCase, testTables, t)
}

func TestProhibitedCharactersInSketchFileName(t *testing.T) {
	testTables := []sketchCheckFunctionTestTable{
		{"Has prohibited characters", "ProhibitedCharactersInFileName", checkresult.Fail, "^Prohibited CharactersInFileName.h$"},
		{"No prohibited characters", "Valid", checkresult.Pass, ""},
	}

	checkSketchCheckFunction(ProhibitedCharactersInSketchFileName, testTables, t)
}

func TestSketchFileNameGTMaxLength(t *testing.T) {
	testTables := []sketchCheckFunctionTestTable{
		{"Has file name > max length", "FileNameTooLong", checkresult.Fail, "^FileNameTooLong12345678901234567890123456789012345678901234567890.h$"},
		{"File names <= max length", "Valid", checkresult.Pass, ""},
	}

	checkSketchCheckFunction(SketchFileNameGTMaxLength, testTables, t)
}

func TestPdeSketchExtension(t *testing.T) {
	testTables := []sketchCheckFunctionTestTable{
		{"Has .pde", "Pde", checkresult.Fail, ""},
		{"No .pde", "Valid", checkresult.Pass, ""},
	}

	checkSketchCheckFunction(PdeSketchExtension, testTables, t)
}

func TestSketchDotJSONJSONFormat(t *testing.T) {
	testTables := []sketchCheckFunctionTestTable{
		{"No metadata file", "NoMetadataFile", checkresult.NotRun, ""},
		{"Valid", "ValidMetadataFile", checkresult.Pass, ""},
		{"Invalid", "InvalidJSONMetadataFile", checkresult.Fail, ""},
	}

	checkSketchCheckFunction(SketchDotJSONJSONFormat, testTables, t)
}
