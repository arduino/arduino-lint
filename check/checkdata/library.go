package checkdata

import (
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

// Initialize gathers the library check data for the specified project.
func InitializeForLibrary(project project.Type) {
	libraryProperties, libraryPropertiesLoadError = libraryproperties.Properties(project.Path)
	if libraryPropertiesLoadError != nil {
		// TODO: can I even do this?
		libraryPropertiesSchemaValidationResult = nil
	} else {
		libraryPropertiesSchemaValidationResult = libraryproperties.Validate(libraryProperties)
	}
}

var libraryPropertiesLoadError error

// LibraryPropertiesLoadError returns the error output from loading the library.properties metadata file.
func LibraryPropertiesLoadError() error {
	return libraryPropertiesLoadError
}

var libraryProperties *properties.Map

// LibraryProperties returns the data from the library.properties metadata file.
func LibraryProperties() *properties.Map {
	return libraryProperties
}

var libraryPropertiesSchemaValidationResult *gojsonschema.Result

// LibraryPropertiesSchemaValidationResult returns the result of validating library.properties against the JSON schema.
// See: https://github.com/xeipuuv/gojsonschema
func LibraryPropertiesSchemaValidationResult() *gojsonschema.Result {
	return libraryPropertiesSchemaValidationResult
}
