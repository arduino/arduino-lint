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

package schema

import (
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/general"
	"github.com/arduino/arduino-lint/internal/rule/schema/testdata"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/ory/jsonschema/v3"
	"github.com/stretchr/testify/require"
)

var validMap = map[string]string{
	"property1":          "foo",
	"property2":          "bar",
	"property3":          "baz",
	"dependentProperty":  "asdf",
	"dependencyProperty": "zxcv",
}

var validPropertiesMap = properties.NewFromHashmap(validMap)

var validSchemaWithReferences Schema

func init() {
	validSchemaWithReferences = Compile(
		"valid-schema-with-references.json",
		[]string{
			"referenced-schema-1.json",
			"referenced-schema-2.json",
		},
		testdata.Asset,
	)
}

func TestCompile(t *testing.T) {
	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{"nonexistent.json"}, testdata.Asset)
	})

	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{"schema-without-id.json"}, testdata.Asset)
	})

	require.Panics(t, func() {
		Compile("invalid-schema.json", []string{}, testdata.Asset)
	})

	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{}, testdata.Asset)
	})

	require.NotPanics(t, func() {
		Compile("valid-schema.json", []string{}, testdata.Asset)
	})

	require.NotPanics(t, func() {
		Compile(
			"valid-schema-with-references.json",
			[]string{
				"referenced-schema-1.json",
				"referenced-schema-2.json",
			},
			testdata.Asset,
		)
	})
}

func TestValidate(t *testing.T) {
	schemaObject := Compile("valid-schema.json", []string{}, testdata.Asset)
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), schemaObject)
	require.Nil(t, validationResult.Result)

	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.Nil(t, validationResult.Result)

	propertiesMap.Set("property1", "a")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), schemaObject)
	require.Equal(t, "#/property1", validationResult.Result.InstancePtr)
	require.Equal(t, "#/properties/property1/minLength", validationResult.Result.SchemaPtr)
}

func TestRequiredPropertyMissing(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, RequiredPropertyMissing("property1", validationResult))

	propertiesMap.Remove("property1")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, RequiredPropertyMissing("property1", validationResult))
}

func TestPropertyPatternMismatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, PropertyPatternMismatch("property2", validationResult))

	propertiesMap.Set("property2", "fOo")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, PropertyPatternMismatch("property2", validationResult))

	require.False(t, PropertyPatternMismatch("property1", validationResult))
}

func TestPropertyLessThanMinLength(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, PropertyLessThanMinLength("property1", validationResult))

	propertiesMap.Set("property1", "a")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, PropertyLessThanMinLength("property1", validationResult))
}

func TestPropertyGreaterThanMaxLength(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, PropertyGreaterThanMaxLength("property1", validationResult))

	propertiesMap.Set("property1", "12345")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, PropertyGreaterThanMaxLength("property1", validationResult))
}

func TestPropertyEnumMismatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, PropertyEnumMismatch("property3", validationResult))

	propertiesMap.Set("property3", "invalid")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, PropertyEnumMismatch("property3", validationResult))
}

func TestPropertyDependenciesMissing(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, PropertyDependenciesMissing("dependentProperty", validationResult))

	propertiesMap.Remove("dependencyProperty")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, PropertyDependenciesMissing("dependentProperty", validationResult))
}

func TestMisspelledOptionalPropertyFound(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, MisspelledOptionalPropertyFound(validationResult))

	propertiesMap.Set("porperties", "foo")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, MisspelledOptionalPropertyFound(validationResult))
}

func TestValidationErrorMatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, ValidationErrorMatch("", "", "", "", validationResult))

	propertiesMap.Set("property2", "fOo")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, ValidationErrorMatch("nomatch", "nomatch", "nomatch", "nomatch", validationResult))
	require.False(t, ValidationErrorMatch("^#/property2$", "nomatch", "nomatch", "nomatch", validationResult))
	require.False(t, ValidationErrorMatch("^#/property2$", "/pattern$", "nomatch", "nomatch", validationResult))
	require.False(t, ValidationErrorMatch("^#/property2$", "/pattern$", `^\^\[a-z\]\+\$$`, "nomatch", validationResult))
	require.True(t, ValidationErrorMatch("^#/property2$", "/pattern$", `^"\^\[a-z\]\+\$"$`, "", validationResult))
	require.True(t, ValidationErrorMatch("", "", "", "", validationResult))

	propertiesMap.Set("property3", "bAz")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.True(t, ValidationErrorMatch("^#/property3$", "/pattern$", "", "", validationResult), "Match pointer below logic inversion keyword")

	propertiesMap = properties.NewFromHashmap(validMap)
	propertiesMap.Remove("property1")
	validationResult = Validate(general.PropertiesToMap(propertiesMap, 0), validSchemaWithReferences)
	require.False(t, ValidationErrorMatch("nomatch", "nomatch", "nomatch", "nomatch", validationResult))
	require.True(t, ValidationErrorMatch("", "", "", "^#/property1$", validationResult))
}

func Test_loadReferencedSchema(t *testing.T) {
	compiler := jsonschema.NewCompiler()

	require.Panics(
		t,
		func() {
			loadReferencedSchema(compiler, "nonexistent.json", testdata.Asset)
		},
	)
	require.Error(t, loadReferencedSchema(compiler, "schema-without-id.json", testdata.Asset))
	require.Nil(t, loadReferencedSchema(compiler, "referenced-schema-2.json", testdata.Asset))
}

func Test_schemaID(t *testing.T) {
	_, err := schemaID("schema-without-id.json", testdata.Asset)
	require.NotNil(t, err)

	id, err := schemaID("valid-schema.json", testdata.Asset)
	require.Equal(t, "https://raw.githubusercontent.com/arduino/arduino-lint/main/internal/rule/schema/testdata/schema-with-references.json", id)
	require.Nil(t, err)
}

func Test_validationErrorSchemaPointerValue(t *testing.T) {
	validationError := ValidationResult{
		Result: &jsonschema.ValidationError{
			SchemaURL: "https://raw.githubusercontent.com/arduino/arduino-lint/main/internal/rule/schema/testdata/referenced-schema-1.json",
			SchemaPtr: "#/definitions/patternObject/pattern",
		},
		dataLoader: testdata.Asset,
	}

	schemaPointerValueInterface := validationErrorSchemaPointerValue(validationError)
	schemaPointerValue, ok := schemaPointerValueInterface.(string)
	require.True(t, ok)
	require.Equal(t, "^[a-z]+$", schemaPointerValue)
}

func Test_validationErrorContextMatch(t *testing.T) {
	validationError := jsonschema.ValidationError{
		Context: nil,
	}

	require.True(t, validationErrorContextMatch(regexp.MustCompile(".*"), &validationError))
	require.False(t, validationErrorContextMatch(regexp.MustCompile("foo"), &validationError))

	validationError.Context = &jsonschema.ValidationErrorContextRequired{
		Missing: []string{
			"#/foo",
			"#/bar",
		},
	}

	require.True(t, validationErrorContextMatch(regexp.MustCompile("^#/bar$"), &validationError))
	require.False(t, validationErrorContextMatch(regexp.MustCompile("nomatch"), &validationError))
}
