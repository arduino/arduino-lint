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
	"time"

	"github.com/arduino/arduino-lint/check/checkdata"
	"github.com/arduino/arduino-lint/check/checkresult"
	"github.com/arduino/arduino-lint/project"
	"github.com/arduino/arduino-lint/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var librariesTestDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	librariesTestDataPath = paths.New(workingDirectory, "testdata", "libraries")
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
			Path:             librariesTestDataPath.Join(testTable.libraryFolderName),
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.Library,
		}

		checkdata.Initialize(testProject)

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
		{"Legacy", "Legacy", checkresult.Skip, ""},
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

func TestLibraryPropertiesVersionFieldBehindTag(t *testing.T) {
	// Set up the test repository folders.
	TagPrereleaseGreaterPath := librariesTestDataPath.Join("TagPrereleaseGreater")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagPrereleaseGreaterPath))
	defer TagPrereleaseGreaterPath.RemoveAll()

	TagGreaterPath := librariesTestDataPath.Join("TagGreater")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagGreaterPath))
	defer TagGreaterPath.RemoveAll()

	LightweightTagGreaterPath := librariesTestDataPath.Join("LightweightTagGreater")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(LightweightTagGreaterPath))
	defer LightweightTagGreaterPath.RemoveAll()

	TagMatchPath := librariesTestDataPath.Join("TagMatch")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagMatchPath))
	defer TagMatchPath.RemoveAll()

	LightweightTagMatchPath := librariesTestDataPath.Join("LightweightTagMatch")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(LightweightTagMatchPath))
	defer LightweightTagMatchPath.RemoveAll()

	TagMatchWithPrefixPath := librariesTestDataPath.Join("TagMatchWithPrefix")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagMatchWithPrefixPath))
	defer TagMatchWithPrefixPath.RemoveAll()

	TagLessThanPath := librariesTestDataPath.Join("TagLessThan")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagLessThanPath))
	defer TagLessThanPath.RemoveAll()

	TagNotVersionPath := librariesTestDataPath.Join("TagNotVersion")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(TagNotVersionPath))
	defer TagNotVersionPath.RemoveAll()

	NoTagsPath := librariesTestDataPath.Join("NoTags")
	require.Nil(t, librariesTestDataPath.Join("Recursive").CopyDirTo(NoTagsPath))
	defer NoTagsPath.RemoveAll()

	// Test repositories are generated on the fly.
	gitInitAndTag := func(t *testing.T, repositoryPath *paths.Path, tagName string, annotated bool) string {
		repository, err := git.PlainInit(repositoryPath.String(), false)
		require.Nil(t, err)
		worktree, err := repository.Worktree()
		require.Nil(t, err)
		_, err = worktree.Add(".")
		require.Nil(t, err)

		signature := &object.Signature{
			Name:  "Jane Developer",
			Email: "janedeveloper@example.com",
			When:  time.Now(),
		}

		_, err = worktree.Commit(
			"Test commit message",
			&git.CommitOptions{
				Author: signature,
			},
		)
		require.Nil(t, err)

		headRef, err := repository.Head()
		require.Nil(t, err)

		if tagName != "" {
			// Annotated and lightweight tags are significantly different, so it's important to ensure the check code works correctly with both.
			if annotated {
				_, err = repository.CreateTag(
					tagName,
					headRef.Hash(),
					&git.CreateTagOptions{
						Tagger:  signature,
						Message: tagName,
					},
				)
			} else {
				_, err = repository.CreateTag(tagName, headRef.Hash(), nil)
			}
			require.Nil(t, err)
		}

		return repositoryPath.Base()
	}

	testTables := []libraryCheckFunctionTestTable{
		// TODO: Test Skip if subproject
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Legacy", "Legacy", checkresult.NotRun, ""},
		{"Unparsable version", "VersionFormatInvalid", checkresult.NotRun, ""},
		{"Not repo", "Recursive", checkresult.Skip, ""},
		{"Tag name not a version", gitInitAndTag(t, TagNotVersionPath, "foo", true), checkresult.Pass, ""},
		{"Match w/ tag prefix", gitInitAndTag(t, TagMatchWithPrefixPath, "1.0.0", true), checkresult.Pass, ""},
		{"Pre-release tag greater", gitInitAndTag(t, TagPrereleaseGreaterPath, "1.0.1-rc1", true), checkresult.Pass, ""},
		{"Tag greater", gitInitAndTag(t, TagGreaterPath, "1.0.1", true), checkresult.Fail, ""},
		{"Lightweight tag greater", gitInitAndTag(t, LightweightTagGreaterPath, "1.0.1", false), checkresult.Fail, ""},
		{"Tag matches", gitInitAndTag(t, TagMatchPath, "1.0.0", true), checkresult.Pass, ""},
		{"Lightweight tag matches", gitInitAndTag(t, LightweightTagMatchPath, "1.0.0", false), checkresult.Pass, ""},
		{"Tag less than version", gitInitAndTag(t, TagLessThanPath, "0.1.0", true), checkresult.Pass, ""},
		{"No tags", gitInitAndTag(t, NoTagsPath, "", true), checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesVersionFieldBehindTag, testTables, t)
}

