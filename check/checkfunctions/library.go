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

// The check functions for libraries.

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkdata/schema"
	"github.com/arduino/arduino-check/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project/library"
	"github.com/arduino/arduino-cli/arduino/libraries"
	"github.com/arduino/arduino-cli/arduino/utils"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/sirupsen/logrus"
)

// LibraryPropertiesFormat checks for invalid library.properties format.
func LibraryPropertiesFormat() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.Fail, checkdata.LibraryPropertiesLoadError().Error()
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesMissing checks for presence of library.properties.
func LibraryPropertiesMissing() (result checkresult.Type, output string) {
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

	path, found := containsMisspelledPathBaseName(directoryListing, "library.properties", "(?i)^librar(y)|(ie)s?[.-_]?propert(y)|(ie)s?$")
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

// LibraryPropertiesNameFieldMissing checks for missing library.properties "name" field.
func LibraryPropertiesNameFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "name" value is less than the minimum length.
func LibraryPropertiesNameFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("name") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldGTMaxLength checks if the library.properties "name" value is greater than the maximum length.
func LibraryPropertiesNameFieldGTMaxLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyGreaterThanMaxLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldGTRecommendedLength checks if the library.properties "name" value is greater than the recommended length.
func LibraryPropertiesNameFieldGTRecommendedLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, checkdata.LibraryProperties().Get("name")
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyGreaterThanMaxLength("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDisallowedCharacters checks for disallowed characters in the library.properties "name" field.
func LibraryPropertiesNameFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.PropertyPatternMismatch("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldHasSpaces checks if the library.properties "name" value contains spaces.
func LibraryPropertiesNameFieldHasSpaces() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSpaces", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldStartsWithArduino checks if the library.properties "name" value starts with "Arduino".
func LibraryPropertiesNameFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldMissingOfficialPrefix checks whether the library.properties `name` value uses the prefix required of all new official Arduino libraries.
func LibraryPropertiesNameFieldMissingOfficialPrefix() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if strings.HasPrefix(name, "Arduino_") {
		return checkresult.Pass, ""
	}
	return checkresult.Fail, name
}

// LibraryPropertiesNameFieldContainsArduino checks if the library.properties "name" value contains "Arduino".
func LibraryPropertiesNameFieldContainsArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldContainsLibrary checks if the library.properties "name" value contains "library".
func LibraryPropertiesNameFieldContainsLibrary() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/name$", "/patternObjects/notContainsSuperfluousTerms", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict], configuration.SchemasPath()) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDuplicate checks whether there is an existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldDuplicate() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, ""
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldNotInIndex checks whether there is no existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, ""
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, name
}

// LibraryPropertiesNameFieldHeaderMismatch checks whether the filename of one of the library's header files matches the Library Manager installation folder name.
func LibraryPropertiesNameFieldHeaderMismatch() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, ok := checkdata.LibraryProperties().GetOk("name")
	if !ok {
		return checkresult.NotRun, ""
	}

	sanitizedName := utils.SanitizeName(name)
	for _, header := range checkdata.SourceHeaders() {
		if strings.TrimSuffix(header, filepath.Ext(header)) == sanitizedName {
			return checkresult.Pass, ""
		}
	}

	return checkresult.Fail, sanitizedName + ".h"
}

