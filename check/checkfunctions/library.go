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

// The check functions for libraries.

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-cli/arduino/libraries"
	"github.com/arduino/arduino-cli/arduino/utils"
	"github.com/arduino/arduino-lint/check/checkdata"
	"github.com/arduino/arduino-lint/check/checkdata/schema"
	"github.com/arduino/arduino-lint/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-lint/check/checkresult"
	"github.com/arduino/arduino-lint/project/library"
	"github.com/arduino/arduino-lint/project/sketch"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	semver "go.bug.st/relaxed-semver"
)

// LibraryInvalid checks whether the provided path is a valid library.
func LibraryInvalid() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() != nil && library.ContainsHeaderFile(checkdata.LoadedLibrary().SourceDir) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, ""
}

// LibraryFolderNameGTMaxLength checks if the library folder name exceeds the maximum length.
func LibraryFolderNameGTMaxLength() (result checkresult.Type, output string) {
	if len(checkdata.ProjectPath().Base()) > 63 {
		return checkresult.Fail, checkdata.ProjectPath().Base()
	}

	return checkresult.Pass, ""
}

// ProhibitedCharactersInLibraryFolderName checks for prohibited characters in the library folder name.
func ProhibitedCharactersInLibraryFolderName() (result checkresult.Type, output string) {
	if !validProjectPathBaseName(checkdata.ProjectPath().Base()) {
		return checkresult.Fail, checkdata.ProjectPath().Base()
	}

	return checkresult.Pass, ""
}

