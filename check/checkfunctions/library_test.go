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

type libraryCheckFunctionTestTable struct {
	testName            string
	libraryFolderName   string
	expectedCheckResult checkresult.Type
	expectedOutputQuery string
}

func checkLibraryCheckFunction(checkFunction Type, testTables []libraryCheckFunctionTestTable, t *testing.T) {
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

func TestMisspelledLibraryPropertiesFileName(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Incorrect", "MisspelledLibraryProperties", checkresult.Fail, ""},
		{"Correct", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(MisspelledLibraryPropertiesFileName, testTables, t)
}

func TestIncorrectLibraryPropertiesFileNameCase(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Incorrect", "IncorrectLibraryPropertiesCase", checkresult.Fail, ""},
		{"Correct", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(IncorrectLibraryPropertiesFileNameCase, testTables, t)
}

func TestLibraryPropertiesMissing(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Invalid non-legacy", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Legacy", "Legacy", checkresult.Fail, ""},
		{"Flat non-legacy", "Flat", checkresult.Pass, ""},
		{"Recursive", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesMissing, testTables, t)
}

func TestRedundantLibraryProperties(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Redundant", "RedundantLibraryProperties", checkresult.Fail, ""},
		{"No redundant", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(RedundantLibraryProperties, testTables, t)
}

func TestLibraryPropertiesFormat(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", checkresult.Fail, ""},
		{"Legacy", "Legacy", checkresult.NotRun, ""},
		{"Valid", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesFormat, testTables, t)
}

func TestLibraryPropertiesNameFieldMissingOfficialPrefix(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Correct prefix", "Arduino_Official", checkresult.Pass, ""},
		{"Incorrect prefix", "Recursive", checkresult.Fail, "^Recursive$"},
	}

	checkLibraryCheckFunction(LibraryPropertiesNameFieldMissingOfficialPrefix, testTables, t)
}

func TestLibraryPropertiesNameFieldDuplicate(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Duplicate", "Indexed", checkresult.Fail, "^Servo$"},
		{"Not duplicate", "NotIndexed", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesNameFieldDuplicate, testTables, t)
}

func TestLibraryPropertiesNameFieldNotInIndex(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"In index", "Indexed", checkresult.Pass, ""},
		{"Not in index", "NotIndexed", checkresult.Fail, "^NotIndexed$"},
	}

	checkLibraryCheckFunction(LibraryPropertiesNameFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesNameFieldHeaderMismatch(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Mismatch", "NameHeaderMismatch", checkresult.Fail, "^NameHeaderMismatch.h$"},
		{"Match", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesNameFieldHeaderMismatch, testTables, t)
}

func TestLibraryPropertiesSentenceFieldSpellCheck(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^grill broccoli now$"},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesSentenceFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldSpellCheck(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^There is a zebra$"},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesParagraphFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldRepeatsSentence(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Repeat", "ParagraphRepeatsSentence", checkresult.Fail, ""},
		{"No repeat", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesParagraphFieldRepeatsSentence, testTables, t)
}
func TestLibraryPropertiesUrlFieldDeadLink(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Bad URL", "BadURL", checkresult.Fail, "^Get \"http://invalid/\": dial tcp: lookup invalid:"},
		{"HTTP error 404", "URL404", checkresult.Fail, "^404 Not Found$"},
		{"Good URL", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesUrlFieldDeadLink, testTables, t)
}

func TestLibraryPropertiesDependsFieldNotInIndex(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Dependency not in index", "DependsNotIndexed", checkresult.Fail, "^NotIndexed$"},
		{"Dependency in index", "DependsIndexed", checkresult.Pass, ""},
		{"No depends", "NoDepends", checkresult.NotRun, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesDependsFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesDotALinkageFieldTrueWithFlatLayout(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Flat layout", "DotALinkageFlat", checkresult.Fail, ""},
		{"Recursive layout", "DotALinkage", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesDotALinkageFieldTrueWithFlatLayout, testTables, t)
}

func TestLibraryPropertiesIncludesFieldItemNotFound(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Missing includes", "MissingIncludes", checkresult.Fail, "^Nonexistent.h$"},
		{"Present includes", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesIncludesFieldItemNotFound, testTables, t)
}

func TestLibraryPropertiesPrecompiledFieldEnabledWithFlatLayout(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.NotRun, ""},
		{"Flat layout", "PrecompiledFlat", checkresult.Fail, "^true$"},
		{"Recursive layout", "Precompiled", checkresult.Pass, ""},
		{"Recursive, not precompiled", "NotPrecompiled", checkresult.NotRun, ""},
		{"Flat, not precompiled", "Flat", checkresult.NotRun, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout, testTables, t)
}

func TestLibraryInvalid(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Invalid library.properties", "InvalidLibraryProperties", checkresult.Fail, ""},
		{"Invalid flat layout", "FlatWithoutHeader", checkresult.Fail, ""},
		{"Invalid recursive layout", "RecursiveWithoutLibraryProperties", checkresult.Fail, ""},
		{"Valid library", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryInvalid, testTables, t)
}

func TestLibraryHasSubmodule(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has submodule", "Submodule", checkresult.Fail, ""},
		{"No submodule", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryHasSubmodule, testTables, t)
}

func TestLibraryContainsSymlinks(t *testing.T) {
	testLibrary := "Recursive"
	symlinkPath := testDataPath.Join(testLibrary, "test-symlink")
	// It's probably most friendly to developers using Windows to create the symlink needed for the test on demand.
	err := os.Symlink(testDataPath.Join(testLibrary, "library.properties").String(), symlinkPath.String())
	require.Nil(t, err, "This test must be run as administrator on Windows to have symlink creation privilege.")
	defer symlinkPath.RemoveAll() // clean up

	testTables := []libraryCheckFunctionTestTable{
		{"Has symlink", testLibrary, checkresult.Fail, ""},
	}

	checkLibraryCheckFunction(LibraryContainsSymlinks, testTables, t)

	err = symlinkPath.RemoveAll()
	require.Nil(t, err)

	testTables = []libraryCheckFunctionTestTable{
		{"No symlink", testLibrary, checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryContainsSymlinks, testTables, t)
}

func TestLibraryHasDotDevelopmentFile(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has .development file", "DotDevelopment", checkresult.Fail, ""},
		{"No .development file", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryHasDotDevelopmentFile, testTables, t)
}

func TestLibraryHasExe(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has .exe file", "Exe", checkresult.Fail, ""},
		{"No .exe files", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryHasExe, testTables, t)
}

func TestLibraryHasStraySketches(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Sketch in root", "SketchInRoot", checkresult.Fail, ""},
		{"Sketch in subfolder", "MisspelledExamplesFolder", checkresult.Fail, ""},
		{"Sketch in legit location", "ExamplesFolder", checkresult.Pass, ""},
		{"No sketches", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryHasStraySketches, testTables, t)
}

func TestProhibitedCharactersInLibraryFolderName(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has prohibited characters", "Prohibited CharactersInFolderName", checkresult.Fail, ""},
		{"No prohibited characters", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(ProhibitedCharactersInLibraryFolderName, testTables, t)
}

func TestLibraryFolderNameGTMaxLength(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has folder name > max length", "FolderNameTooLong12345678901234567890123456789012345678901234567890", checkresult.Fail, "^FolderNameTooLong12345678901234567890123456789012345678901234567890$"},
		{"Folder name <= max length", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryFolderNameGTMaxLength, testTables, t)
}

func TestIncorrectLibrarySrcFolderNameCase(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Flat, not precompiled", "Flat", checkresult.NotRun, ""},
		{"Incorrect case", "IncorrectSrcFolderNameCase", checkresult.Fail, ""},
		{"Correct case", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(IncorrectLibrarySrcFolderNameCase, testTables, t)
}

func TestMisspelledExamplesFolderName(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correctly spelled", "ExamplesFolder", checkresult.Pass, ""},
		{"Misspelled", "MisspelledExamplesFolder", checkresult.Fail, ""},
		{"No examples folder", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(MisspelledExamplesFolderName, testTables, t)
}

func TestIncorrectExamplesFolderNameCase(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correct case", "ExamplesFolder", checkresult.Pass, ""},
		{"Incorrect case", "IncorrectExamplesFolderCase", checkresult.Fail, ""},
		{"No examples folder", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(IncorrectExamplesFolderNameCase, testTables, t)
}

func TestMisspelledExtrasFolderName(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correctly spelled", "ExtrasFolder", checkresult.Pass, ""},
		{"Misspelled", "MisspelledExtrasFolder", checkresult.Fail, ""},
		{"No extras folder", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(MisspelledExtrasFolderName, testTables, t)
}

func TestIncorrectExtrasFolderNameCase(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correct case", "ExtrasFolder", checkresult.Pass, ""},
		{"Incorrect case", "IncorrectExtrasFolderCase", checkresult.Fail, ""},
		{"No extras folder", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(IncorrectExtrasFolderNameCase, testTables, t)
}

func TestRecursiveLibraryWithUtilityFolder(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Flat", "Flat", checkresult.NotRun, ""},
		{"Recursive with utility", "RecursiveWithUtilityFolder", checkresult.Fail, ""},
		{"Recursive without utility", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(RecursiveLibraryWithUtilityFolder, testTables, t)
}

func TestMissingReadme(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Readme", "Readme", checkresult.Pass, ""},
		{"No readme", "NoReadme", checkresult.Fail, ""},
	}

	checkLibraryCheckFunction(MissingReadme, testTables, t)
}
