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
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkdata/schema"
	"github.com/arduino/arduino-check/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
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
