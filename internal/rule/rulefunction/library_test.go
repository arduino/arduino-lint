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
	"time"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
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

type libraryRuleFunctionTestTable struct {
	testName            string
	libraryFolderName   string
	expectedRuleResult  ruleresult.Type
	expectedOutputQuery string
}

func checkLibraryRuleFunction(ruleFunction Type, testTables []libraryRuleFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             librariesTestDataPath.Join(testTable.libraryFolderName),
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.Library,
		}

		projectdata.Initialize(testProject)

		result, output := ruleFunction()
		assert.Equal(t, testTable.expectedRuleResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), fmt.Sprintf("%s (output: %s, assertion regex: %s)", testTable.testName, output, testTable.expectedOutputQuery))
	}
}

func TestLibraryInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid library.properties", "InvalidLibraryProperties", ruleresult.Fail, ""},
		{"Invalid flat layout", "FlatWithoutHeader", ruleresult.Fail, ""},
		{"Invalid recursive layout", "RecursiveWithoutLibraryProperties", ruleresult.Fail, ""},
		{"Valid library", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryInvalid, testTables, t)
}

func TestLibraryFolderNameGTMaxLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Has folder name > max length", "FolderNameTooLong12345678901234567890123456789012345678901234567890", ruleresult.Fail, "^FolderNameTooLong12345678901234567890123456789012345678901234567890$"},
		{"Folder name <= max length", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryFolderNameGTMaxLength, testTables, t)
}

