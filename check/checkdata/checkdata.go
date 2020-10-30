/*
Package checkdata handles the collection of data specific to a project before running the checks on it.
This is for data required by multiple checks.
*/
package checkdata

import (
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
)

// Initialize gathers the check data for the specified project.
func Initialize(project project.Type) {
	projectType = project.ProjectType
	projectPath = project.Path
	switch project.ProjectType {
	case projecttype.Sketch:
	case projecttype.Library:
		InitializeForLibrary(project)
	case projecttype.Platform:
	case projecttype.PackageIndex:
	}
}

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
