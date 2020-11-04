// Package schema contains code for working with JSON schema.
package schema

import (
	"os"
	"regexp"

	"github.com/arduino/arduino-check/util"
	"github.com/arduino/go-paths-helper"
	"github.com/xeipuuv/gojsonschema"
)

// SchemasPath returns the path to the folder containing the JSON schemas.
func SchemasPath() *paths.Path {
	workingPath, _ := os.Getwd()
	return paths.New(workingPath)
}

// Compile compiles the schema files specified by the filename arguments and returns the compiled schema.
func Compile(schemaFilename string, referencedSchemaFilenames []string, schemasPath *paths.Path) *gojsonschema.Schema {
	schemaLoader := gojsonschema.NewSchemaLoader()

	// Load the referenced schemas.
	for _, referencedSchemaFilename := range referencedSchemaFilenames {
		referencedSchemaPath := schemasPath.Join(referencedSchemaFilename)
		referencedSchemaURI := util.PathURI(referencedSchemaPath)
		err := schemaLoader.AddSchemas(gojsonschema.NewReferenceLoader(referencedSchemaURI))
		if err != nil {
			panic(err.Error())
		}
	}

	// Compile the schema.
	schemaPath := schemasPath.Join(schemaFilename)
	schemaURI := util.PathURI(schemaPath)
	compiledSchema, err := schemaLoader.Compile(gojsonschema.NewReferenceLoader(schemaURI))
	if err != nil {
		panic(err.Error())
	}
	return compiledSchema
}

// Validate validates an instance against a JSON schema and returns the gojsonschema.Result object.
func Validate(instanceObject interface{}, schemaObject *gojsonschema.Schema) *gojsonschema.Result {
	result, err := schemaObject.Validate(gojsonschema.NewGoLoader(instanceObject))
	if err != nil {
		panic(err.Error())
	}

	return result
}

// RequiredPropertyMissing returns whether the given required property is missing from the document.
func RequiredPropertyMissing(propertyName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("required", "(root)", propertyName+" is required", validationResult)
}

// PropertyPatternMismatch returns whether the given property did not match the regular expression defined in the JSON schema.
func PropertyPatternMismatch(propertyName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("pattern", propertyName, "", validationResult)
}

// ValidationErrorMatch returns whether the given query matches against the JSON schema validation error.
// See: https://github.com/xeipuuv/gojsonschema#working-with-errors
func ValidationErrorMatch(typeQuery string, fieldQuery string, descriptionQueryRegexp string, validationResult *gojsonschema.Result) bool {
	if validationResult.Valid() {
		// No error, so nothing to match
		return false
	}
	for _, validationError := range validationResult.Errors() {
		if typeQuery == "" || typeQuery == validationError.Type() {
			if fieldQuery == "" || fieldQuery == validationError.Field() {
				descriptionQuery := regexp.MustCompile(descriptionQueryRegexp)
				return descriptionQuery.MatchString(validationError.Description())
			}
		}
	}

	return false
}