func TestProhibitedCharactersInLibraryFolderName(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Has prohibited characters", "Prohibited CharactersInFolderName", ruleresult.Fail, ""},
		{"No prohibited characters", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(ProhibitedCharactersInLibraryFolderName, testTables, t)
}

func TestLibraryHasSubmodule(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Has submodule", "Submodule", ruleresult.Fail, ""},
		{"No submodule", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryHasSubmodule, testTables, t)
}

func TestLibraryContainsSymlinks(t *testing.T) {
	testLibrary := "Recursive"
	// Set up a library with a file target symlink.
	symlinkPath := librariesTestDataPath.Join(testLibrary, "test-symlink")
	// It's probably most friendly to developers using Windows to create the symlink needed for the test on demand.
	err := os.Symlink(librariesTestDataPath.Join(testLibrary, "library.properties").String(), symlinkPath.String())
	require.Nil(t, err, "This test must be run as administrator on Windows to have symlink creation privilege.")
	defer symlinkPath.RemoveAll() // clean up

	testTables := []libraryRuleFunctionTestTable{
		{"Has file target symlink", testLibrary, ruleresult.Fail, ""},
	}

	checkLibraryRuleFunction(LibraryContainsSymlinks, testTables, t)

	err = symlinkPath.RemoveAll()
	require.Nil(t, err)

	// Set up a library with a folder target symlink.
	err = os.Symlink(librariesTestDataPath.Join(testLibrary, "src").String(), symlinkPath.String())
	require.Nil(t, err)

	testTables = []libraryRuleFunctionTestTable{
		{"Has folder target symlink", testLibrary, ruleresult.Fail, ""},
	}

	checkLibraryRuleFunction(LibraryContainsSymlinks, testTables, t)

	err = symlinkPath.RemoveAll()
	require.Nil(t, err)

	testTables = []libraryRuleFunctionTestTable{
		{"No symlink", testLibrary, ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryContainsSymlinks, testTables, t)
}

func TestLibraryHasDotDevelopmentFile(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Has .development file", "DotDevelopment", ruleresult.Fail, ""},
		{"No .development file", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryHasDotDevelopmentFile, testTables, t)
}

func TestLibraryHasExe(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Has .exe file", "Exe", ruleresult.Fail, ""},
		{"No .exe files", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryHasExe, testTables, t)
}

func TestLibraryPropertiesNameFieldHeaderMismatch(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Mismatch", "NameHeaderMismatch", ruleresult.Fail, "^NameHeaderMismatch.h$"},
		{"Match", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldHeaderMismatch, testTables, t)
}

func TestIncorrectLibrarySrcFolderNameCase(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Flat, not precompiled", "Flat", ruleresult.Skip, ""},
		{"Incorrect case", "IncorrectSrcFolderNameCase", ruleresult.Fail, ""},
		{"Correct case", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(IncorrectLibrarySrcFolderNameCase, testTables, t)
}

func TestRecursiveLibraryWithUtilityFolder(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Flat", "Flat", ruleresult.Skip, ""},
		{"Recursive with utility", "RecursiveWithUtilityFolder", ruleresult.Fail, ""},
		{"Recursive without utility", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(RecursiveLibraryWithUtilityFolder, testTables, t)
}

func TestMisspelledExtrasFolderName(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Correctly spelled", "ExtrasFolder", ruleresult.Pass, ""},
		{"Misspelled", "MisspelledExtrasFolder", ruleresult.Fail, ""},
		{"No extras folder", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(MisspelledExtrasFolderName, testTables, t)
}

func TestIncorrectExtrasFolderNameCase(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Correct case", "ExtrasFolder", ruleresult.Pass, ""},
		{"Incorrect case", "IncorrectExtrasFolderCase", ruleresult.Fail, ""},
		{"No extras folder", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(IncorrectExtrasFolderNameCase, testTables, t)
}

func TestLibraryPropertiesMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid non-legacy", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Fail, ""},
		{"Flat non-legacy", "Flat", ruleresult.Pass, ""},
		{"Recursive", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMissing, testTables, t)
}

func TestMisspelledLibraryPropertiesFileName(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Incorrect", "MisspelledLibraryProperties", ruleresult.Fail, ""},
		{"Correct", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(MisspelledLibraryPropertiesFileName, testTables, t)
}

func TestIncorrectLibraryPropertiesFileNameCase(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Incorrect", "IncorrectLibraryPropertiesCase", ruleresult.Fail, ""},
		{"Correct", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(IncorrectLibraryPropertiesFileNameCase, testTables, t)
}

func TestRedundantLibraryProperties(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Redundant", "RedundantLibraryProperties", ruleresult.Fail, ""},
		{"No redundant", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(RedundantLibraryProperties, testTables, t)
}

func TestLibraryPropertiesFormat(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.Fail, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesFormat, testTables, t)
}

func TestLibraryPropertiesNameFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldMissing, testTables, t)
}

func TestLibraryPropertiesNameFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field too short", "NameLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesNameFieldGTMaxLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field too long", "NameGTMaxLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldGTMaxLength, testTables, t)
}

func TestLibraryPropertiesNameFieldGTRecommendedLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field longer than recommended", "NameIsBiggerThanRecommendedLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldGTRecommendedLength, testTables, t)
}

func TestLibraryPropertiesNameFieldDisallowedCharacters(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field has disallowed characters", "NameHasBadChars", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldDisallowedCharacters, testTables, t)
}

func TestLibraryPropertiesNameFieldStartsWithArduino(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field starts with Arduino", "Arduino_Official", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldStartsWithArduino, testTables, t)
}

func TestLibraryPropertiesNameFieldMissingOfficialPrefix(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.NotRun, ""},
		{"Correct prefix", "Arduino_Official", ruleresult.Pass, ""},
		{"Incorrect prefix", "Recursive", ruleresult.Fail, "^Recursive$"},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldMissingOfficialPrefix, testTables, t)
}

func TestLibraryPropertiesNameFieldContainsArduino(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field contains Arduino", "NameContainsArduino", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldContainsArduino, testTables, t)
}

func TestLibraryPropertiesNameFieldHasSpaces(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field contains spaces", "NameHasSpaces", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldHasSpaces, testTables, t)
}

func TestLibraryPropertiesNameFieldContainsLibrary(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Name field contains library", "NameHasLibrary", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldContainsLibrary, testTables, t)
}

func TestLibraryPropertiesNameFieldDuplicate(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Duplicate", "Indexed", ruleresult.Fail, "^Servo$"},
		{"Not duplicate", "NotIndexed", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldDuplicate, testTables, t)
}

func TestLibraryPropertiesNameFieldNotInIndex(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"In index", "Indexed", ruleresult.Pass, ""},
		{"Not in index", "NotIndexed", ruleresult.Fail, "^NotIndexed$"},
	}

	checkLibraryRuleFunction(LibraryPropertiesNameFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesVersionFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Version field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesVersionFieldMissing, testTables, t)
}

func TestLibraryPropertiesVersionFieldNonRelaxedSemver(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Version not relaxed semver compliant", "VersionNotRelaxedSemver", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesVersionFieldNonRelaxedSemver, testTables, t)
}

