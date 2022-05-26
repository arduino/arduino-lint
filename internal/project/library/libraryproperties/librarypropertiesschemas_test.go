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

// This file contains tests for the library.properties JSON schemas.
package libraryproperties_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/library/libraryproperties"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-properties-orderedmap"

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

type propertyValueTestTable struct {
	testName        string
	propertyValue   string
	complianceLevel compliancelevel.Type
	assertion       assert.BoolAssertionFunc
}

func checkPropertyPatternMismatch(propertyName string, testTables []propertyValueTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyPatternMismatch(propertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func checkPropertyEnumMismatch(propertyName string, testTables []propertyValueTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyEnumMismatch(propertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func checkPropertyFormatMismatch(propertyName string, testTables []propertyValueTestTable, t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyFormatMismatch(propertyName, validationResult[testTable.complianceLevel]))
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
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		validationResult = changeValueUpdateValidationResult(propertyName, testTable.propertyValue, libraryProperties, validationResult)

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.ValidationErrorMatch("#/"+propertyName, testTable.schemaPointerQuery, "", "", validationResult[testTable.complianceLevel]))
		})
	}
}

func changeValueUpdateValidationResult(propertyName string, propertyValue string, libraryProperties *properties.Map, validationResult map[compliancelevel.Type]schema.ValidationResult) map[compliancelevel.Type]schema.ValidationResult {
	if validationResult == nil || libraryProperties.Get(propertyName) != propertyValue {
		libraryProperties.Set(propertyName, propertyValue)
		return libraryproperties.Validate(libraryProperties)
	}

	// No change to property, return the previous validationResult.
	return validationResult
}

func TestPropertiesValid(t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	validationResult := libraryproperties.Validate(libraryProperties)
	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)
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
	var validationResult map[compliancelevel.Type]schema.ValidationResult

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
				validationResult = libraryproperties.Validate(libraryProperties)
			}
		}

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", tt.propertyName, tt.minLength, tt.complianceLevel), func(t *testing.T) {
			assertion(t, schema.PropertyLessThanMinLength(tt.propertyName, validationResult[tt.complianceLevel]))
		})
	}

	// Test schema validation results with minimum value length.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) < tt.minLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.minLength))
			validationResult = libraryproperties.Validate(libraryProperties)
		}

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", tt.propertyName, tt.minLength, tt.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(tt.propertyName, validationResult[tt.complianceLevel]))
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
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	// Test schema validation results with value length > maximum.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) <= tt.maxLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.maxLength+1))
			validationResult = libraryproperties.Validate(libraryProperties)
		}

		t.Run(fmt.Sprintf("%s greater than maximum length of %d (%s)", tt.propertyName, tt.maxLength, tt.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyGreaterThanMaxLength(tt.propertyName, validationResult[tt.complianceLevel]))
		})
	}

	// Test schema validation results with minimum value length.
	for _, tt := range tests {
		if len(libraryProperties.Get(tt.propertyName)) > tt.maxLength {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Set(tt.propertyName, strings.Repeat("a", tt.maxLength))
			validationResult = libraryproperties.Validate(libraryProperties)
		}

		t.Run(fmt.Sprintf("%s at maximum length of %d (%s)", tt.propertyName, tt.maxLength, tt.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyGreaterThanMaxLength(tt.propertyName, validationResult[tt.complianceLevel]))
		})
	}
}

func TestPropertiesNamePattern(t *testing.T) {
	testTables := []validationErrorTestTable{
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Permissive, assert.True},
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Specification, assert.True},
		{"Disallowed character", "-foo", "/patternObjects/allowedCharacters", compliancelevel.Strict, assert.True},

		// The "minLength" schema will enforce the minimum length, so this is not the responsibility of the pattern schema.
		{"Empty", "", "/patternObjects/allowedCharacters", compliancelevel.Permissive, assert.False},
		{"Empty", "", "/patternObjects/allowedCharacters", compliancelevel.Specification, assert.False},
		{"Empty", "", "/patternObjects/allowedCharacters", compliancelevel.Strict, assert.False},

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
		{"X.Y", "1.0", compliancelevel.Specification, assert.False},
		{"X.Y", "1.0", compliancelevel.Strict, assert.True},

		{"X", "1", compliancelevel.Permissive, assert.False},
		{"X", "1", compliancelevel.Specification, assert.False},
		{"X", "1", compliancelevel.Strict, assert.True},
	}

	checkPropertyPatternMismatch("version", testTables, t)
}

func TestPropertiesMaintainerPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Starts with arduino", "arduinofoo", compliancelevel.Permissive, assert.False},
		{"Contains arduino", "fooarduinobar", compliancelevel.Permissive, assert.False},
		{"Starts with arduino", "arduinofoo", compliancelevel.Specification, assert.True},
		{"Contains arduino", "fooarduinobar", compliancelevel.Specification, assert.False},
		{"Starts with arduino", "arduinofoo", compliancelevel.Strict, assert.True},
		{"Contains arduino", "fooarduinobar", compliancelevel.Strict, assert.True},
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
	testTables := []propertyValueTestTable{
		{"Invalid URL format", "foo", compliancelevel.Permissive, assert.False},
		{"Invalid URL format", "foo", compliancelevel.Specification, assert.True},
		{"Invalid URL format", "foo", compliancelevel.Strict, assert.True},
	}

	checkPropertyFormatMismatch("url", testTables, t)
}

func TestPropertiesDependsPattern(t *testing.T) {
	testTables := []propertyValueTestTable{
		{"Valid name", "foo", compliancelevel.Permissive, assert.False},
		{"Valid name", "foo", compliancelevel.Specification, assert.False},
		{"Valid name", "foo", compliancelevel.Strict, assert.False},

		{"Valid names", "foo,bar", compliancelevel.Permissive, assert.False},
		{"Valid names", "foo,bar", compliancelevel.Specification, assert.False},
		{"Valid names", "foo,bar", compliancelevel.Strict, assert.False},

		{"Trailing comma", "foo,", compliancelevel.Permissive, assert.True},
		{"Trailing comma", "foo,", compliancelevel.Specification, assert.True},
		{"Trailing comma", "foo,", compliancelevel.Strict, assert.True},

		{"Invalid characters", "-foo", compliancelevel.Permissive, assert.True},
		{"Invalid characters", "-foo", compliancelevel.Specification, assert.True},
		{"Invalid characters", "-foo", compliancelevel.Strict, assert.True},

		{"Empty", "", compliancelevel.Permissive, assert.False},
		{"Empty", "", compliancelevel.Specification, assert.False},
		{"Empty", "", compliancelevel.Strict, assert.False},

		{"<version", "foo (<1.2.3)", compliancelevel.Permissive, assert.False},
		{"<version", "foo (<1.2.3)", compliancelevel.Specification, assert.False},
		{"<version", "foo (<1.2.3)", compliancelevel.Strict, assert.False},
		{"<=version", "foo (<=1.2.3)", compliancelevel.Permissive, assert.False},
		{"<=version", "foo (<=1.2.3)", compliancelevel.Specification, assert.False},
		{"<=version", "foo (<=1.2.3)", compliancelevel.Strict, assert.False},
		{"=version", "foo (=1.2.3)", compliancelevel.Permissive, assert.False},
		{"=version", "foo (=1.2.3)", compliancelevel.Specification, assert.False},
		{"=version", "foo (=1.2.3)", compliancelevel.Strict, assert.False},
		{">=version", "foo (>=1.2.3)", compliancelevel.Permissive, assert.False},
		{">=version", "foo (>=1.2.3)", compliancelevel.Specification, assert.False},
		{">=version", "foo (>=1.2.3)", compliancelevel.Strict, assert.False},
		{">version", "foo (>1.2.3)", compliancelevel.Permissive, assert.False},
		{">version", "foo (>1.2.3)", compliancelevel.Specification, assert.False},
		{">version", "foo (>1.2.3)", compliancelevel.Strict, assert.False},

		{"Relaxed version", "foo (=1.2)", compliancelevel.Permissive, assert.False},
		{"Relaxed version", "foo (=1.2)", compliancelevel.Specification, assert.False},
		{"Relaxed version", "foo (=1.2)", compliancelevel.Strict, assert.False},
		{"Pre-release version", "foo (=1.2.3-rc1)", compliancelevel.Permissive, assert.False},
		{"Pre-release version", "foo (=1.2.3-rc1)", compliancelevel.Specification, assert.False},
		{"Pre-release version", "foo (=1.2.3-rc1)", compliancelevel.Strict, assert.False},

		{"Version w/o space", "foo(>1.2.3)", compliancelevel.Permissive, assert.True},
		{"Version w/o space", "foo(>1.2.3)", compliancelevel.Specification, assert.True},
		{"Version w/o space", "foo(>1.2.3)", compliancelevel.Strict, assert.True},

		{"Names w/ version", "foo (<=1.2.3),bar", compliancelevel.Permissive, assert.False},
		{"Names w/ version", "foo (<=1.2.3),bar", compliancelevel.Specification, assert.False},
		{"Names w/ version", "foo (<=1.2.3),bar", compliancelevel.Strict, assert.False},

		{"Names w/ parenthesized version constraints", "foo ((>0.1.0 && <2.0.0) || >2.1.0),bar", compliancelevel.Permissive, assert.False},
		{"Names w/ parenthesized version constraints", "foo ((>0.1.0 && <2.0.0) || >2.1.0),bar", compliancelevel.Specification, assert.False},
		{"Names w/ parenthesized version constraints", "foo ((>0.1.0 && <2.0.0) || >2.1.0),bar", compliancelevel.Strict, assert.False},

		{"Names w/ empty version constraint", "foo (),bar", compliancelevel.Permissive, assert.False},
		{"Names w/ empty version constraint", "foo (),bar", compliancelevel.Specification, assert.False},
		{"Names w/ empty version constraint", "foo (),bar", compliancelevel.Strict, assert.False},
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
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		_, removePropertyPresent := libraryProperties.GetOk(testTable.removePropertyName)
		_, addPropertyPresent := libraryProperties.GetOk(testTable.addPropertyName)
		if validationResult == nil || removePropertyPresent || !addPropertyPresent {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Remove(testTable.removePropertyName)
			libraryProperties.Set(testTable.addPropertyName, "foo")
			validationResult = libraryproperties.Validate(libraryProperties)
		}

		t.Run(fmt.Sprintf("%s (%s)", testTable.testName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.MisspelledOptionalPropertyFound(validationResult[testTable.complianceLevel]))
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
	var validationResult map[compliancelevel.Type]schema.ValidationResult

	for _, testTable := range testTables {
		_, propertyExists := libraryProperties.GetOk(testTable.propertyName)
		if propertyExists {
			libraryProperties = properties.NewFromHashmap(validLibraryPropertiesMap)
			libraryProperties.Remove(testTable.propertyName)
			validationResult = libraryproperties.Validate(libraryProperties)
		}

		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.RequiredPropertyMissing(testTable.propertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestPropertiesMaintainerOrEmailRequired(t *testing.T) {
	libraryProperties := properties.NewFromHashmap(validLibraryPropertiesMap)
	libraryProperties.Remove("maintainer")
	libraryProperties.Set("email", "foo@example.com")
	validationResult := libraryproperties.Validate(libraryProperties)
	assert.False(
		t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Permissive]),
		"maintainer property is not required when email property is defined.",
	)
	assert.True(
		t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Specification]),
		"maintainer property is unconditionally required.",
	)
	assert.True(t,
		schema.RequiredPropertyMissing("maintainer", validationResult[compliancelevel.Strict]),
		"maintainer property is unconditionally required.",
	)
}
