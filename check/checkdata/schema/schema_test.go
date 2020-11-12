package schema

import (
	"os"
	"regexp"
	"runtime"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/ory/jsonschema/v3"
	"github.com/stretchr/testify/require"
)

var schemasPath *paths.Path

var validMap = map[string]string{
	"property1": "foo",
	"property2": "bar",
	"property3": "baz",
}

var validPropertiesMap = properties.NewFromHashmap(validMap)

var validSchemaWithReferences *jsonschema.Schema

func init() {
	workingPath, _ := os.Getwd()
	schemasPath = paths.New(workingPath).Join("testdata")

	validSchemaWithReferences = Compile(
		"valid-schema-with-references.json",
		[]string{
			"referenced-schema-1.json",
			"referenced-schema-2.json",
		},
		schemasPath,
	)
}

func TestCompile(t *testing.T) {
	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{"nonexistent.json"}, schemasPath)
	})

	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{"schema-without-id.json"}, schemasPath)
	})

	require.Panics(t, func() {
		Compile("invalid-schema.json", []string{}, schemasPath)
	})

	require.Panics(t, func() {
		Compile("valid-schema-with-references.json", []string{}, schemasPath)
	})

	require.NotPanics(t, func() {
		Compile("valid-schema.json", []string{}, schemasPath)
	})

	require.NotPanics(t, func() {
		Compile(
			"valid-schema-with-references.json",
			[]string{
				"referenced-schema-1.json",
				"referenced-schema-2.json",
			},
			schemasPath,
		)
	})
}

func TestValidate(t *testing.T) {
	schemaObject := Compile("valid-schema.json", []string{}, schemasPath)
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, schemaObject, schemasPath)
	require.Nil(t, validationResult)

	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.Nil(t, validationResult)

	propertiesMap.Set("property1", "a")
	validationResult = Validate(propertiesMap, schemaObject, schemasPath)
	require.Equal(t, "#/property1", validationResult.InstancePtr)
	require.Equal(t, "#/properties/property1/minLength", validationResult.SchemaPtr)
}

func TestRequiredPropertyMissing(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, RequiredPropertyMissing("property1", validationResult, schemasPath))

	propertiesMap.Remove("property1")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, RequiredPropertyMissing("property1", validationResult, schemasPath))
}

func TestPropertyPatternMismatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, PropertyPatternMismatch("property2", validationResult, schemasPath))

	propertiesMap.Set("property2", "fOo")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, PropertyPatternMismatch("property2", validationResult, schemasPath))

	require.False(t, PropertyPatternMismatch("property1", validationResult, schemasPath))
}

func TestPropertyLessThanMinLength(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, PropertyLessThanMinLength("property1", validationResult, schemasPath))

	propertiesMap.Set("property1", "a")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, PropertyLessThanMinLength("property1", validationResult, schemasPath))
}

func TestPropertyGreaterThanMaxLength(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, PropertyGreaterThanMaxLength("property1", validationResult, schemasPath))

	propertiesMap.Set("property1", "12345")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, PropertyGreaterThanMaxLength("property1", validationResult, schemasPath))
}

func TestPropertyEnumMismatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, PropertyEnumMismatch("property3", validationResult, schemasPath))

	propertiesMap.Set("property3", "invalid")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, PropertyEnumMismatch("property3", validationResult, schemasPath))
}

func TestMisspelledOptionalPropertyFound(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, MisspelledOptionalPropertyFound(validationResult, schemasPath))

	propertiesMap.Set("porperties", "foo")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.True(t, MisspelledOptionalPropertyFound(validationResult, schemasPath))
}

func TestValidationErrorMatch(t *testing.T) {
	propertiesMap := properties.NewFromHashmap(validMap)
	validationResult := Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, ValidationErrorMatch("", "", "", "", validationResult, schemasPath))

	propertiesMap.Set("property2", "fOo")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, ValidationErrorMatch("nomatch", "nomatch", "nomatch", "nomatch", validationResult, schemasPath))
	require.False(t, ValidationErrorMatch("^#/property2$", "nomatch", "nomatch", "nomatch", validationResult, schemasPath))
	require.False(t, ValidationErrorMatch("^#/property2$", "/pattern$", "nomatch", "nomatch", validationResult, schemasPath))
	require.False(t, ValidationErrorMatch("^#/property2$", "/pattern$", `^\^\[a-z\]\+\$$`, "nomatch", validationResult, schemasPath))
	require.True(t, ValidationErrorMatch("^#/property2$", "/pattern$", `^"\^\[a-z\]\+\$"$`, "", validationResult, schemasPath))
	require.True(t, ValidationErrorMatch("", "", "", "", validationResult, schemasPath))

	propertiesMap = properties.NewFromHashmap(validMap)
	propertiesMap.Remove("property1")
	validationResult = Validate(propertiesMap, validSchemaWithReferences, schemasPath)
	require.False(t, ValidationErrorMatch("nomatch", "nomatch", "nomatch", "nomatch", validationResult, schemasPath))
	require.True(t, ValidationErrorMatch("", "", "", "^#/property1$", validationResult, schemasPath))
}

func Test_loadReferencedSchema(t *testing.T) {
	compiler := jsonschema.NewCompiler()

	require.Error(t, loadReferencedSchema(compiler, "nonexistent.json", schemasPath))
	require.Error(t, loadReferencedSchema(compiler, "schema-without-id.json", schemasPath))
	require.Nil(t, loadReferencedSchema(compiler, "referenced-schema-2.json", schemasPath))
}

func Test_schemaID(t *testing.T) {
	_, err := schemaID("schema-without-id.json", schemasPath)
	require.NotNil(t, err)

	id, err := schemaID("valid-schema.json", schemasPath)
	require.Equal(t, "https://raw.githubusercontent.com/arduino/arduino-check/main/check/checkdata/schema/testdata/schema-with-references.json", id)
	require.Nil(t, err)
}

func Test_pathURI(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		require.Equal(t, "file:///c:/foo%20bar", pathURI(paths.New("c:/foo bar")))
	default:
		require.Equal(t, "file:///foo%20bar", pathURI(paths.New("/foo bar")))
	}
}

func Test_validationErrorSchemaPointerValue(t *testing.T) {
	validationError := jsonschema.ValidationError{
		SchemaURL: "https://raw.githubusercontent.com/arduino/arduino-check/main/check/checkdata/schema/testdata/referenced-schema-1.json",
		SchemaPtr: "#/definitions/patternObject/pattern",
	}

	schemaPointerValueInterface := validationErrorSchemaPointerValue(&validationError, schemasPath)
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
