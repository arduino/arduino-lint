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

// This file contains tests for the library.properties JSON schemas.
package libraryproperties_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/arduino/arduino-check/check/checkdata/schema"
	"github.com/arduino/arduino-check/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/ory/jsonschema/v3"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"

	"github.com/stretchr/testify/assert"
)

var validLibraryPropertiesMap = map[string]string{
	"name":          "WebServer",
	"version":       "1.0.0",
	"author":        "Cristian Maglie <c.maglie@example.com>, Pippo Pluto <pippo@example.com>",
	"maintainer":    "Cristian Maglie <c.maglie@example.com>",
	"sentence":      "A library that makes coding a Webserver a breeze.",
	"paragraph":     "Supports HTTP1.1 and you can do GET and POST.",
	"category":      "Communication",
	"url":           "http://example.com/",
	"architectures": "avr",
	"depends":       "ArduinoHttpClient",
	"dot_a_linkage": "true",
	"includes":      "WebServer.h",
	"precompiled":   "full",
	"ldflags":       "-lm",
}

var schemasPath *paths.Path

func init() {
	workingPath, _ := os.Getwd()
	schemasPath = paths.New(workingPath).Join("..", "..", "..", "etc", "schemas")
}

type propertyValueTestTable struct {
	testName        string
	propertyValue   string
	complianceLevel compliancelevel.Type
	assertion       assert.BoolAssertionFunc
}

func checkPropertyPatternMismatch(propertyName string, testTables []propertyValueTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyPatternMismatch(propertyName, validationResult[testTable.complianceLevel], schemasPath))
		})
	}
}

func checkPropertyEnumMismatch(propertyName string, testTables []propertyValueTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyEnumMismatch(propertyName, validationResult[testTable.complianceLevel], schemasPath))
		})
	}
}

type validationErrorTestTable struct {
	testName           string
	propertyValue      string
	schemaPointerQuery string
	complianceLevel    compliancelevel.Type
	assertion          assert.BoolAssertionFunc
}

func checkValidationErrorMatch(propertyName string, testTables []validationErrorTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.ValidationErrorMatch("#/"+propertyName, testTable.schemaPointerQuery, "", "", validationResult[testTable.complianceLevel], schemasPath))
		})
	}
}

func changeValueUpdateValidationResult(propertyName string, propertyValue string, libraryProperties *properties.Map, validationResult map[compliancelevel.Type]*jsonschema.ValidationError) map[compliancelevel.Type]*jsonschema.ValidationError {
	if validationResult == nil || libraryProperties.Get(propertyName) != propertyValue {
		libraryProperties.Set(propertyName, propertyValue)
		return libraryproperties.Validate(libraryProperties, schemasPath)
	}

	// No change to property, return the previous validationResult.
	return validationResult
}

func TestCompile(t *testing.T) {
	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Validate = true // Enable meta-schema validation when schemas are added and compiled

	directoryListing, _ := schemasPath.ReadDir()
	directoryListing.FilterOutDirs()
	directoryListing.FilterSuffix(".json")

	referencedSchemaFilenames := []string{}
	// Generate a list of referenced schemas
	logrus.Trace("Discovering definition schemas:")
	for _, schemaPath := range directoryListing {
		if schemaPath.HasSuffix("definitions-schema.json") {
			logrus.Trace(schemaPath)
			referencedSchemaFilenames = append(referencedSchemaFilenames, schemaPath.Base())
		}
	}

	// Compile the parent schemas
	logrus.Trace("Validating schemas:")
	for _, schemaPath := range directoryListing {
		if !schemaPath.HasSuffix("definitions-schema.json") {
			logrus.Trace(schemaPath)
			assert.NotPanics(t, func() {
				schema.Compile(schemaPath.Base(), referencedSchemaFilenames, schemasPath)
			})
		}
	}
}

func TestPropertiesValid(t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	validationResult := libraryproperties.Validate(libraryProperties, schemasPath)
	assert.Nil(t, validationResult[compliancelevel.Permissive])
	assert.Nil(t, validationResult[compliancelevel.Specification])
	assert.Nil(t, validationResult[compliancelevel.Strict])
}

