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
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/internal/check/checkresult"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/checkdata"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var packageIndexesTestDataPath *paths.Path

func init() {
	workingDirectory, _ := paths.Getwd()
	packageIndexesTestDataPath = workingDirectory.Join("testdata", "packageindexes")
}

type packageIndexCheckFunctionTestTable struct {
	testName               string
	packageIndexFolderName string
	expectedCheckResult    checkresult.Type
	expectedOutputQuery    string
}

func checkPackageIndexCheckFunction(checkFunction Type, testTables []packageIndexCheckFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             packageIndexesTestDataPath.Join(testTable.packageIndexFolderName),
			ProjectType:      projecttype.PackageIndex,
			SuperprojectType: projecttype.PackageIndex,
		}

		checkdata.Initialize(testProject)

		result, output := checkFunction()
		assert.Equal(t, testTable.expectedCheckResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestPackageIndexJSONFormat(t *testing.T) {
	testTables := []packageIndexCheckFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", checkresult.Fail, ""},
		{"Not valid package index", "invalid-package-index", checkresult.Pass, ""},
		{"Valid package index", "valid-package-index", checkresult.Pass, ""},
	}

	checkPackageIndexCheckFunction(PackageIndexJSONFormat, testTables, t)
}

func TestPackageIndexFormat(t *testing.T) {
	testTables := []packageIndexCheckFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", checkresult.Fail, ""},
		{"Not valid package index", "invalid-package-index", checkresult.Fail, ""},
		{"Valid package index", "valid-package-index", checkresult.Pass, ""},
	}

	checkPackageIndexCheckFunction(PackageIndexFormat, testTables, t)
}