func TestLibraryPropertiesVersionFieldNonSemver(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Version not semver compliant", "VersionNotSemver", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesVersionFieldNonSemver, testTables, t)
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
			// Annotated and lightweight tags are significantly different, so it's important to ensure the rule code works correctly with both.
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

	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Unparsable version", "VersionFormatInvalid", ruleresult.NotRun, ""},
		{"Not repo", "Recursive", ruleresult.Skip, ""},
		{"Tag name not a version", gitInitAndTag(t, TagNotVersionPath, "foo", true), ruleresult.Pass, ""},
		{"Match w/ tag prefix", gitInitAndTag(t, TagMatchWithPrefixPath, "1.0.0", true), ruleresult.Pass, ""},
		{"Pre-release tag greater", gitInitAndTag(t, TagPrereleaseGreaterPath, "1.0.1-rc1", true), ruleresult.Pass, ""},
		{"Tag greater", gitInitAndTag(t, TagGreaterPath, "1.0.1", true), ruleresult.Fail, ""},
		{"Lightweight tag greater", gitInitAndTag(t, LightweightTagGreaterPath, "1.0.1", false), ruleresult.Fail, ""},
		{"Tag matches", gitInitAndTag(t, TagMatchPath, "1.0.0", true), ruleresult.Pass, ""},
		{"Lightweight tag matches", gitInitAndTag(t, LightweightTagMatchPath, "1.0.0", false), ruleresult.Pass, ""},
		{"Tag less than version", gitInitAndTag(t, TagLessThanPath, "0.1.0", true), ruleresult.Pass, ""},
		{"No tags", gitInitAndTag(t, NoTagsPath, "", true), ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesVersionFieldBehindTag, testTables, t)
}

func TestLibraryPropertiesAuthorFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesAuthorFieldMissing, testTables, t)
}

func TestLibraryPropertiesAuthorFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Author field too short", "AuthorLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesAuthorFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesMaintainerFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMaintainerFieldMissing, testTables, t)
}

func TestLibraryPropertiesMaintainerFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Maintainer field too short", "MaintainerLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMaintainerFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesMaintainerFieldStartsWithArduino(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Maintainer field starts w/ Arduino", "MaintainerStartsWithArduino", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMaintainerFieldStartsWithArduino, testTables, t)
}

func TestLibraryPropertiesMaintainerFieldContainsArduino(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Maintainer field contains Arduino", "MaintainerContainsArduino", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMaintainerFieldContainsArduino, testTables, t)
}

func TestLibraryPropertiesEmailFieldAsMaintainerAlias(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"No email field", "MissingFields", ruleresult.Skip, ""},
		{"email in place of maintainer", "EmailOnly", ruleresult.Fail, ""},
		{"email and maintainer", "EmailAndMaintainer", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesEmailFieldAsMaintainerAlias, testTables, t)
}

func TestLibraryPropertiesEmailFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Email field too short", "EmailLTMinLength", ruleresult.Fail, ""},
		{"Valid", "EmailOnly", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesEmailFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesEmailFieldStartsWithArduino(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not an alias", "EmailAndMaintainer", ruleresult.Skip, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Email field starts w/ Arduino", "EmailStartsWithArduino", ruleresult.Fail, ""},
		{"Valid", "EmailOnly", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesEmailFieldStartsWithArduino, testTables, t)
}

func TestLibraryPropertiesSentenceFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesSentenceFieldMissing, testTables, t)
}

func TestLibraryPropertiesSentenceFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Sentence field too short", "SentenceLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesSentenceFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesSentenceFieldSpellCheck(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", ruleresult.Fail, "^grill broccoli now$"},
		{"Non-nil diff but no typos", "SpuriousMisspelledSentenceParagraphValue", ruleresult.Pass, ""},
		{"Correct spelling", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesSentenceFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesParagraphFieldMissing, testTables, t)
}

func TestLibraryPropertiesParagraphFieldSpellCheck(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Misspelled word", "MisspelledSentenceParagraphValue", ruleresult.Fail, "^There is a zebra$"},
		{"Non-nil diff but no typos", "SpuriousMisspelledSentenceParagraphValue", ruleresult.Pass, ""},
		{"Correct spelling", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesParagraphFieldSpellCheck, testTables, t)
}