func TestPropertiesMinLength(t *testing.T) {
	tests := []struct {
		propertyName    string
		minLength       int
		complianceLevel compliancelevel.Type
	}{
		{"name", 1, compliancelevel.Permissive},
		{"name", 1, compliancelevel.Specification},
		{"name", 1, compliancelevel.Strict},

		{"author", 1, compliancelevel.Permissive},
		{"author", 1, compliancelevel.Specification},
		{"author", 1, compliancelevel.Strict},

		{"maintainer", 1, compliancelevel.Permissive},
		{"maintainer", 1, compliancelevel.Specification},
		{"maintainer", 1, compliancelevel.Strict},

		{"email", 1, compliancelevel.Permissive},
		{"email", 1, compliancelevel.Specification},
		{"email", 1, compliancelevel.Strict},

		{"sentence", 1, compliancelevel.Permissive},
		{"sentence", 1, compliancelevel.Specification},
		{"sentence", 1, compliancelevel.Strict},

		{"architectures", 0, compliancelevel.Permissive},
		{"architectures", 1, compliancelevel.Specification},
		{"architectures", 1, compliancelevel.Strict},

		{"includes", 1, compliancelevel.Permissive},
		{"includes", 1, compliancelevel.Specification},
		{"includes", 1, compliancelevel.Strict},

		{"ldflags", 0, compliancelevel.Permissive},
		{"ldflags", 3, compliancelevel.Specification},
		{"ldflags", 3, compliancelevel.Strict},
	}

	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	// Test schema validation results with value length < minimum.
	for _, tt := range tests {
		var assertion assert.BoolAssertionFunc
		if tt.minLength == 0 {
			assertion = assert.False
		} else {
			assertion = assert.True
			value, propertyExists := libraryProperties.GetOk(tt.propertyName)
			if !propertyExists || len(value) >= tt.minLength {
				libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
				libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.minLength-1))
				validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
			}
		}

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", tt.propertyName, tt.minLength, tt.complianceLevel), func(t *testing.T) {
			assertion(t, schema.PropertyLessThanMinLength(tt.propertyName, validationResult[tt.complianceLevel], schemasPath))
		})
	}

	// Test schema validation results with minimum value length.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) < tt.minLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.minLength))
			validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
		}

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", tt.propertyName, tt.minLength, tt.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(tt.propertyName, validationResult[tt.complianceLevel], schemasPath))
		})
	}
}

func TestPropertiesMaxLength(t *testing.T) {
	tests := []struct {
		propertyName    string
		maxLength       int
		complianceLevel compliancelevel.Type
	}{
		{"name", 63, compliancelevel.Permissive},
		{"name", 63, compliancelevel.Specification},
		{"name", 16, compliancelevel.Strict},
	}

	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	// Test schema validation results with value length > maximum.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) <= tt.maxLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.maxLength+1))
			validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
		}

		t.Run(fmt.Sprintf("%s greater than maximum length of %d (%s)", tt.propertyName, tt.maxLength, tt.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyGreaterThanMaxLength(tt.propertyName, validationResult[tt.complianceLevel], schemasPath))
		})
	}

	// Test schema validation results with minimum value length.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) > tt.maxLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.maxLength))
			validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
		}

		t.Run(fmt.Sprintf("%s at maximum length of %d (%s)", tt.propertyName, tt.maxLength, tt.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyGreaterThanMaxLength(tt.propertyName, validationResult[tt.complianceLevel], schemasPath))
		})
	}
}

func TestPropertiesNamePattern(t *testing.T) {
	testTables := []validationErrorTestTable{
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Permissive, assert.True},
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Specification, assert.True},
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Strict, assert.True},

		{"Starts with arduino", "arduinofoo", "/patternObjects/notStartsWithArduino", compliancelevel.Permissive, assert.False},
		{"Starts with arduino", "arduinofoo", "/patternObjects/notStartsWithArduino", compliancelevel.Specification, assert.True},
		{"Starts with arduino", "arduinofoo", "/patternObjects/notStartsWithArduino", compliancelevel.Strict, assert.True},

		{"Contains spaces", "foo bar", "/patternObjects/notContainsSpaces", compliancelevel.Permissive, assert.False},
		{"Contains spaces", "foo bar", "/patternObjects/notContainsSpaces", compliancelevel.Specification, assert.False},
		{"Contains spaces", "foo bar", "/patternObjects/notContainsSpaces", compliancelevel.Strict, assert.True},

		{"Contains superfluous terms", "foo library", "/patternObjects/notContainsSuperfluousTerms", compliancelevel.Permissive, assert.False},
		{"Contains superfluous terms", "foolibrary", "/patternObjects/notContainsSuperfluousTerms", compliancelevel.Specification, assert.False},
		{"Contains superfluous terms", "foolibrary", "/patternObjects/notContainsSuperfluousTerms", compliancelevel.Strict, assert.True},
	}

	checkValidationErrorMatch("name", testTables, t)
}

func TestPropertiesVersionPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"X.Y.Z-prerelease", "1.0.0-pre1", compliancelevel.Permissive, assert.False},
		{"X.Y.Z-prerelease", "1.0.0-pre1", compliancelevel.Specification, assert.False},
		{"X.Y.Z-prerelease", "1.0.0-pre1", compliancelevel.Strict, assert.False},

		{"X.Y.Z+build", "1.0.0+build", compliancelevel.Permissive, assert.False},
		{"X.Y.Z+build", "1.0.0+build", compliancelevel.Specification, assert.False},
		{"X.Y.Z+build", "1.0.0+build", compliancelevel.Strict, assert.False},

		{"vX.Y.Z", "v1.0.0", compliancelevel.Permissive, assert.True},
		{"vX.Y.Z", "v1.0.0", compliancelevel.Specification, assert.True},
		{"vX.Y.Z", "v1.0.0", compliancelevel.Strict, assert.True},

		{"X.Y", "1.0", compliancelevel.Permissive, assert.False},
		{"X.Y", "1.0", compliancelevel.Specification, assert.True},
		{"X.Y", "1.0", compliancelevel.Strict, assert.True},

		{"X", "1", compliancelevel.Permissive, assert.False},
		{"X", "1", compliancelevel.Specification, assert.True},
		{"X", "1", compliancelevel.Strict, assert.True},
	}

	checkPropertyPatternMismatch("version", testTables, t)
}

func TestPropertiesMaintainerPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Starts with arduino", "arduinofoo", compliancelevel.Permissive, assert.False},
		{"Starts with arduino", "arduinofoo", compliancelevel.Specification, assert.True},
		{"Starts with arduino", "arduinofoo", compliancelevel.Strict, assert.True},
	}

	checkPropertyPatternMismatch("maintainer", testTables, t)
}

func TestPropertiesEmailPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Starts with arduino", "arduinofoo", compliancelevel.Permissive, assert.False},
		{"Starts with arduino", "arduinofoo", compliancelevel.Specification, assert.True},
		{"Starts with arduino", "arduinofoo", compliancelevel.Strict, assert.True},
	}

	checkPropertyPatternMismatch("email", testTables, t)
}

func TestPropertiesCategoryEnum(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Invalid category", "foo", compliancelevel.Permissive, assert.False},
		{"Invalid category", "foo", compliancelevel.Specification, assert.True},
		{"Invalid category", "foo", compliancelevel.Strict, assert.True},
	}

	checkPropertyEnumMismatch("category", testTables, t)
}

func TestPropertiesUrlFormat(t *testing.T) {
	testTables := []validationErrorTestTable{
		{"Invalid URL format", "foo", "/format$", compliancelevel.Permissive, assert.False},
		{"Invalid URL format", "foo", "/format$", compliancelevel.Specification, assert.True},
		{"Invalid URL format", "foo", "/format$", compliancelevel.Strict, assert.True},
	}

	checkValidationErrorMatch("url", testTables, t)
}

func TestPropertiesDependsPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Invalid characters", "-foo", compliancelevel.Permissive, assert.True},
		{"Invalid characters", "-foo", compliancelevel.Permissive, assert.True},
		{"Invalid characters", "-foo", compliancelevel.Permissive, assert.True},
	}

	checkPropertyPatternMismatch("depends", testTables, t)
}

func TestPropertiesDotALinkageEnum(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Invalid enum value", "foo", compliancelevel.Permissive, assert.True},
		{"Invalid enum value", "foo", compliancelevel.Specification, assert.True},
		{"Invalid enum value", "foo", compliancelevel.Strict, assert.True},
	}

	checkPropertyEnumMismatch("dot_a_linkage", testTables, t)
}

func TestPropertiesPrecompiledEnum(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Invalid enum value", "foo", compliancelevel.Permissive, assert.True},
		{"Invalid enum value", "foo", compliancelevel.Specification, assert.True},
		{"Invalid enum value", "foo", compliancelevel.Strict, assert.True},
	}

	checkPropertyEnumMismatch("precompiled", testTables, t)
}

