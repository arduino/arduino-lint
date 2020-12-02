// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

/*
Package checkdata handles the collection of data specific to a project before running the checks on it.
This is for data required by multiple checks.
*/
package checkdata

import (
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/packageindex"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
)

// Initialize gathers the check data for the specified project.
func Initialize(project project.Type, schemasPath *paths.Path) {
	projectType = project.ProjectType
	projectPath = project.Path
	switch project.ProjectType {
	case projecttype.Sketch:
	case projecttype.Library:
		InitializeForLibrary(project, schemasPath)
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
