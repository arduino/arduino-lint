// Package libraryproperties provides functions for working with the library.properties Arduino library metadata file.
package libraryproperties

import (
	"github.com/arduino/arduino-check/check/checkdata/schema"
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
	referencedSchemaFilenames := []string{}
	schemaObject := schema.Compile("arduino-library-properties-schema.json", referencedSchemaFilenames)

	return schema.Validate(libraryProperties, schemaObject)
}
