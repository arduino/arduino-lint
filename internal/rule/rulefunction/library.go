// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package rulefunction

// The rule functions for libraries.

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-cli/arduino/libraries"
	"github.com/arduino/arduino-cli/arduino/utils"
	"github.com/arduino/arduino-lint/internal/project/library"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/sketch"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	semver "go.bug.st/relaxed-semver"
)

// LibraryInvalid checks whether the provided path is a valid library.
func LibraryInvalid() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() != nil && library.ContainsHeaderFile(projectdata.LoadedLibrary().SourceDir) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, ""
}

// LibraryFolderNameGTMaxLength checks if the library folder name exceeds the maximum length.
func LibraryFolderNameGTMaxLength() (result ruleresult.Type, output string) {
	if len(projectdata.ProjectPath().Base()) > 63 {
		return ruleresult.Fail, projectdata.ProjectPath().Base()
	}

	return ruleresult.Pass, ""
}

// ProhibitedCharactersInLibraryFolderName checks for prohibited characters in the library folder name.
func ProhibitedCharactersInLibraryFolderName() (result ruleresult.Type, output string) {
	if !validProjectPathBaseName(projectdata.ProjectPath().Base()) {
		return ruleresult.Fail, projectdata.ProjectPath().Base()
	}

	return ruleresult.Pass, ""
}

