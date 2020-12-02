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
	"regexp"
	"testing"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var platformTestDataPath *paths.Path

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	platformTestDataPath = workingDirectory.Join("testdata", "platforms")
}

type platformCheckFunctionTestTable struct {
	testName            string
	platformFolderName  string
	expectedCheckResult checkresult.Type
	expectedOutputQuery string
}

func checkPlatformCheckFunction(checkFunction Type, testTables []platformCheckFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             platformTestDataPath.Join(testTable.platformFolderName),
			ProjectType:      projecttype.Platform,
			SuperprojectType: projecttype.Platform,
		}

		checkdata.Initialize(testProject, nil)

		result, output := checkFunction()
		assert.Equal(t, testTable.expectedCheckResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestBoardsTxtMissing(t *testing.T) {
	testTables := []platformCheckFunctionTestTable{
		{"Present", "valid-boards.txt", checkresult.Pass, ""},
		{"Missing", "missing-boards.txt", checkresult.Fail, ""},
	}

	checkPlatformCheckFunction(BoardsTxtMissing, testTables, t)
}

func TestBoardsTxtFormat(t *testing.T) {
	testTables := []platformCheckFunctionTestTable{
		{"Missing", "missing-boards.txt", checkresult.NotRun, ""},
		{"Valid", "valid-boards.txt", checkresult.Pass, ""},
		{"Invalid", "invalid-boards.txt", checkresult.Fail, ""},
	}

	checkPlatformCheckFunction(BoardsTxtFormat, testTables, t)
}