// LibraryPropertiesVersionFieldMissing checks for missing library.properties "version" field.
func LibraryPropertiesVersionFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonRelaxedSemver checks whether the library.properties "version" value is "relaxed semver" compliant.
func LibraryPropertiesVersionFieldNonRelaxedSemver() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	version, ok := checkdata.LibraryProperties().GetOk("version")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyPatternMismatch("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, version
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldNonSemver checks whether the library.properties "version" value is semver compliant.
func LibraryPropertiesVersionFieldNonSemver() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	version, ok := checkdata.LibraryProperties().GetOk("version")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyPatternMismatch("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Strict], configuration.SchemasPath()) {
		return checkresult.Fail, version
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesAuthorFieldMissing checks for missing library.properties "author" field.
func LibraryPropertiesAuthorFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("author", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesAuthorFieldLTMinLength checks if the library.properties "author" value is less than the minimum length.
func LibraryPropertiesAuthorFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("author") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("author", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldMissing checks for missing library.properties "maintainer" field.
func LibraryPropertiesMaintainerFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("maintainer", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldLTMinLength checks if the library.properties "maintainer" value is less than the minimum length.
func LibraryPropertiesMaintainerFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("maintainer", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldStartsWithArduino checks if the library.properties "maintainer" value starts with "Arduino".
func LibraryPropertiesMaintainerFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	maintainer, ok := checkdata.LibraryProperties().GetOk("maintainer")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/maintainer$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, maintainer
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesEmailFieldAsMaintainerAlias checks whether the library.properties "email" field is being used as an alias for the "maintainer" field.
func LibraryPropertiesEmailFieldAsMaintainerAlias() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("email") {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "email" value is less than the minimum length.
func LibraryPropertiesEmailFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if checkdata.LibraryProperties().ContainsKey("maintainer") || !checkdata.LibraryProperties().ContainsKey("email") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("email", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMaintainerFieldStartsWithArduino checks if the library.properties "email" value starts with "Arduino".
func LibraryPropertiesEmailFieldStartsWithArduino() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if checkdata.LibraryProperties().ContainsKey("maintainer") {
		return checkresult.NotRun, ""
	}

	email, ok := checkdata.LibraryProperties().GetOk("email")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/email$", "/patternObjects/notStartsWithArduino", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, email
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesSentenceFieldMissing checks for missing library.properties "sentence" field.
func LibraryPropertiesSentenceFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("sentence", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesSentenceFieldLTMinLength checks if the library.properties "sentence" value is less than the minimum length.
func LibraryPropertiesSentenceFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("sentence") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("sentence", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
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
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("paragraph", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
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
		return checkresult.NotRun, ""
	}

	sentence, hasSentence := checkdata.LibraryProperties().GetOk("sentence")
	paragraph, hasParagraph := checkdata.LibraryProperties().GetOk("paragraph")

	if !hasSentence || !hasParagraph {
		return checkresult.NotRun, ""
	}

	if strings.HasPrefix(paragraph, sentence) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldMissing checks for missing library.properties "category" field.
func LibraryPropertiesCategoryFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("category", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldInvalid checks for invalid category in the library.properties "category" field.
func LibraryPropertiesCategoryFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	category, ok := checkdata.LibraryProperties().GetOk("category")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyEnumMismatch("category", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, category
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesCategoryFieldUncategorized checks whether the library.properties "category" value is "Uncategorized".
func LibraryPropertiesCategoryFieldUncategorized() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	category, ok := checkdata.LibraryProperties().GetOk("category")
	if !ok {
		return checkresult.NotRun, ""
	}

	if category == "Uncategorized" {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldMissing checks for missing library.properties "url" field.
func LibraryPropertiesUrlFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("url", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldInvalid checks whether the library.properties "url" value has a valid URL format.
func LibraryPropertiesUrlFieldInvalid() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	url, ok := checkdata.LibraryProperties().GetOk("url")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.ValidationErrorMatch("^#/url$", "/format$", "", "", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, url
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesUrlFieldDeadLink checks whether the URL in the library.properties `url` field can be loaded.
func LibraryPropertiesUrlFieldDeadLink() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	url, ok := checkdata.LibraryProperties().GetOk("url")
	if !ok {
		return checkresult.NotRun, ""
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
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("architectures", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "architectures" value is less than the minimum length.
func LibraryPropertiesArchitecturesFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("architectures") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("architectures", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDependsFieldDisallowedCharacters checks for disallowed characters in the library.properties "depends" field.
func LibraryPropertiesDependsFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	depends, ok := checkdata.LibraryProperties().GetOk("depends")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyPatternMismatch("depends", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, depends
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDependsFieldNotInIndex checks whether the libraries listed in the library.properties `depends` field are in the Library Manager index.
func LibraryPropertiesDependsFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	depends, hasDepends := checkdata.LibraryProperties().GetOk("depends")
	if !hasDepends {
		return checkresult.NotRun, ""
	}

	dependencies, err := properties.SplitQuotedString(depends, "", false)
	if err != nil {
		panic(err)
	}
	dependenciesNotInIndex := []string{}
	for _, dependency := range dependencies {
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
		return checkresult.NotRun, ""
	}

	dotALinkage, ok := checkdata.LibraryProperties().GetOk("dot_a_linkage")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyEnumMismatch("dot_a_linkage", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, dotALinkage
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesDotALinkageFieldTrueWithFlatLayout checks whether a library using the "dot_a_linkage" feature has the required recursive layout type.
func LibraryPropertiesDotALinkageFieldTrueWithFlatLayout() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("dot_a_linkage") {
		return checkresult.NotRun, ""
	}

	if checkdata.LoadedLibrary().DotALinkage && checkdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldLTMinLength checks if the library.properties "includes" value is less than the minimum length.
func LibraryPropertiesIncludesFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if !checkdata.LibraryProperties().ContainsKey("includes") {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("includes", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesIncludesFieldItemNotFound checks whether the header files specified in the library.properties `includes` field are in the library.
func LibraryPropertiesIncludesFieldItemNotFound() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	includes, ok := checkdata.LibraryProperties().GetOk("includes")
	if !ok {
		return checkresult.NotRun, ""
	}

	includesList, err := properties.SplitQuotedString(includes, "", false)
	if err != nil {
		panic(err)
	}

	findInclude := func(include string) bool {
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
		return checkresult.NotRun, ""
	}

	precompiled, ok := checkdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return checkresult.NotRun, ""
	}

	if schema.PropertyEnumMismatch("precompiled", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, precompiled
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout checks whether a precompiled library has the required recursive layout type.
func LibraryPropertiesPrecompiledFieldEnabledWithFlatLayout() (result checkresult.Type, output string) {
	if checkdata.LoadedLibrary() == nil || checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	precompiled, ok := checkdata.LibraryProperties().GetOk("precompiled")
	if !ok {
		return checkresult.NotRun, ""
	}

	if checkdata.LoadedLibrary().Precompiled && checkdata.LoadedLibrary().Layout == libraries.FlatLayout {
		return checkresult.Fail, precompiled
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesLdflagsFieldLTMinLength checks if the library.properties "ldflags" value is less than the minimum length.
func LibraryPropertiesLdflagsFieldLTMinLength() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.PropertyLessThanMinLength("ldflags", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesMisspelledOptionalField checks if library.properties contains common misspellings of optional fields.
func LibraryPropertiesMisspelledOptionalField() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.MisspelledOptionalPropertyFound(checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryInvalid checks whether the provided path is a valid library.
func LibraryInvalid() (result checkresult.Type, output string) {
	directoryListing, err := checkdata.LoadedLibrary().SourceDir.ReadDir()
	if err != nil {
		panic(err)
	}

	directoryListing.FilterOutDirs()
	for _, potentialHeaderFile := range directoryListing {
		if library.HasHeaderFileValidExtension(potentialHeaderFile) {
			return checkresult.Pass, ""
		}
	}

	return checkresult.Fail, ""
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

// ProhibitedCharactersInLibraryFolderName checks for prohibited characters in the library folder name.
func ProhibitedCharactersInLibraryFolderName() (result checkresult.Type, output string) {
	if !validProjectPathBaseName(checkdata.ProjectPath().Base()) {
		return checkresult.Fail, checkdata.ProjectPath().Base()
	}

	return checkresult.Pass, ""
}

// LibraryFolderNameGTMaxLength checks if the library folder name exceeds the maximum length.
func LibraryFolderNameGTMaxLength() (result checkresult.Type, output string) {
	if len(checkdata.ProjectPath().Base()) > 63 {
		return checkresult.Fail, checkdata.ProjectPath().Base()
	}

	return checkresult.Pass, ""
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

// spellCheckLibraryPropertiesFieldValue returns the value of the provided library.properties field with commonly misspelled words corrected.
func spellCheckLibraryPropertiesFieldValue(fieldName string) (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	fieldValue, ok := checkdata.LibraryProperties().GetOk(fieldName)
	if !ok {
		return checkresult.NotRun, ""
	}

	replaced, diff := checkdata.MisspelledWordsReplacer().Replace(fieldValue)
	if diff != nil {
		return checkresult.Fail, replaced
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
