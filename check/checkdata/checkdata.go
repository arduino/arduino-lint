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

func ProjectType() projecttype.Type {
	return projectType
}

var projectPath *paths.Path

func ProjectPath() *paths.Path {
	return projectPath
}

var libraryPropertiesLoadError error

func LibraryPropertiesLoadError() error {
	return libraryPropertiesLoadError
}

var libraryProperties *properties.Map

func LibraryProperties() *properties.Map {
	return libraryProperties
}

var libraryPropertiesSchemaValidationResult *gojsonschema.Result

func LibraryPropertiesSchemaValidationResult() *gojsonschema.Result {
	return libraryPropertiesSchemaValidationResult
}

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
