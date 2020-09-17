package library

import (
	"os"

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

func ValidateProperties(libraryProperties *properties.Map) *gojsonschema.Result {
	workingPath, _ := os.Getwd()
	schemaPath := paths.New(workingPath).Join("arduino-library-properties-schema.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath.String())
	documentLoader := gojsonschema.NewGoLoader(libraryProperties)

	result, _ := gojsonschema.Validate(schemaLoader, documentLoader)
	return result
}
