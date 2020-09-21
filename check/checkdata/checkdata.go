/*
Package checkdata handles the collection of data specific to a project before running the checks on it.
This is for data required by multiple checks.
*/
package checkdata

import (
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

var projectType projecttype.Type

// ProjectType returns the type of the project being checked.
func ProjectType() projecttype.Type {
	return projectType
}

var projectPath *paths.Path

// ProjectPath returns the path to the project being checked.
func ProjectPath() *paths.Path {
	return projectPath
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

// Initialize gathers the check data for the specified project.
func Initialize(project project.Type) {
	projectType = project.ProjectType
	projectPath = project.Path
	switch project.ProjectType {
	case projecttype.Sketch:
	case projecttype.Library:
		libraryProperties, libraryPropertiesLoadError = libraryproperties.Properties(project.Path)
		if libraryPropertiesLoadError != nil {
			// TODO: can I even do this?
			libraryPropertiesSchemaValidationResult = nil
		} else {
			libraryPropertiesSchemaValidationResult = libraryproperties.Validate(libraryProperties)
		}
	case projecttype.Platform:
	case projecttype.PackageIndex:
	}
}