func TestLibraryPropertiesParagraphFieldRepeatsSentence(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Repeat", "ParagraphRepeatsSentence", ruleresult.Fail, ""},
		{"No repeat", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesParagraphFieldRepeatsSentence, testTables, t)
}

func TestLibraryPropertiesCategoryFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesCategoryFieldMissing, testTables, t)
}

func TestLibraryPropertiesCategoryFieldInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Unsupported category name", "CategoryInvalid", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesCategoryFieldInvalid, testTables, t)
}

func TestLibraryPropertiesCategoryFieldUncategorized(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"No category field", "MissingFields", ruleresult.Skip, ""},
		{"Uncategorized category", "UncategorizedCategoryValue", ruleresult.Fail, ""},
		{"Valid category value", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesCategoryFieldUncategorized, testTables, t)
}

func TestLibraryPropertiesUrlFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesURLFieldMissing, testTables, t)
}

func TestLibraryPropertiesUrlFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"url field too short", "UrlLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesURLFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesUrlFieldInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.NotRun, ""},
		{"Invalid URL format", "UrlFormatInvalid", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesURLFieldInvalid, testTables, t)
}

func TestLibraryPropertiesUrlFieldDeadLink(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.NotRun, ""},
		{"Bad URL", "BadURL", ruleresult.Fail, "^Head \"http://invalid/\": dial tcp: lookup invalid:"},
		{"HTTP error 404", "URL404", ruleresult.Fail, "^404 Not Found$"},
		{"Good URL", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesURLFieldDeadLink, testTables, t)
}

func TestLibraryPropertiesArchitecturesFieldMissing(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Field missing", "MissingFields", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesArchitecturesFieldMissing, testTables, t)
}

func TestLibraryPropertiesArchitecturesFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Architectures field too short", "ArchitecturesLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesArchitecturesFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesArchitecturesFieldSoloAlias(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Solo alias", "ArchitectureAliasSolo", ruleresult.Fail, ""},
		{"Alias w/ true", "ArchitectureAliasWithTrue", ruleresult.Pass, ""},
		{"No alias", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesArchitecturesFieldSoloAlias, testTables, t)
}

func TestLibraryPropertiesArchitecturesFieldValueCase(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Miscased", "ArchitectureMiscased", ruleresult.Fail, ""},
		{"Miscased w/ correct case", "ArchitectureMiscasedWithCorrect", ruleresult.Pass, ""},
		{"Correct case", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesArchitecturesFieldValueCase, testTables, t)
}

func TestLibraryPropertiesDependsFieldInvalidFormat(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Depends field has disallowed characters", "DependsHasBadChars", ruleresult.Fail, ""},
		{"Valid", "DependsValid", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesDependsFieldInvalidFormat, testTables, t)
}

func TestLibraryPropertiesDependsFieldNotInIndex(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"No depends field", "MissingFields", ruleresult.Skip, ""},
		{"Dependency not in index", "DependsNotIndexed", ruleresult.Fail, "^NotIndexed$"},
		{"Dependency constraint not in index", "DependsConstraintNotIndexed", ruleresult.Fail, "^Servo \\(=0\\.0\\.1\\)$"},
		{"Dependencies in index", "DependsIndexed", ruleresult.Pass, ""},
		{"Depends field empty", "DependsEmpty", ruleresult.Pass, ""},
		{"No depends", "NoDepends", ruleresult.Skip, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesDependsFieldNotInIndex, testTables, t)
}

func TestLibraryPropertiesDependsFieldConstraintInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"No depends field", "NoDepends", ruleresult.Skip, ""},
		{"Depends field empty", "DependsEmpty", ruleresult.Pass, ""},
		{"Invalid depends field format", "DependsHasBadChars", ruleresult.Pass, ""},
		{"Invalid constraint syntax", "DependsConstraintInvalid", ruleresult.Fail, "^BarLib \\(nope\\), QuxLib \\(huh\\)$"},
		{"Valid constraint syntax", "DependsConstraintValid", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesDependsFieldConstraintInvalid, testTables, t)
}

func TestLibraryPropertiesDotALinkageFieldInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"dot_a_linkage field invalid value", "DotALinkageInvalid", ruleresult.Fail, ""},
		{"Valid", "DotALinkage", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesDotALinkageFieldInvalid, testTables, t)
}

