// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

/*
Package projectdata handles the collection of data specific to a project before running the rules on it.
This is for data required by multiple rules.
*/
package projectdata

import (
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/packageindex"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/go-paths-helper"
)

// Initialize gathers the check data for the specified project.
func Initialize(project project.Type) {
	superprojectType = project.SuperprojectType
	projectType = project.ProjectType
	projectPath = project.Path
	switch project.ProjectType {
	case projecttype.Sketch:
		InitializeForSketch(project)
	case projecttype.Library:
		InitializeForLibrary(project)
	case projecttype.Platform:
		InitializeForPlatform(project)
	case projecttype.PackageIndex:
		var err error
		// Because a package index project is a file, but project.Path may be a folder, an extra discovery step is needed for this project type.
		projectPath, err = packageindex.Find(project.Path)
		if err != nil {
			panic(err)
		}

		InitializeForPackageIndex()
	}
}

var superprojectType projecttype.Type

// SuperProjectType returns the type of the project being checked.
func SuperProjectType() projecttype.Type {
	return superprojectType
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
