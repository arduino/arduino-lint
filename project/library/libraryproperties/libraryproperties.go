// Package libraryproperties provides functions for working with the library.properties Arduino library metadata file.
package libraryproperties

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

// Properties parses the library.properties from the given path and returns the data.
func Properties(libraryPath *paths.Path) (*properties.Map, error) {
	libraryProperties, err := properties.Load(libraryPath.Join("library.properties").String())
	if err != nil {
		return nil, err
	}
	return libraryProperties, nil
}

// Validate validates library.properties data against the JSON schema.
func Validate(libraryProperties *properties.Map) *gojsonschema.Result {
	workingPath, _ := os.Getwd()
	schemaPath := paths.New(workingPath).Join("arduino-library-properties-schema.json")
	uriFriendlySchemaPath := filepath.ToSlash(schemaPath.String())
	schemaPathURI := url.URL{
		Scheme: "file",
		Path:   uriFriendlySchemaPath,
	}
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPathURI.String())

	documentLoader := gojsonschema.NewGoLoader(libraryProperties)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	return result
}

// FieldMissing returns whether the given required field is missing from library.properties.
func FieldMissing(fieldName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("required", "(root)", fieldName+" is required", validationResult)
}

// FieldPatternMismatch returns whether the given field did not match the regular expression defined in the JSON schema.
func FieldPatternMismatch(fieldName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("pattern", fieldName, "", validationResult)
}

// ValidationErrorMatch returns whether the given query matches against the JSON schema validation error.
// See: https://github.com/xeipuuv/gojsonschema#working-with-errors
func ValidationErrorMatch(typeQuery string, fieldQuery string, descriptionQuery string, validationResult *gojsonschema.Result) bool {
	if validationResult.Valid() {
		// No error, so nothing to match
		return false
	}
	for _, validationError := range validationResult.Errors() {
		if typeQuery == "" || typeQuery == validationError.Type() {
			if fieldQuery == "" || fieldQuery == validationError.Field() {
				if descriptionQuery == "" || descriptionQuery == validationError.Description() {
					return true
				}
			}
		}
	}

	return false
}