func TestLibraryPropertiesDotALinkageFieldTrueWithFlatLayout(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Flat layout", "DotALinkageFlat", ruleresult.Fail, ""},
		{"Recursive layout", "DotALinkage", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesDotALinkageFieldTrueWithFlatLayout, testTables, t)
}

func TestLibraryPropertiesIncludesFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Includes field too short", "IncludesLTMinLength", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesIncludesFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesIncludesFieldItemNotFound(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Missing includes", "MissingIncludes", ruleresult.Fail, "^Nonexistent.h$"},
		{"Double comma in includes list", "IncludesListSkip", ruleresult.Pass, ""},
		{"Present includes", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesIncludesFieldItemNotFound, testTables, t)
}

func TestLibraryPropertiesPrecompiledFieldInvalid(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"Precompiled field invalid value", "PrecompiledInvalid", ruleresult.Fail, ""},
		{"Valid", "Precompiled", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesPrecompiledFieldInvalid, testTables, t)
}

func TestLibraryPropertiesPrecompiledFieldEnabledWithFlatLayout(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Unable to load", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Not defined", "MissingFields", ruleresult.Skip, ""},
		{"Flat layout", "PrecompiledFlat", ruleresult.Fail, "^true$"},
		{"Recursive layout", "Precompiled", ruleresult.Pass, ""},
		{"Recursive, not precompiled", "NotPrecompiled", ruleresult.Skip, ""},
		{"Flat, not precompiled", "Flat", ruleresult.Skip, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout, testTables, t)
}

func TestLibraryPropertiesLdflagsFieldLTMinLength(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Legacy", "Legacy", ruleresult.Skip, ""},
		{"ldflags field too short", "LdflagsLTMinLength", ruleresult.Fail, ""},
		{"Valid", "LdflagsValid", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesLdflagsFieldLTMinLength, testTables, t)
}

func TestLibraryPropertiesMisspelledOptionalField(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Invalid", "InvalidLibraryProperties", ruleresult.NotRun, ""},
		{"Misspelled depends field name", "DependsFieldMisspelled", ruleresult.Fail, ""},
		{"Misspelled dot_a_linkage field name", "DotALinkageFieldMisspelled", ruleresult.Fail, ""},
		{"Misspelled includes field name", "IncludesFieldMisspelled", ruleresult.Fail, ""},
		{"Misspelled precompiled field name", "PrecompiledFieldMisspelled", ruleresult.Fail, ""},
		{"Misspelled ldflags field name", "LdflagsFieldMisspelled", ruleresult.Fail, ""},
		{"Valid", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryPropertiesMisspelledOptionalField, testTables, t)
}

func TestLibraryHasStraySketches(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Sketch in root", "SketchInRoot", ruleresult.Fail, ""},
		{"Sketch in subfolder", "MisspelledExamplesFolder", ruleresult.Fail, ""},
		{"Sketch in legit location", "ExamplesFolder", ruleresult.Pass, ""},
		{"No sketches", "Recursive", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(LibraryHasStraySketches, testTables, t)
}

func TestMissingExamples(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"File name collision", "ExamplesFile", ruleresult.Fail, ""},
		{"Has examples", "ExamplesFolder", ruleresult.Pass, ""},
		{`Has examples (in "example" folder)`, "ExampleFolder", ruleresult.Pass, ""},
		{"No examples", "NoExamples", ruleresult.Fail, ""},
	}

	checkLibraryRuleFunction(MissingExamples, testTables, t)
}

func TestMisspelledExamplesFolderName(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Correctly spelled", "ExamplesFolder", ruleresult.Pass, ""},
		{"Misspelled", "MisspelledExamplesFolder", ruleresult.Fail, ""},
		{"No examples folder", "NoExamples", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(MisspelledExamplesFolderName, testTables, t)
}

func TestIncorrectExamplesFolderNameCase(t *testing.T) {
	testTables := []libraryRuleFunctionTestTable{
		{"Correct case", "ExamplesFolder", ruleresult.Pass, ""},
		{"Incorrect case", "IncorrectExamplesFolderCase", ruleresult.Fail, ""},
		{"No examples folder", "NoExamples", ruleresult.Pass, ""},
	}

	checkLibraryRuleFunction(IncorrectExamplesFolderNameCase, testTables, t)
}
