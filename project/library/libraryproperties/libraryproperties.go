package libraryproperties

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

func Properties(libraryPath *paths.Path) (*properties.Map, error) {
	libraryProperties, err := properties.Load(libraryPath.Join("library.properties").String())
	if err != nil {
		return nil, err
	}
	return libraryProperties, nil
}

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

func FieldMissing(fieldName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("required", "(root)", fieldName+" is required", validationResult)
}

func FieldPatternMismatch(fieldName string, validationResult *gojsonschema.Result) bool {
	return ValidationErrorMatch("pattern", fieldName, "", validationResult)
}

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