// LibraryHasSubmodule checks whether the library contains a Git submodule.
func LibraryHasSubmodule() (result checkresult.Type, output string) {
	dotGitmodulesPath := checkdata.ProjectPath().Join(".gitmodules")
	hasDotGitmodules, err := dotGitmodulesPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if hasDotGitmodules && dotGitmodulesPath.IsNotDir() {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryContainsSymlinks checks if the library folder contains symbolic links.
func LibraryContainsSymlinks() (result checkresult.Type, output string) {
	projectPathListing, err := checkdata.ProjectPath().ReadDirRecursive()
	if err != nil {
		panic(err)
	}
	projectPathListing.FilterOutDirs()

	symlinkPaths := []string{}
	for _, projectPathItem := range projectPathListing {
		projectPathItemStat, err := os.Lstat(projectPathItem.String())
		if err != nil {
			panic(err)
		}

		if projectPathItemStat.Mode()&os.ModeSymlink != 0 {
			symlinkPaths = append(symlinkPaths, projectPathItem.String())
		}
	}

	if len(symlinkPaths) > 0 {
		return checkresult.Fail, strings.Join(symlinkPaths, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryHasDotDevelopmentFile checks whether the library contains a .development flag file.
func LibraryHasDotDevelopmentFile() (result checkresult.Type, output string) {
	dotDevelopmentPath := checkdata.ProjectPath().Join(".development")
	hasDotDevelopment, err := dotDevelopmentPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if hasDotDevelopment && dotDevelopmentPath.IsNotDir() {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryHasExe checks whether the library contains files with .exe extension.
func LibraryHasExe() (result checkresult.Type, output string) {
	projectPathListing, err := checkdata.ProjectPath().ReadDirRecursive()
	if err != nil {
		panic(err)
	}
	projectPathListing.FilterOutDirs()

	exePaths := []string{}
	for _, projectPathItem := range projectPathListing {
		if projectPathItem.Ext() == ".exe" {
			exePaths = append(exePaths, projectPathItem.String())
		}
	}

	if len(exePaths) > 0 {
		return checkresult.Fail, strings.Join(exePaths, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldHeaderMismatch checks whether the filename of one of the library's header files matches the Library Manager installation folder name.
func LibraryPropertiesNameFieldHeaderMismatch() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	sanitizedName := utils.SanitizeName(name)
	for _, header := range checkdata.SourceHeaders() {
		if strings.TrimSuffix(header, filepath.Ext(header)) == sanitizedName {
			return checkresult.Pass, ""
		}
	}

	return checkresult.Fail, sanitizedName + ".h"
}

// IncorrectLibrarySrcFolderNameCase checks for incorrect case of src subfolder name in recursive format libraries.
func IncorrectLibrarySrcFolderNameCase() (result checkresult.Type, output string) {
	if library.ContainsMetadataFile(checkdata.ProjectPath()) && library.ContainsHeaderFile(checkdata.ProjectPath()) {
		// Flat layout, so no special treatment of src subfolder.
		return checkresult.Skip, "Not applicable due to layout type"
	}

	// The library is intended to have the recursive layout.
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "src")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// RecursiveLibraryWithUtilityFolder checks for presence of a `utility` subfolder in a recursive layout library.
func RecursiveLibraryWithUtilityFolder() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil {
		return checkresult.NotRun, "Library not loaded"
	}

	if checkdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return checkresult.Skip, "Not applicable due to layout type"
	}

	if checkdata.ProjectPath().Join("utility").Exist() {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// MisspelledExtrasFolderName checks for incorrectly spelled `extras` folder name.
func MisspelledExtrasFolderName() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "extras", "(?i)^extra$")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// IncorrectExtrasFolderNameCase checks for incorrect `extras` folder name case.
func IncorrectExtrasFolderNameCase() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "extras")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMissing checks for presence of library.properties.
func LibraryPropertiesMissing() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil {
		return checkresult.NotRun, "Couldn't load library."
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// MisspelledLibraryPropertiesFileName checks for incorrectly spelled library.properties file name.
func MisspelledLibraryPropertiesFileName() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterOutDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "library.properties", "(?i)^librar((y)|(ie))s?[.-_]?propert((y)|(ie))s?$")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// IncorrectLibraryPropertiesFileNameCase checks for incorrect library.properties file name case.
func IncorrectLibraryPropertiesFileNameCase() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterOutDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "library.properties")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// RedundantLibraryProperties checks for redundant copies of the library.properties file.
func RedundantLibraryProperties() (result checkresult.Type, output string) {
	redundantLibraryPropertiesPath := checkdata.ProjectPath().Join("src", "library.properties")
	if redundantLibraryPropertiesPath.Exist() {
		return checkresult.Fail, redundantLibraryPropertiesPath.String()
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesFormat checks for invalid library.properties format.
func LibraryPropertiesFormat() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() != nil && checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has no library.properties"
	}

	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.Fail, checkdata.LibraryPropertiesLoadError().Error()
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldMissing checks for missing library.properties "name" field.
func LibraryPropertiesNameFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "name" value is less than the minimum length.
func LibraryPropertiesNameFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("name") {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldGTMaxLength checks if the library.properties "name" value is greater than the maximum length.
func LibraryPropertiesNameFieldGTMaxLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyGreaterThanMaxLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldGTRecommendedLength checks if the library.properties "name" value is greater than the recommended length.
func LibraryPropertiesNameFieldGTRecommendedLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyGreaterThanMaxLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDisallowedCharacters checks for disallowed characters in the library.properties "name" field.
func LibraryPropertiesNameFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/allowedCharacters", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldStartsWithArduino checks if the library.properties "name" value starts with "Arduino".
func LibraryPropertiesNameFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldMissingOfficialPrefix checks whether the library.properties `name` value uses the prefix required of all new official Arduino libraries.
func LibraryPropertiesNameFieldMissingOfficialPrefix() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if strings.HasPrefix(name, "Arduino_") {
		return checkresult.Pass, ""
	}
	return checkresult.Fail, name
}

// LibraryPropertiesNameFieldContainsArduino checks if the library.properties "name" value contains "Arduino".
func LibraryPropertiesNameFieldContainsArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldHasSpaces checks if the library.properties "name" value contains spaces.
func LibraryPropertiesNameFieldHasSpaces() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSpaces", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldContainsLibrary checks if the library.properties "name" value contains "library".
func LibraryPropertiesNameFieldContainsLibrary() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSuperfluousTerms", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDuplicate checks whether there is an existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldDuplicate() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, "Field not present"
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldNotInIndex checks whether there is no existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, "Field not present"
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, name
}

// LibraryPropertiesVersionFieldMissing checks for missing library.properties "version" field.
func LibraryPropertiesVersionFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonRelaxedSemver checks whether the library.properties "version" value is "relaxed semver" compliant.
func LibraryPropertiesVersionFieldNonRelaxedSemver() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	version, ok := checkdata.LibraryProperties().GetOk("version")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyPatternMismatch("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, version
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonSemver checks whether the library.properties "version" value is semver compliant.
func LibraryPropertiesVersionFieldNonSemver() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	version, ok := checkdata.LibraryProperties().GetOk("version")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyPatternMismatch("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return checkresult.Fail, version
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldBehindTag checks whether a release tag was made without first bumping the library.properties version value.
func LibraryPropertiesVersionFieldBehindTag() (result checkresult.Type, output string) {
	if checkdata.ProjectType() != checkdata.SuperProjectType() {
		return checkresult.Skip, "Not relevant for subprojects"
	}

	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	versionString, ok := checkdata.LibraryProperties().GetOk("version")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	version, err := semver.Parse(versionString)
	if err != nil {
		return checkresult.NotRun, "Can't parse version value"
	}
	logrus.Tracef("version value: %s", version)

	repository, err := git.PlainOpen(checkdata.ProjectPath().String())
	if err != nil {
		return checkresult.Skip, "Project path is not a repository"
	}

	headRef, err := repository.Head()
	if err != nil {
		panic(err)
	}

	headCommit, err := repository.CommitObject(headRef.Hash())
	if err != nil {
		panic(err)
	}

	commits := object.NewCommitIterCTime(headCommit, nil, nil) // Get iterator for the head commit and all its parents in chronological commit time order.

	tagRefs, err := repository.Tags() // Get an iterator of the refs of the repository's tags. These are not in a useful order, so it's necessary to cross-reference them against the commits, which are.

	for { // Iterate over all commits in reverse chronological order.
		commit, err := commits.Next()
		if err != nil {
			// Reached end of commits.
			break
		}

		for { // Iterate over all tag refs.
			tagRef, err := tagRefs.Next()
			if err != nil {
				// Reached end of tags
				break
			}

			// Annotated tags have their own hash, different from the commit hash, so they must be resolved before comparing with the commit.
			resolvedTagRef, err := repository.ResolveRevision(plumbing.Revision(tagRef.Hash().String()))
			if err != nil {
				panic(err)
			}

			if commit.Hash == *resolvedTagRef {
				logrus.Tracef("Found tag: %s", tagRef.Name())

				tagName := strings.TrimPrefix(tagRef.Name().String(), "refs/tags/")
				tagName = strings.TrimPrefix(tagName, "v") // It's common practice to prefix release tag names with "v".
				tagVersion, err := semver.Parse(tagName)
				if err != nil {
					// The normalized tag name is not a recognizable "relaxed semver" version.
					logrus.Tracef("Can't parse tag name.")
					break // Disregard unparsable tags.
				}
				logrus.Tracef("Tag version: %s", tagVersion)

				if tagVersion.GreaterThan(version) {
					logrus.Tracef("Tag is greater.")

					if strings.Contains(tagVersion.String(), "-") {
						// The lack of version bump may have been intentional.
						logrus.Tracef("Tag is pre-release.")
						break
					}

					return checkresult.Fail, fmt.Sprintf("%s vs %s", tagName, versionString)
				}

				return checkresult.Pass, "" // Tag is less than or equal to version field value, all is well.
			}
		}
	}

	return checkresult.Pass, "" // No problems were found.
}

// LibraryPropertiesAuthorFieldMissing checks for missing library.properties "author" field.
func LibraryPropertiesAuthorFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("author", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesAuthorFieldLTMinLength checks if the library.properties "author" value is less than the minimum length.
func LibraryPropertiesAuthorFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("author") {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("author", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldMissing checks for missing library.properties "maintainer" field.
func LibraryPropertiesMaintainerFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("maintainer", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldLTMinLength checks if the library.properties "maintainer" value is less than the minimum length.
func LibraryPropertiesMaintainerFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("maintainer", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldStartsWithArduino checks if the library.properties "maintainer" value starts with "Arduino".
func LibraryPropertiesMaintainerFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	maintainer, ok := checkdata.LibraryProperties().GetOk("maintainer")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/maintainer$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, maintainer
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesEmailFieldAsMaintainerAlias checks whether the library.properties "email" field is being used as an alias for the "maintainer" field.
func LibraryPropertiesEmailFieldAsMaintainerAlias() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("email") {
		return checkresult.Skip, "Field not present"
	}

	if !checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesEmailFieldLTMinLength checks if the library.properties "email" value is less than the minimum length.
func LibraryPropertiesEmailFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LibraryProperties().ContainsKey("maintainer") || !checkdata.LibraryProperties().ContainsKey("email") {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("email", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesEmailFieldStartsWithArduino checks if the library.properties "email" value starts with "Arduino".
func LibraryPropertiesEmailFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.NotRun, "Field not present"
	}

	email, ok := checkdata.LibraryProperties().GetOk("email")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/email$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, email
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesSentenceFieldMissing checks for missing library.properties "sentence" field.
func LibraryPropertiesSentenceFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("sentence", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesSentenceFieldLTMinLength checks if the library.properties "sentence" value is less than the minimum length.
func LibraryPropertiesSentenceFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("sentence") {
		return checkresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("sentence", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesSentenceFieldSpellCheck checks for commonly misspelled words in the library.properties `sentence` field value.
func LibraryPropertiesSentenceFieldSpellCheck() (result checkresult.Type, output string) {
	return spellCheckLibraryPropertiesFieldValue("sentence")
}

// LibraryPropertiesParagraphFieldMissing checks for missing library.properties "paragraph" field.
func LibraryPropertiesParagraphFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("paragraph", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesParagraphFieldSpellCheck checks for commonly misspelled words in the library.properties `paragraph` field value.
func LibraryPropertiesParagraphFieldSpellCheck() (result checkresult.Type, output string) {
	return spellCheckLibraryPropertiesFieldValue("paragraph")
}

// LibraryPropertiesParagraphFieldRepeatsSentence checks whether the library.properties `paragraph` value repeats the `sentence` value.
func LibraryPropertiesParagraphFieldRepeatsSentence() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	sentence, hasSentence := checkdata.LibraryProperties().GetOk("sentence")
	paragraph, hasParagraph := checkdata.LibraryProperties().GetOk("paragraph")

	if !hasSentence || !hasParagraph {
		return checkresult.NotRun, "Field not present"
	}

	if strings.HasPrefix(paragraph, sentence) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldMissing checks for missing library.properties "category" field.
func LibraryPropertiesCategoryFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("category", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldInvalid checks for invalid category in the library.properties "category" field.
func LibraryPropertiesCategoryFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	category, ok := checkdata.LibraryProperties().GetOk("category")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("category", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, category
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldUncategorized checks whether the library.properties "category" value is "Uncategorized".
func LibraryPropertiesCategoryFieldUncategorized() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	category, ok := checkdata.LibraryProperties().GetOk("category")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if category == "Uncategorized" {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldMissing checks for missing library.properties "url" field.
func LibraryPropertiesUrlFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("url", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldInvalid checks whether the library.properties "url" value has a valid URL format.
func LibraryPropertiesUrlFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	url, ok := checkdata.LibraryProperties().GetOk("url")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/url$", "/format$", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, url
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldDeadLink checks whether the URL in the library.properties `url` field can be loaded.
func LibraryPropertiesUrlFieldDeadLink() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	url, ok := checkdata.LibraryProperties().GetOk("url")
	if !ok {
		return checkresult.NotRun, "Field not present"
	}

	logrus.Tracef("Checking URL: %s", url)
	httpResponse, err := http.Get(url)
	if err != nil {
		return checkresult.Fail, err.Error()
	}

	if httpResponse.StatusCode == http.StatusOK {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, httpResponse.Status
}

// LibraryPropertiesArchitecturesFieldMissing checks for missing library.properties "architectures" field.
func LibraryPropertiesArchitecturesFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if checkdata.LoadedLibrary().IsLegacy {
		return checkresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("architectures", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldLTMinLength checks if the library.properties "architectures" value is less than the minimum length.
func LibraryPropertiesArchitecturesFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	if !checkdata.LibraryProperties().ContainsKey("architectures") {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("architectures", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldAlias checks whether an alias architecture name is present, but not its true Arduino architecture name.
func LibraryPropertiesArchitecturesFieldSoloAlias() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	architectures, ok := checkdata.LibraryProperties().GetOk("architectures")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	architecturesList := commaSeparatedToList(strings.ToLower(architectures))

	// Must be all lowercase (there is a separate check for incorrect architecture case).
	var aliases = map[string][]string{
		"atmelavr":      {"avr"},
		"atmelmegaavr":  {"megaavr"},
		"atmelsam":      {"sam", "samd"},
		"espressif32":   {"esp32"},
		"espressif8266": {"esp8266"},
		"intel_arc32":   {"arc32"},
		"nordicnrf52":   {"nRF5", "nrf52", "mbed"},
	}

	trueArchitecturePresent := func(trueArchitecturesQuery []string) bool {
		for _, trueArchitectureQuery := range trueArchitecturesQuery {
			for _, architecture := range architecturesList {
				if architecture == trueArchitectureQuery {
					return true
				}
			}
		}

		return false
	}

	soloAliases := []string{}
	for _, architecture := range architecturesList {
		trueEquivalents, isAlias := aliases[architecture]
		if isAlias && !trueArchitecturePresent(trueEquivalents) {
			soloAliases = append(soloAliases, architecture)
		}
	}

	if len(soloAliases) > 0 {
		return checkresult.Fail, strings.Join(soloAliases, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldValueCase checks for incorrect case of common architectures.
func LibraryPropertiesArchitecturesFieldValueCase() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	architectures, ok := checkdata.LibraryProperties().GetOk("architectures")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	architecturesList := commaSeparatedToList(architectures)

	var commonArchitecturesList = []string{
		"apollo3",
		"arc32",
		"avr",
		"esp32",
		"esp8266",
		"i586",
		"i686",
		"k210",
		"mbed",
		"megaavr",
		"mraa",
		"nRF5",
		"nrf52",
		"pic32",
		"sam",
		"samd",
		"wiced",
		"win10",
	}

	correctArchitecturePresent := func(correctArchitectureQuery string) bool {
		for _, architecture := range architecturesList {
			if architecture == correctArchitectureQuery {
				return true
			}
		}

		return false
	}

	miscasedArchitectures := []string{}
	for _, architecture := range architecturesList {
		for _, commonArchitecture := range commonArchitecturesList {
			if architecture == commonArchitecture {
				break
			}

			if strings.EqualFold(architecture, commonArchitecture) && !correctArchitecturePresent(commonArchitecture) {
				// The architecture has incorrect case and the correctly cased name is not present in the architectures field.
				miscasedArchitectures = append(miscasedArchitectures, architecture)
				break
			}
		}
	}

	if len(miscasedArchitectures) > 0 {
		return checkresult.Fail, strings.Join(miscasedArchitectures, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDependsFieldDisallowedCharacters checks for disallowed characters in the library.properties "depends" field.
func LibraryPropertiesDependsFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	depends, ok := checkdata.LibraryProperties().GetOk("depends")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyPatternMismatch("depends", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, depends
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDependsFieldNotInIndex checks whether the libraries listed in the library.properties `depends` field are in the Library Manager index.
func LibraryPropertiesDependsFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	depends, hasDepends := checkdata.LibraryProperties().GetOk("depends")
	if !hasDepends {
		return checkresult.Skip, "Field not present"
	}

	dependencies := commaSeparatedToList(depends)

	dependenciesNotInIndex := []string{}
	for _, dependency := range dependencies {
		if dependency == "" {
			continue
		}
		logrus.Tracef("Checking if dependency %s is in index.", dependency)
		if !nameInLibraryManagerIndex(dependency) {
			dependenciesNotInIndex = append(dependenciesNotInIndex, dependency)
		}
	}

	if len(dependenciesNotInIndex) > 0 {
		return checkresult.Fail, strings.Join(dependenciesNotInIndex, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDotALinkageFieldInvalid checks for invalid value in the library.properties "dot_a_linkage" field.
func LibraryPropertiesDotALinkageFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Couldn't load library.properties"
	}

	dotALinkage, ok := checkdata.LibraryProperties().GetOk("dot_a_linkage")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("dot_a_linkage", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, dotALinkage
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDotALinkageFieldTrueWithFlatLayout checks whether a library using the "dot_a_linkage" feature has the required recursive layout type.
func LibraryPropertiesDotALinkageFieldTrueWithFlatLayout() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil {
		return checkresult.NotRun, "Library not loaded"
	}

	if !checkdata.LibraryProperties().ContainsKey("dot_a_linkage") {
		return checkresult.Skip, "Field not present"
	}

	if checkdata.LoadedLibrary().DotALinkage && checkdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesIncludesFieldLTMinLength checks if the library.properties "includes" value is less than the minimum length.
func LibraryPropertiesIncludesFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	if !checkdata.LibraryProperties().ContainsKey("includes") {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("includes", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesIncludesFieldItemNotFound checks whether the header files specified in the library.properties `includes` field are in the library.
func LibraryPropertiesIncludesFieldItemNotFound() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	includes, ok := checkdata.LibraryProperties().GetOk("includes")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	includesList := commaSeparatedToList(includes)

	findInclude := func(include string) bool {
		if include == "" {
			return true
		}
		for _, header := range checkdata.SourceHeaders() {
			logrus.Tracef("Comparing include %s with header file %s", include, header)
			if include == header {
				logrus.Tracef("match!")
				return true
			}
		}
		return false
	}

	includesNotInLibrary := []string{}
	for _, include := range includesList {
		if !findInclude(include) {
			includesNotInLibrary = append(includesNotInLibrary, include)
		}
	}

	if len(includesNotInLibrary) > 0 {
		return checkresult.Fail, strings.Join(includesNotInLibrary, ", ")
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesPrecompiledFieldInvalid checks for invalid value in the library.properties "precompiled" field.
func LibraryPropertiesPrecompiledFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	precompiled, ok := checkdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("precompiled", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, precompiled
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout checks whether a precompiled library has the required recursive layout type.
func LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil || checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	precompiled, ok := checkdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	if checkdata.LoadedLibrary().Precompiled && checkdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return checkresult.Fail, precompiled
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesLdflagsFieldLTMinLength checks if the library.properties "ldflags" value is less than the minimum length.
func LibraryPropertiesLdflagsFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	if !checkdata.LibraryProperties().ContainsKey("ldflags") {
		return checkresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("ldflags", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMisspelledOptionalField checks if library.properties contains common misspellings of optional fields.
func LibraryPropertiesMisspelledOptionalField() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	if schema.MisspelledOptionalPropertyFound(checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryHasStraySketches checks for sketches outside the `examples` and `extras` folders.
func LibraryHasStraySketches() (result checkresult.Type, output string) {
	straySketchPaths := []string{}
	if sketch.ContainsMainSketchFile(checkdata.ProjectPath()) { // Check library root.
		straySketchPaths = append(straySketchPaths, checkdata.ProjectPath().String())
	}

	// Check subfolders.
	projectPathListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	projectPathListing.FilterDirs()

	for _, topLevelSubfolder := range projectPathListing {
		if topLevelSubfolder.Base() == "examples" || topLevelSubfolder.Base() == "extras" {
			continue // Skip valid sketch locations.
		}

		topLevelSubfolderRecursiveListing, err := topLevelSubfolder.ReadDirRecursive()
		if err != nil {
			panic(err)
		}
		topLevelSubfolderRecursiveListing.FilterDirs()

		for _, subfolder := range topLevelSubfolderRecursiveListing {
			if sketch.ContainsMainSketchFile(subfolder) {
				straySketchPaths = append(straySketchPaths, subfolder.String())
			}
		}
	}

	if len(straySketchPaths) > 0 {
		return checkresult.Fail, strings.Join(straySketchPaths, ", ")
	}

	return checkresult.Pass, ""
}

// MissingExamples checks whether the library is missing examples.
func MissingExamples() (result checkresult.Type, output string) {
	for _, examplesFolderName := range library.ExamplesFolderSupportedNames() {
		examplesPath := checkdata.ProjectPath().Join(examplesFolderName)
		exists, err := examplesPath.IsDirCheck()
		if err != nil {
			panic(err)
		}
		if exists {
			directoryListing, _ := examplesPath.ReadDirRecursive()
			directoryListing.FilterDirs()
			for _, potentialExamplePath := range directoryListing {
				if sketch.ContainsMainSketchFile(potentialExamplePath) {
					return checkresult.Pass, ""
				}
			}
		}
	}

	return checkresult.Fail, ""
}

// MisspelledExamplesFolderName checks for incorrectly spelled `examples` folder name.
func MisspelledExamplesFolderName() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "examples", "(?i)^e((x)|(xs)|(s))((am)|(ma))p((le)|(el))s?$")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// IncorrectExamplesFolderNameCase checks for incorrect `examples` folder name case.
func IncorrectExamplesFolderNameCase() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "examples")
	if found {
		return checkresult.Fail, path.String()
	}

	return checkresult.Pass, ""
}

// nameInLibraryManagerIndex returns whether there is a library in Library Manager index using the given name.
func nameInLibraryManagerIndex(name string) bool {
	libraries := checkdata.LibraryManagerIndex()["libraries"].([]interface{})
	for _, libraryInterface := range libraries {
		library := libraryInterface.(map[string]interface{})
		if library["name"].(string) == name {
			return true
		}
	}

	return false
}

// spellCheckLibraryPropertiesFieldValue returns the value of the provided library.properties field with commonly misspelled words corrected.
func spellCheckLibraryPropertiesFieldValue(fieldName string) (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, "Library not loaded"
	}

	fieldValue, ok := checkdata.LibraryProperties().GetOk(fieldName)
	if !ok {
		return checkresult.Skip, "Field not present"
	}

	replaced, diff := checkdata.MisspelledWordsReplacer().Replace(fieldValue)
	if len(diff) > 0 {
		return checkresult.Fail, replaced
	}

	return checkresult.Pass, ""
}

// commaSeparatedToList returns the list equivalent of a comma-separated string.
func commaSeparatedToList(commaSeparated string) []string {
	list := []string{}
	for _, item := range strings.Split(commaSeparated, ",") {
		list = append(list, strings.TrimSpace(item))
	}

	return list
}