func TestPropertyNames(t *testing.T) {
	testTables := []struct {
		testName           string
		removePropertyName string
		addPropertyName    string
		complianceLevel    compliancelevel.Type
		assertion          assert.BoolAssertionFunc
	}{
		{"depends: singular", "depends", "depend", compliancelevel.Permissive, assert.False},
		{"depends: singular", "depends", "depend", compliancelevel.Specification, assert.False},
		{"depends: singular", "depends", "depend", compliancelevel.Strict, assert.True},

		{"depends: miscapitalized", "depends", "Depends", compliancelevel.Permissive, assert.False},
		{"depends: miscapitalized", "depends", "Depends", compliancelevel.Specification, assert.False},
		{"depends: miscapitalized", "depends", "Depends", compliancelevel.Strict, assert.True},

		{"dot_a_linkage: mis-hyphenated", "dot_a_linkage", "dot-a-linkage", compliancelevel.Permissive, assert.False},
		{"dot_a_linkage: mis-hyphenated", "dot_a_linkage", "dot-a-linkage", compliancelevel.Specification, assert.False},
		{"dot_a_linkage: mis-hyphenated", "dot_a_linkage", "dot-a-linkage", compliancelevel.Strict, assert.True},

		{"dot_a_linkage: plural", "dot_a_linkage", "dot_a_linkages", compliancelevel.Permissive, assert.False},
		{"dot_a_linkage: plural", "dot_a_linkage", "dot_a_linkages", compliancelevel.Specification, assert.False},
		{"dot_a_linkage: plural", "dot_a_linkage", "dot_a_linkages", compliancelevel.Strict, assert.True},

		{"dot_a_linkage: miscapitalized", "dot_a_linkage", "Dot_a_linkage", compliancelevel.Permissive, assert.False},
		{"dot_a_linkage: miscapitalized", "dot_a_linkage", "Dot_a_linkage", compliancelevel.Specification, assert.False},
		{"dot_a_linkage: miscapitalized", "dot_a_linkage", "Dot_a_linkage", compliancelevel.Strict, assert.True},

		{"includes: singular", "includes", "include", compliancelevel.Permissive, assert.False},
		{"includes: singular", "includes", "include", compliancelevel.Specification, assert.False},
		{"includes: singular", "includes", "include", compliancelevel.Strict, assert.True},

		{"includes: miscapitalized", "includes", "Includes", compliancelevel.Permissive, assert.False},
		{"includes: miscapitalized", "includes", "Includes", compliancelevel.Specification, assert.False},
		{"includes: miscapitalized", "includes", "Includes", compliancelevel.Strict, assert.True},

		{"precompiled: tense", "precompiled", "precompile", compliancelevel.Permissive, assert.False},
		{"precompiled: tense", "precompiled", "precompile", compliancelevel.Specification, assert.False},
		{"precompiled: tense", "precompiled", "precompile", compliancelevel.Strict, assert.True},

		{"precompiled: mis-hyphenated", "precompiled", "pre-compiled", compliancelevel.Permissive, assert.False},
		{"precompiled: mis-hyphenated", "precompiled", "pre-compiled", compliancelevel.Specification, assert.False},
		{"precompiled: mis-hyphenated", "precompiled", "pre-compiled", compliancelevel.Strict, assert.True},

		{"precompiled: miscapitalized", "precompiled", "Precompiled", compliancelevel.Permissive, assert.False},
		{"precompiled: miscapitalized", "precompiled", "Precompiled", compliancelevel.Specification, assert.False},
		{"precompiled: miscapitalized", "precompiled", "Precompiled", compliancelevel.Strict, assert.True},

		{"ldflags: mis-hyphenated", "ldflags", "ld_flags", compliancelevel.Permissive, assert.False},
		{"ldflags: mis-hyphenated", "ldflags", "ld_flags", compliancelevel.Specification, assert.False},
		{"ldflags: mis-hyphenated", "ldflags", "ld_flags", compliancelevel.Strict, assert.True},

		{"ldflags: singular", "ldflags", "ldflag", compliancelevel.Permissive, assert.False},
		{"ldflags: singular", "ldflags", "ldflag", compliancelevel.Specification, assert.False},
		{"ldflags: singular", "ldflags", "ldflag", compliancelevel.Strict, assert.True},
	}

	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	for _, testTable := range testTables {
		_, removePropertyPresent := libraryProperties.GetOk(testTable.removePropertyName)
		_, addPropertyPresent := libraryProperties.GetOk(testTable.addPropertyName)
		if validationResult == nil || removePropertyPresent || !addPropertyPresent {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Remove(testTable.removePropertyName)
			libraryProperties.Set(testTable.addPropertyName, "foo")
			validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
		}

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.MisspelledOptionalPropertyFound(validationResult[testTable.complianceLevel], schemasPath))
		})
	}
}