// LibraryHasSubmodule checks whether the library contains a Git submodule.
func LibraryHasSubmodule() (result ruleresult.Type, output string) {
	dotGitmodulesPath := projectdata.ProjectPath().Join(".gitmodules")
	hasDotGitmodules, err := dotGitmodulesPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if hasDotGitmodules && dotGitmodulesPath.IsNotDir() {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryContainsSymlinks checks if the library folder contains symbolic links.
func LibraryContainsSymlinks() (result ruleresult.Type, output string) {
	projectPathListing, err := projectdata.ProjectPath().ReadDirRecursive()
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
		return ruleresult.Fail, strings.Join(symlinkPaths, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryHasDotDevelopmentFile checks whether the library contains a .development flag file.
func LibraryHasDotDevelopmentFile() (result ruleresult.Type, output string) {
	dotDevelopmentPath := projectdata.ProjectPath().Join(".development")
	hasDotDevelopment, err := dotDevelopmentPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if hasDotDevelopment && dotDevelopmentPath.IsNotDir() {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryHasExe checks whether the library contains files with .exe extension.
func LibraryHasExe() (result ruleresult.Type, output string) {
	projectPathListing, err := projectdata.ProjectPath().ReadDirRecursive()
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
		return ruleresult.Fail, strings.Join(exePaths, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldHeaderMismatch checks whether the filename of one of the library's header files matches the Library Manager installation folder name.
func LibraryPropertiesNameFieldHeaderMismatch() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	sanitizedName := utils.SanitizeName(name)
	for _, header := range projectdata.SourceHeaders() {
		if strings.TrimSuffix(header, filepath.Ext(header)) == sanitizedName {
			return ruleresult.Pass, ""
		}
	}

	return ruleresult.Fail, sanitizedName + ".h"
}

// IncorrectLibrarySrcFolderNameCase checks for incorrect case of src subfolder name in recursive format libraries.
func IncorrectLibrarySrcFolderNameCase() (result ruleresult.Type, output string) {
	if library.ContainsMetadataFile(projectdata.ProjectPath()) && library.ContainsHeaderFile(projectdata.ProjectPath()) {
		// Flat layout, so no special treatment of src subfolder.
		return ruleresult.Skip, "Not applicable due to layout type"
	}

	// The library is intended to have the recursive layout.
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "src")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// RecursiveLibraryWithUtilityFolder checks for presence of a `utility` subfolder in a recursive layout library.
func RecursiveLibraryWithUtilityFolder() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() == nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	if projectdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return ruleresult.Skip, "Not applicable due to layout type"
	}

	if projectdata.ProjectPath().Join("utility").Exist() {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// MisspelledExtrasFolderName checks for incorrectly spelled `extras` folder name.
func MisspelledExtrasFolderName() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "extras", "(?i)^extra$")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// IncorrectExtrasFolderNameCase checks for incorrect `extras` folder name case.
func IncorrectExtrasFolderNameCase() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "extras")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesMissing checks for presence of library.properties.
func LibraryPropertiesMissing() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() == nil {
		return ruleresult.NotRun, "Couldn't load library."
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// MisspelledLibraryPropertiesFileName checks for incorrectly spelled library.properties file name.
func MisspelledLibraryPropertiesFileName() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterOutDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "library.properties", "(?i)^librar((y)|(ie))s?[.-_]?propert((y)|(ie))s?$")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// IncorrectLibraryPropertiesFileNameCase checks for incorrect library.properties file name case.
func IncorrectLibraryPropertiesFileNameCase() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterOutDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "library.properties")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// RedundantLibraryProperties checks for redundant copies of the library.properties file.
func RedundantLibraryProperties() (result ruleresult.Type, output string) {
	redundantLibraryPropertiesPath := projectdata.ProjectPath().Join("src", "library.properties")
	if redundantLibraryPropertiesPath.Exist() {
		return ruleresult.Fail, redundantLibraryPropertiesPath.String()
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesFormat checks for invalid library.properties format.
func LibraryPropertiesFormat() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() != nil && projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has no library.properties"
	}

	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.Fail, projectdata.LibraryPropertiesLoadError().Error()
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldMissing checks for missing library.properties "name" field.
func LibraryPropertiesNameFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("name", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "name" value is less than the minimum length.
func LibraryPropertiesNameFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("name") {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("name", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldGTMaxLength checks if the library.properties "name" value is greater than the maximum length.
func LibraryPropertiesNameFieldGTMaxLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyGreaterThanMaxLength("name", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldGTRecommendedLength checks if the library.properties "name" value is greater than the recommended length.
func LibraryPropertiesNameFieldGTRecommendedLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyGreaterThanMaxLength("name", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldDisallowedCharacters checks for disallowed characters in the library.properties "name" field.
func LibraryPropertiesNameFieldDisallowedCharacters() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/allowedCharacters", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldStartsWithArduino checks if the library.properties "name" value starts with "Arduino".
func LibraryPropertiesNameFieldStartsWithArduino() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notStartsWithArduino", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldMissingOfficialPrefix checks whether the library.properties `name` value uses the prefix required of all new official Arduino libraries.
func LibraryPropertiesNameFieldMissingOfficialPrefix() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if strings.HasPrefix(name, "Arduino_") {
		return ruleresult.Pass, ""
	}
	return ruleresult.Fail, name
}

// LibraryPropertiesNameFieldContainsArduino checks if the library.properties "name" value contains "Arduino".
func LibraryPropertiesNameFieldContainsArduino() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsArduino", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldHasSpaces checks if the library.properties "name" value contains spaces.
func LibraryPropertiesNameFieldHasSpaces() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSpaces", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldContainsLibrary checks if the library.properties "name" value contains "library".
func LibraryPropertiesNameFieldContainsLibrary() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, ok := projectdata.LibraryProperties().GetOk("name")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSuperfluousTerms", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldDuplicate checks whether there is an existing entry in the Library Manager index using the library.properties `name` value.
func LibraryPropertiesNameFieldDuplicate() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, hasName := projectdata.LibraryProperties().GetOk("name")
	if !hasName {
		return ruleresult.NotRun, "Field not present"
	}

	if nameInLibraryManagerIndex(name) {
		return ruleresult.Fail, name
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesNameFieldNotInIndex checks whether there is no existing entry in the Library Manager index using the library.properties `name` value.
func LibraryPropertiesNameFieldNotInIndex() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	name, hasName := projectdata.LibraryProperties().GetOk("name")
	if !hasName {
		return ruleresult.NotRun, "Field not present"
	}

	if nameInLibraryManagerIndex(name) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, name
}

// LibraryPropertiesVersionFieldMissing checks for missing library.properties "version" field.
func LibraryPropertiesVersionFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("version", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonRelaxedSemver checks whether the library.properties "version" value is "relaxed semver" compliant.
func LibraryPropertiesVersionFieldNonRelaxedSemver() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	version, ok := projectdata.LibraryProperties().GetOk("version")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyPatternMismatch("version", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, version
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonSemver checks whether the library.properties "version" value is semver compliant.
func LibraryPropertiesVersionFieldNonSemver() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	version, ok := projectdata.LibraryProperties().GetOk("version")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyPatternMismatch("version", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, version
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesVersionFieldBehindTag checks whether a release tag was made without first bumping the library.properties version value.
func LibraryPropertiesVersionFieldBehindTag() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	versionString, ok := projectdata.LibraryProperties().GetOk("version")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	version, err := semver.Parse(versionString)
	if err != nil {
		return ruleresult.NotRun, "Can't parse version value"
	}
	logrus.Tracef("version value: %s", version)

	repository, err := git.PlainOpen(projectdata.ProjectPath().String())
	if err != nil {
		return ruleresult.Skip, "Project path is not a repository"
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

					return ruleresult.Fail, fmt.Sprintf("%s vs %s", tagName, versionString)
				}

				return ruleresult.Pass, "" // Tag is less than or equal to version field value, all is well.
			}
		}
	}

	return ruleresult.Pass, "" // No problems were found.
}

// LibraryPropertiesAuthorFieldMissing checks for missing library.properties "author" field.
func LibraryPropertiesAuthorFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("author", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesAuthorFieldLTMinLength checks if the library.properties "author" value is less than the minimum length.
func LibraryPropertiesAuthorFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("author") {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("author", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldMissing checks for missing library.properties "maintainer" field.
func LibraryPropertiesMaintainerFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("maintainer", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldLTMinLength checks if the library.properties "maintainer" value is less than the minimum length.
func LibraryPropertiesMaintainerFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("maintainer") {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("maintainer", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldStartsWithArduino checks if the library.properties "maintainer" value starts with "Arduino".
func LibraryPropertiesMaintainerFieldStartsWithArduino() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	maintainer, ok := projectdata.LibraryProperties().GetOk("maintainer")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/maintainer$", "/patternObjects/notStartsWithArduino", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, maintainer
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesEmailFieldAsMaintainerAlias checks whether the library.properties "email" field is being used as an alias for the "maintainer" field.
func LibraryPropertiesEmailFieldAsMaintainerAlias() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("email") {
		return ruleresult.Skip, "Field not present"
	}

	if !projectdata.LibraryProperties().ContainsKey("maintainer") {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesEmailFieldLTMinLength checks if the library.properties "email" value is less than the minimum length.
func LibraryPropertiesEmailFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LibraryProperties().ContainsKey("maintainer") || !projectdata.LibraryProperties().ContainsKey("email") {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("email", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesEmailFieldStartsWithArduino checks if the library.properties "email" value starts with "Arduino".
func LibraryPropertiesEmailFieldStartsWithArduino() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LibraryProperties().ContainsKey("maintainer") {
		return ruleresult.Skip, "No email alias field"
	}

	email, ok := projectdata.LibraryProperties().GetOk("email")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/email$", "/patternObjects/notStartsWithArduino", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, email
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesSentenceFieldMissing checks for missing library.properties "sentence" field.
func LibraryPropertiesSentenceFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("sentence", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesSentenceFieldLTMinLength checks if the library.properties "sentence" value is less than the minimum length.
func LibraryPropertiesSentenceFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("sentence") {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("sentence", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesSentenceFieldSpellCheck checks for commonly misspelled words in the library.properties `sentence` field value.
func LibraryPropertiesSentenceFieldSpellCheck() (result ruleresult.Type, output string) {
	return spellCheckLibraryPropertiesFieldValue("sentence")
}

// LibraryPropertiesParagraphFieldMissing checks for missing library.properties "paragraph" field.
func LibraryPropertiesParagraphFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("paragraph", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesParagraphFieldSpellCheck checks for commonly misspelled words in the library.properties `paragraph` field value.
func LibraryPropertiesParagraphFieldSpellCheck() (result ruleresult.Type, output string) {
	return spellCheckLibraryPropertiesFieldValue("paragraph")
}

// LibraryPropertiesParagraphFieldRepeatsSentence checks whether the library.properties `paragraph` value repeats the `sentence` value.
func LibraryPropertiesParagraphFieldRepeatsSentence() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	sentence, hasSentence := projectdata.LibraryProperties().GetOk("sentence")
	paragraph, hasParagraph := projectdata.LibraryProperties().GetOk("paragraph")

	if !hasSentence || !hasParagraph {
		return ruleresult.NotRun, "Field not present"
	}

	if strings.HasPrefix(paragraph, sentence) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesCategoryFieldMissing checks for missing library.properties "category" field.
func LibraryPropertiesCategoryFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("category", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesCategoryFieldInvalid checks for invalid category in the library.properties "category" field.
func LibraryPropertiesCategoryFieldInvalid() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	category, ok := projectdata.LibraryProperties().GetOk("category")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("category", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, category
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesCategoryFieldUncategorized checks whether the library.properties "category" value is "Uncategorized".
func LibraryPropertiesCategoryFieldUncategorized() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	category, ok := projectdata.LibraryProperties().GetOk("category")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if category == "Uncategorized" {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesUrlFieldMissing checks for missing library.properties "url" field.
func LibraryPropertiesUrlFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("url", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesUrlFieldLTMinLength checks if the library.properties "url" value is less than the minimum length.
func LibraryPropertiesUrlFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("url") {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.PropertyLessThanMinLength("url", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Permissive]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesUrlFieldInvalid checks whether the library.properties "url" value has a valid URL format.
func LibraryPropertiesUrlFieldInvalid() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	url, ok := projectdata.LibraryProperties().GetOk("url")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	if schema.ValidationErrorMatch("^#/url$", "/format$", "", "", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, url
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesUrlFieldDeadLink checks whether the URL in the library.properties `url` field can be loaded.
func LibraryPropertiesUrlFieldDeadLink() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	url, ok := projectdata.LibraryProperties().GetOk("url")
	if !ok {
		return ruleresult.NotRun, "Field not present"
	}

	logrus.Tracef("Checking URL: %s", url)
	httpResponse, err := http.Get(url)
	if err != nil {
		return ruleresult.Fail, err.Error()
	}

	if httpResponse.StatusCode == http.StatusOK {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, httpResponse.Status
}

// LibraryPropertiesArchitecturesFieldMissing checks for missing library.properties "architectures" field.
func LibraryPropertiesArchitecturesFieldMissing() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if projectdata.LoadedLibrary().IsLegacy {
		return ruleresult.Skip, "Library has legacy format"
	}

	if schema.RequiredPropertyMissing("architectures", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}
	return ruleresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldLTMinLength checks if the library.properties "architectures" value is less than the minimum length.
func LibraryPropertiesArchitecturesFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	if !projectdata.LibraryProperties().ContainsKey("architectures") {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("architectures", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldSoloAlias checks whether an alias architecture name is present, but not its true Arduino architecture name.
func LibraryPropertiesArchitecturesFieldSoloAlias() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	architectures, ok := projectdata.LibraryProperties().GetOk("architectures")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	architecturesList := commaSeparatedToList(strings.ToLower(architectures))

	// Must be all lowercase (there is a separate rule for incorrect architecture case).
	var aliases = map[string][]string{
		"atmelavr":      {"avr"},
		"atmelmegaavr":  {"megaavr"},
		"atmelsam":      {"sam", "samd"},
		"espressif32":   {"esp32"},
		"espressif8266": {"esp8266"},
		"intel_arc32":   {"arc32"},
		"nordicnrf52":   {"nRF5", "nrf52", "mbed", "mbed_edge", "mbed_nano"},
		"raspberrypi":   {"mbed_nano", "mbed_rp2040", "rp2040"},
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
		return ruleresult.Fail, strings.Join(soloAliases, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesArchitecturesFieldValueCase checks for incorrect case of common architectures.
func LibraryPropertiesArchitecturesFieldValueCase() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	architectures, ok := projectdata.LibraryProperties().GetOk("architectures")
	if !ok {
		return ruleresult.Skip, "Field not present"
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
		"mbed_edge",
		"mbed_nano",
		"mbed_portenta",
		"mbed_rp2040",
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
		return ruleresult.Fail, strings.Join(miscasedArchitectures, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesDependsFieldDisallowedCharacters checks for disallowed characters in the library.properties "depends" field.
func LibraryPropertiesDependsFieldDisallowedCharacters() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	depends, ok := projectdata.LibraryProperties().GetOk("depends")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyPatternMismatch("depends", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, depends
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesDependsFieldNotInIndex checks whether the libraries listed in the library.properties `depends` field are in the Library Manager index.
func LibraryPropertiesDependsFieldNotInIndex() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	depends, hasDepends := projectdata.LibraryProperties().GetOk("depends")
	if !hasDepends {
		return ruleresult.Skip, "Field not present"
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
		return ruleresult.Fail, strings.Join(dependenciesNotInIndex, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesDotALinkageFieldInvalid checks for invalid value in the library.properties "dot_a_linkage" field.
func LibraryPropertiesDotALinkageFieldInvalid() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load library.properties"
	}

	dotALinkage, ok := projectdata.LibraryProperties().GetOk("dot_a_linkage")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("dot_a_linkage", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, dotALinkage
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesDotALinkageFieldTrueWithFlatLayout checks whether a library using the "dot_a_linkage" feature has the required recursive layout type.
func LibraryPropertiesDotALinkageFieldTrueWithFlatLayout() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() == nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	if !projectdata.LibraryProperties().ContainsKey("dot_a_linkage") {
		return ruleresult.Skip, "Field not present"
	}

	if projectdata.LoadedLibrary().DotALinkage && projectdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesIncludesFieldLTMinLength checks if the library.properties "includes" value is less than the minimum length.
func LibraryPropertiesIncludesFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	if !projectdata.LibraryProperties().ContainsKey("includes") {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("includes", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesIncludesFieldItemNotFound checks whether the header files specified in the library.properties `includes` field are in the library.
func LibraryPropertiesIncludesFieldItemNotFound() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	includes, ok := projectdata.LibraryProperties().GetOk("includes")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	includesList := commaSeparatedToList(includes)

	findInclude := func(include string) bool {
		if include == "" {
			return true
		}
		for _, header := range projectdata.SourceHeaders() {
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
		return ruleresult.Fail, strings.Join(includesNotInLibrary, ", ")
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesPrecompiledFieldInvalid checks for invalid value in the library.properties "precompiled" field.
func LibraryPropertiesPrecompiledFieldInvalid() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	precompiled, ok := projectdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyEnumMismatch("precompiled", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, precompiled
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout checks whether a precompiled library has the required recursive layout type.
func LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout() (result ruleresult.Type, output string) {
	if projectdata.LoadedLibrary() == nil || projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	precompiled, ok := projectdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	if projectdata.LoadedLibrary().Precompiled && projectdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return ruleresult.Fail, precompiled
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesLdflagsFieldLTMinLength checks if the library.properties "ldflags" value is less than the minimum length.
func LibraryPropertiesLdflagsFieldLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	if !projectdata.LibraryProperties().ContainsKey("ldflags") {
		return ruleresult.Skip, "Field not present"
	}

	if schema.PropertyLessThanMinLength("ldflags", projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryPropertiesMisspelledOptionalField checks if library.properties contains common misspellings of optional fields.
func LibraryPropertiesMisspelledOptionalField() (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	if schema.MisspelledOptionalPropertyFound(projectdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// LibraryHasStraySketches checks for sketches outside the `examples` and `extras` folders.
func LibraryHasStraySketches() (result ruleresult.Type, output string) {
	straySketchPaths := []string{}
	if sketch.ContainsMainSketchFile(projectdata.ProjectPath()) { // Check library root.
		straySketchPaths = append(straySketchPaths, projectdata.ProjectPath().String())
	}

	// Check subfolders.
	projectPathListing, err := projectdata.ProjectPath().ReadDir()
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
		return ruleresult.Fail, strings.Join(straySketchPaths, ", ")
	}

	return ruleresult.Pass, ""
}

// MissingExamples checks whether the library is missing examples.
func MissingExamples() (result ruleresult.Type, output string) {
	for _, examplesFolderName := range library.ExamplesFolderSupportedNames() {
		examplesPath := projectdata.ProjectPath().Join(examplesFolderName)
		exists, err := examplesPath.IsDirCheck()
		if err != nil {
			panic(err)
		}
		if exists {
			directoryListing, _ := examplesPath.ReadDirRecursive()
			directoryListing.FilterDirs()
			for _, potentialExamplePath := range directoryListing {
				if sketch.ContainsMainSketchFile(potentialExamplePath) {
					return ruleresult.Pass, ""
				}
			}
		}
	}

	return ruleresult.Fail, ""
}

// MisspelledExamplesFolderName checks for incorrectly spelled `examples` folder name.
func MisspelledExamplesFolderName() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsMisspelledPathBaseName(directoryListing, "examples", "(?i)^e((x)|(xs)|(s))((am)|(ma))p((le)|(el))s?$")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// IncorrectExamplesFolderNameCase checks for incorrect `examples` folder name case.
func IncorrectExamplesFolderNameCase() (result ruleresult.Type, output string) {
	directoryListing, err := projectdata.ProjectPath().ReadDir()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterDirs()

	path, found := containsIncorrectPathBaseCase(directoryListing, "examples")
	if found {
		return ruleresult.Fail, path.String()
	}

	return ruleresult.Pass, ""
}

// nameInLibraryManagerIndex returns whether there is a library in Library Manager index using the given name.
func nameInLibraryManagerIndex(name string) bool {
	library := projectdata.LibraryManagerIndex().Index.FindIndexedLibrary(&libraries.Library{Name: name})
	return library != nil
}

// spellCheckLibraryPropertiesFieldValue returns the value of the provided library.properties field with commonly misspelled words corrected.
func spellCheckLibraryPropertiesFieldValue(fieldName string) (result ruleresult.Type, output string) {
	if projectdata.LibraryPropertiesLoadError() != nil {
		return ruleresult.NotRun, "Library not loaded"
	}

	fieldValue, ok := projectdata.LibraryProperties().GetOk(fieldName)
	if !ok {
		return ruleresult.Skip, "Field not present"
	}

	replaced, diff := projectdata.MisspelledWordsReplacer().Replace(fieldValue)
	if len(diff) > 0 {
		return ruleresult.Fail, replaced
	}

	return ruleresult.Pass, ""
}

// commaSeparatedToList returns the list equivalent of a comma-separated string.
func commaSeparatedToList(commaSeparated string) []string {
	list := []string{}
	for _, item := range strings.Split(commaSeparated, ",") {
		list = append(list, strings.TrimSpace(item))
	}

	return list
}