func TestLibraryPropertiesSentenceFieldSpellCheck(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.Skip, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^grill broccoli now$"},
		{"Non-nil diff but no typos", "SpuriousMisspelledSentenceParagraphValue", checkresult.Pass, ""},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesSentenceFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldSpellCheck(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.Skip, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", checkresult.Fail, "^There is a zebra$"},
		{"Non-nil diff but no typos", "SpuriousMisspelledSentenceParagraphValue", checkresult.Pass, ""},
		{"Correct spelling", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesParagraphFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesEmailFieldAsMaintainerAlias(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"No email field", "MissingFields", checkresult.Skip, ""},
		{"email in place of maintainer", "EmailOnly", checkresult.Fail, ""},
		{"email and maintainer", "EmailAndMaintainer", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesEmailFieldAsMaintainerAlias, testTables, t)
}

func TestLibraryPropertiesParagraphFieldRepeatsSentence(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Repeat", "ParagraphRepeatsSentence", checkresult.Fail, ""},
		{"No repeat", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesParagraphFieldRepeatsSentence, testTables, t)
}

func TestLibraryPropertiesCategoryFieldUncategorized(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"No category field", "MissingFields", checkresult.Skip, ""},
		{"Uncategorized category", "UncategorizedCategoryValue", checkresult.Fail, ""},
		{"Valid category value", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesCategoryFieldUncategorized, testTables, t)
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
		{"Depends field empty", "DependsEmpty", checkresult.Pass, ""},
		{"No depends", "NoDepends", checkresult.Skip, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesDependsFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesDotALinkageFieldTrueWithFlatLayout(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.Skip, ""},
		{"Flat layout", "DotALinkageFlat", checkresult.Fail, ""},
		{"Recursive layout", "DotALinkage", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesDotALinkageFieldTrueWithFlatLayout, testTables, t)
}

func TestLibraryPropertiesIncludesFieldItemNotFound(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.Skip, ""},
		{"Missing includes", "MissingIncludes", checkresult.Fail, "^Nonexistent.h$"},
		{"Present includes", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(LibraryPropertiesIncludesFieldItemNotFound, testTables, t)
}

func TestLibraryPropertiesPrecompiledFieldEnabledWithFlatLayout(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", checkresult.NotRun, ""},
		{"Not defined", "MissingFields", checkresult.Skip, ""},
		{"Flat layout", "PrecompiledFlat", checkresult.Fail, "^true$"},
		{"Recursive layout", "Precompiled", checkresult.Pass, ""},
		{"Recursive, not precompiled", "NotPrecompiled", checkresult.Skip, ""},
		{"Flat, not precompiled", "Flat", checkresult.Skip, ""},
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
	symlinkPath := librariesTestDataPath.Join(testLibrary, "test-symlink")
	// It's probably most friendly to developers using Windows to create the symlink needed for the test on demand.
	err := os.Symlink(librariesTestDataPath.Join(testLibrary, "library.properties").String(), symlinkPath.String())
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
		{"Flat, not precompiled", "Flat", checkresult.Skip, ""},
		{"Incorrect case", "IncorrectSrcFolderNameCase", checkresult.Fail, ""},
		{"Correct case", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(IncorrectLibrarySrcFolderNameCase, testTables, t)
}

func TestMissingExamples(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Has examples", "ExamplesFolder", checkresult.Pass, ""},
		{`Has examples (in "example" folder)`, "ExampleFolder", checkresult.Pass, ""},
		{"No examples", "NoExamples", checkresult.Fail, ""},
	}

	checkLibraryCheckFunction(MissingExamples, testTables, t)
}

func TestMisspelledExamplesFolderName(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correctly spelled", "ExamplesFolder", checkresult.Pass, ""},
		{"Misspelled", "MisspelledExamplesFolder", checkresult.Fail, ""},
		{"No examples folder", "NoExamples", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(MisspelledExamplesFolderName, testTables, t)
}

func TestIncorrectExamplesFolderNameCase(t *testing.T) {
	testTables := []libraryCheckFunctionTestTable{
		{"Correct case", "ExamplesFolder", checkresult.Pass, ""},
		{"Incorrect case", "IncorrectExamplesFolderCase", checkresult.Fail, ""},
		{"No examples folder", "NoExamples", checkresult.Pass, ""},
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
		{"Flat", "Flat", checkresult.Skip, ""},
		{"Recursive with utility", "RecursiveWithUtilityFolder", checkresult.Fail, ""},
		{"Recursive without utility", "Recursive", checkresult.Pass, ""},
	}

	checkLibraryCheckFunction(RecursiveLibraryWithUtilityFolder, testTables, t)
}