func TestRequired(t *testing.T) {
	testTables := []struct {
		propertyName    string
		complianceLevel compliancelevel.Type
		assertion       assert.BoolAssertionFunc
	}{
		{"name", compliancelevel.Permissive, assert.True},
		{"name", compliancelevel.Specification, assert.True},
		{"name", compliancelevel.Strict, assert.True},

		{"version", compliancelevel.Permissive, assert.True},
		{"version", compliancelevel.Specification, assert.True},
		{"version", compliancelevel.Strict, assert.True},

		{"author", compliancelevel.Permissive, assert.True},
		{"author", compliancelevel.Specification, assert.True},
		{"author", compliancelevel.Strict, assert.True},

		{"maintainer", compliancelevel.Permissive, assert.True},
		{"maintainer", compliancelevel.Specification, assert.True},
		{"maintainer", compliancelevel.Strict, assert.True},

		{"sentence", compliancelevel.Permissive, assert.True},
		{"sentence", compliancelevel.Specification, assert.True},
		{"sentence", compliancelevel.Strict, assert.True},

		{"paragraph", compliancelevel.Permissive, assert.True},
		{"paragraph", compliancelevel.Specification, assert.True},
		{"paragraph", compliancelevel.Strict, assert.True},

		{"category", compliancelevel.Permissive, assert.False},
		{"category", compliancelevel.Specification, assert.False},
		{"category", compliancelevel.Strict, assert.True},

		{"url", compliancelevel.Permissive, assert.True},
		{"url", compliancelevel.Specification, assert.True},
		{"url", compliancelevel.Strict, assert.True},

		{"architectures", compliancelevel.Permissive, assert.False},
		{"architectures", compliancelevel.Specification, assert.False},
		{"architectures", compliancelevel.Strict, assert.True},

		{"depends", compliancelevel.Permissive, assert.False},
		{"depends", compliancelevel.Specification, assert.False},
		{"depends", compliancelevel.Strict, assert.False},

		{"dot_a_linkage", compliancelevel.Permissive, assert.False},
		{"dot_a_linkage", compliancelevel.Specification, assert.False},
		{"dot_a_linkage", compliancelevel.Strict, assert.False},

		{"includes", compliancelevel.Permissive, assert.False},
		{"includes", compliancelevel.Specification, assert.False},
		{"includes", compliancelevel.Strict, assert.False},

		{"precompiled", compliancelevel.Permissive, assert.False},
		{"precompiled", compliancelevel.Specification, assert.False},
		{"precompiled", compliancelevel.Strict, assert.False},

		{"ldflags", compliancelevel.Permissive, assert.False},
		{"ldflags", compliancelevel.Specification, assert.False},
		{"ldflags", compliancelevel.Strict, assert.False},
	}

	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]*jsonschema.ValidationError

	for _, testTable := range testTables {
		_, propertyExists := libraryProperties.GetOk(testTable.propertyName)
		if propertyExists {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Remove(testTable.propertyName)
			validationResult = libraryproperties.Validate(libraryProperties, schemasPath)
		}

		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.RequiredPropertyMissing(testTable.propertyName, validationResult[testTable.complianceLevel], schemasPath))
		})
	}
}

func TestPropertiesMaintainerOrEmailRequired(t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	libraryProperties.Remove("maintainer")
	libraryProperties.Set("email", "foo@example.com")
	validationResult := libraryproperties.Validate(libraryProperties, schemasPath)
	assert.False(
		t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Permissive], schemasPath),
		"maintainer property is not required when email property is defined.",
	)
	assert.True(
		t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Specification], schemasPath),
		"maintainer property is unconditionally required.",
	)
	assert.True(t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Strict], schemasPath),
		"maintainer property is unconditionally required.",
	)
}
