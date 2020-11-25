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
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path
var schemasPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	testDataPath = paths.New(workingDirectory, "testdata", "libraries")
	schemasPath = paths.New(workingDirectory, "..", "..", "etc", "schemas")
}

type checkFunctionTestTable struct {
	testName            string
	libraryFolderName   string
	expectedCheckResult checkresult.Type
	expectedOutputQuery string
}

func checkCheckFunction(checkFunction Type, testTables []checkFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             testDataPath.Join(testTable.libraryFolderName),
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.Library,
		}

		checkdata.Initialize(testProject, schemasPath)

		result, output := checkFunction()
		assert.Equal(t, testTable.expectedCheckResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestLibraryPropertiesNameFieldMissingOfficialPrefix(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Correct prefix", "Arduino_Official", checkresult.Pass, ""},
		{"Incorrect prefix", "Recursive", checkresult.Fail, "^Recursive$"},
	}

	checkCheckFunction(LibraryPropertiesNameFieldMissingOfficialPrefix, testTables, t)
}

func TestLibraryPropertiesNameFieldDuplicate(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Duplicate", "Indexed", checkresult.Fail, "^Servo$"},
		{"Not duplicate", "NotIndexed", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesNameFieldDuplicate, testTables, t)
}

func TestLibraryPropertiesNameFieldNotInIndex(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"In index", "Indexed", checkresult.Pass, ""},
		{"Not in index", "NotIndexed", checkresult.Fail, "^NotIndexed$"},
	}

	checkCheckFunction(LibraryPropertiesNameFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesNameFieldHeaderMismatch(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Mismatch", "NameHeaderMismatch", checkresult.Fail, "^NameHeaderMismatch.h$"},
		{"Match", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesNameFieldHeaderMismatch, testTables, t)
}

func TestLibraryPropertiesSentenceFieldSpellCheck(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^grill broccoli now$"},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesSentenceFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldSpellCheck(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^There is a zebra$"},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesParagraphFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldRepeatsSentence(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Repeat", "ParagraphRepeatsSentence", checkresult.Fail, ""},
		{"No repeat", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesParagraphFieldRepeatsSentence, testTables, t)
}
func TestLibraryPropertiesUrlFieldDeadLink(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Bad URL", "BadURL", checkresult.Fail, "^Get \"http://invalid/\": dial tcp: lookup invalid:"},
		{"HTTP error 404", "URL404", checkresult.Fail, "^404 Not Found$"},
		{"Good URL", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesUrlFieldDeadLink, testTables, t)
}

func TestLibraryPropertiesDependsFieldNotInIndex(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Dependency not in index", "DependsNotIndexed", checkresult.Fail, "^NotIndexed$"},
		{"Dependency in index", "DependsIndexed", checkresult.Pass, ""},
		{"No depends", "NoDepends", checkresult.NotRun, ""},
	}

	checkCheckFunction(LibraryPropertiesDependsFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesDotALinkageFieldTrueWithFlatLayout(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Flat layout", "DotALinkageFlat", checkresult.Fail, ""},
		{"Recursive layout", "DotALinkage", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesDotALinkageFieldTrueWithFlatLayout, testTables, t)
}

func TestLibraryPropertiesIncludesFieldItemNotFound(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Missing includes", "MissingIncludes", checkresult.Fail, "^Nonexistent.h$"},
		{"Present includes", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryPropertiesIncludesFieldItemNotFound, testTables, t)
}

func TestLibraryPropertiesPrecompiledFieldEnabledWithFlatLayout(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Flat layout", "PrecompiledFlat", checkresult.Fail, "^true$"},
		{"Recursive layout", "Precompiled", checkresult.Pass, ""},
		{"Recursive, not precompiled", "NotPrecompiled", checkresult.NotRun, ""},
		{"Flat, not precompiled", "NotPrecompiledFlat", checkresult.NotRun, ""},
	}

	checkCheckFunction(LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout, testTables, t)
}

func TestLibraryHasSubmodule(t *testing.T) {
	testTables := []checkFunctionTestTable{
		{"Has submodule", "Submodule", checkresult.Fail, ""},
		{"No submodule", "Recursive", checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryHasSubmodule, testTables, t)
}

func TestLibraryContainsSymlinks(t *testing.T) {
	testLibrary := "Recursive"
	symlinkPath := testDataPath.Join(testLibrary, "test-symlink")
	// It's probably most friendly to developers using Windows to create the symlink needed for the test on demand.
	err := os.Symlink(testDataPath.Join(testLibrary, "library.properties").String(), symlinkPath.String())
	require.Nil(t, err, "This test must be run as administrator on Windows to have symlink creation privilege.")
	defer symlinkPath.RemoveAll() // clean up

	testTables := []checkFunctionTestTable{
		{"Has symlink", testLibrary, checkresult.Fail, ""},
	}

	checkCheckFunction(LibraryContainsSymlinks, testTables, t)

	err = symlinkPath.RemoveAll()
	require.Nil(t, err)

	testTables = []checkFunctionTestTable{
		{"No symlink", testLibrary, checkresult.Pass, ""},
	}

	checkCheckFunction(LibraryContainsSymlinks, testTables, t)
}
