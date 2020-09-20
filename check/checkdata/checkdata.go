package checkdata

import (
	"github.com/arduino/arduino-check/projects"
	"github.com/arduino/arduino-check/projects/library"
	"github.com/arduino/arduino-check/projects/projecttype"
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

func Initialize(project projects.Type) {
	projectType = project.ProjectType
	projectPath = project.Path
	switch project.ProjectType {
	case projecttype.Sketch:
	case projecttype.Library:
		libraryProperties, libraryPropertiesLoadError = library.Properties(project.Path)
		if libraryPropertiesLoadError != nil {
			libraryPropertiesSchemaValidationResult = library.ValidateProperties(libraryProperties)
		} else {
			// TODO: can I even do this?
			libraryPropertiesSchemaValidationResult = nil
		}
	case projecttype.Platform:
	case projecttype.PackageIndex:
	}
}
