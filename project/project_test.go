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

package project

import (
	"os"
	"reflect"
	"testing"

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = paths.New(workingDirectory, "testdata")
}

func TestFindProjects(t *testing.T) {
	sketchPath := testDataPath.Join("Sketch")
	libraryPath := testDataPath.Join("Library")
	libraryExamplePath := testDataPath.Join("Library", "examples", "Example")
	platformPath := testDataPath.Join("Platform")
	platformBundledLibraryPath := testDataPath.Join("Platform", "libraries", "Library")
	platformBundledLibraryExamplePath := testDataPath.Join("Platform", "libraries", "Library", "examples", "Example")
	packageIndexPath := testDataPath.Join("PackageIndex")
	projectsPath := testDataPath.Join("Projects")
	projectsPathSketch := testDataPath.Join("Projects", "Sketch")
	projectsPathLibrary := testDataPath.Join("Projects", "Library")
	projectsPathLibraryExample := testDataPath.Join("Projects", "Library", "examples", "Example")

	testTables := []struct {
		testName               string
		superprojectTypeFilter string
		recursive              string
		projectPaths           []string
		errorAssertion         assert.ErrorAssertionFunc
		expectedProjects       []Type
	}{
		{
			"Sketch file",
			"all",
			"",
			[]string{sketchPath.Join("Sketch.ino").String()},
			assert.NoError,
			[]Type{
				{
					Path:             sketchPath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Library file",
			"all",
			"",
			[]string{libraryPath.Join("Library.h").String()},
			assert.NoError,
			[]Type{
				{
					Path:             libraryPath,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             libraryExamplePath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Library,
				},
			},
		},
		{
			"Platform file",
			"all",
			"",
			[]string{platformPath.Join("boards.txt").String()},
			assert.NoError,
			[]Type{
				{
					Path:             platformPath,
					ProjectType:      projecttype.Platform,
					SuperprojectType: projecttype.Platform,
				},
				{
					Path:             platformBundledLibraryPath,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Platform,
				},
				{
					Path:             platformBundledLibraryExamplePath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Platform,
				},
			},
		},
		{
			"Package index file",
			"all",
			"",
			[]string{packageIndexPath.Join("package_foo_index.json").String()},
			assert.NoError,
			[]Type{
				{
					Path:             packageIndexPath,
					ProjectType:      projecttype.PackageIndex,
					SuperprojectType: projecttype.PackageIndex,
				},
			},
		},
		{
			"Sketch folder",
			"all",
			"",
			[]string{sketchPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             sketchPath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Library folder",
			"all",
			"",
			[]string{libraryPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             libraryPath,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             libraryExamplePath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Library,
				},
			},
		},
		{
			"Platform folder",
			"all",
			"",
			[]string{platformPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             platformPath,
					ProjectType:      projecttype.Platform,
					SuperprojectType: projecttype.Platform,
				},
				{
					Path:             platformBundledLibraryPath,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Platform,
				},
				{
					Path:             platformBundledLibraryExamplePath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Platform,
				},
			},
		},
		{
			"Package index folder",
			"all",
			"",
			[]string{packageIndexPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             packageIndexPath,
					ProjectType:      projecttype.PackageIndex,
					SuperprojectType: projecttype.PackageIndex,
				},
			},
		},
		{
			"Projects folder, recursive",
			"all",
			"true",
			[]string{projectsPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             projectsPathLibrary,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             projectsPathLibraryExample,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             projectsPathSketch,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Projects folder, non-recursive",
			"all",
			"false",
			[]string{projectsPath.String()},
			assert.Error,
			[]Type{},
		},
		{
			"Multiple target folders",
			"all",
			"true",
			[]string{projectsPath.String(), sketchPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             projectsPathLibrary,
					ProjectType:      projecttype.Library,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             projectsPathLibraryExample,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Library,
				},
				{
					Path:             projectsPathSketch,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
				{
					Path:             sketchPath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"superproject type filter",
			"sketch",
			"true",
			[]string{projectsPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             projectsPathLibraryExample,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
				{
					Path:             projectsPathSketch,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
	}

	for _, testTable := range testTables {
		flags := test.ConfigurationFlags()
		flags.Set("project-type", testTable.superprojectTypeFilter)
		if testTable.recursive != "" {
			flags.Set("recursive", testTable.recursive)
		}
		configuration.Initialize(flags, testTable.projectPaths)
		foundProjects, err := FindProjects()
		testTable.errorAssertion(t, err)
		if err == nil {
			assert.True(t, reflect.DeepEqual(foundProjects, testTable.expectedProjects), testTable.testName)
		}
	}
}
