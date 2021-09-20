// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License, either
// version 3 of the License, or (at your option) any later version.
// This license covers the main part of Arduino Lint.
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
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = paths.New(workingDirectory, "testdata")
}

func TestSymlinkLoop(t *testing.T) {
	// Set up directory structure of test library.
	libraryPath, err := paths.TempDir().MkTempDir("TestSymlinkLoop")
	defer libraryPath.RemoveAll() // Clean up after the test.
	require.Nil(t, err)
	err = libraryPath.Join("TestSymlinkLoop.h").WriteFile([]byte{})
	require.Nil(t, err)
	examplesPath := libraryPath.Join("examples")
	err = examplesPath.Mkdir()
	require.Nil(t, err)

	// It's probably most friendly for contributors using Windows to create the symlinks needed for the test on demand.
	err = os.Symlink(examplesPath.Join("..").String(), examplesPath.Join("UpGoer1").String())
	require.Nil(t, err, "This test must be run as administrator on Windows to have symlink creation privilege.")
	// It's necessary to have multiple symlinks to a parent directory to create the loop.
	err = os.Symlink(examplesPath.Join("..").String(), examplesPath.Join("UpGoer2").String())
	require.Nil(t, err)

	configuration.Initialize(test.ConfigurationFlags(), []string{libraryPath.String()})

	assert.Panics(t, func() { FindProjects() }, "Infinite symlink loop encountered during project discovery")
}

func TestFindProjects(t *testing.T) {
	sketchPath := testDataPath.Join("Sketch")
	libraryPath := testDataPath.Join("Library")
	libraryExamplePath := testDataPath.Join("Library", "examples", "Example")
	platformPath := testDataPath.Join("Platform")
	platformBundledLibraryPath := testDataPath.Join("Platform", "libraries", "Library")
	platformBundledLibraryExamplePath := testDataPath.Join("Platform", "libraries", "Library", "examples", "Example")
	packageIndexFolderPath := testDataPath.Join("PackageIndex")
	packageIndexFilePath := packageIndexFolderPath.Join("package_foo_index.json")
	projectsPath := testDataPath.Join("Projects")
	projectsPathSketch := testDataPath.Join("Projects", "Sketch")
	projectsPathLibrary := testDataPath.Join("Projects", "Library")
	projectsPathLibraryExample := testDataPath.Join("Projects", "Library", "examples", "Example")

	testTables := []struct {
		testName               string
		superprojectTypeFilter []string
		recursive              []string
		projectPaths           []string
		errorAssertion         assert.ErrorAssertionFunc
		expectedProjects       []Type
	}{
		{
			"Sketch file",
			[]string{"all", "sketch"},
			[]string{"true", "false"},
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
			[]string{"all", "library"},
			[]string{"true", "false"},
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
			[]string{"all", "platform"},
			[]string{"true", "false"},
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
			[]string{"all", "package-index"},
			[]string{"true", "false"},
			[]string{packageIndexFilePath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             packageIndexFilePath,
					ProjectType:      projecttype.PackageIndex,
					SuperprojectType: projecttype.PackageIndex,
				},
			},
		},
		{
			"Explicit file",
			[]string{"sketch"},
			[]string{"true", "false"},
			[]string{libraryPath.Join("Library.h").String()},
			assert.NoError,
			[]Type{
				{
					Path:             libraryPath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Sketch folder",
			[]string{"all", "sketch"},
			[]string{"true", "false"},
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
			[]string{"all", "library"},
			[]string{"true", "false"},
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
			[]string{"all", "platform"},
			[]string{"true", "false"},
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
			[]string{"all", "package-index"},
			[]string{"true", "false"},
			[]string{packageIndexFolderPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             packageIndexFolderPath,
					ProjectType:      projecttype.PackageIndex,
					SuperprojectType: projecttype.PackageIndex,
				},
			},
		},
		{
			"Explicit folder",
			[]string{"sketch"},
			[]string{"false"},
			[]string{libraryPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             libraryPath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Explicit folder",
			[]string{"sketch"},
			[]string{"true"},
			[]string{libraryPath.String()},
			assert.NoError,
			[]Type{
				{
					Path:             libraryExamplePath,
					ProjectType:      projecttype.Sketch,
					SuperprojectType: projecttype.Sketch,
				},
			},
		},
		{
			"Projects folder",
			[]string{"all"},
			[]string{"true"},
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
			"Projects folder",
			[]string{"sketch"},
			[]string{"true"},
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
		{
			"Projects folder, non-recursive",
			[]string{"all"},
			[]string{"false"},
			[]string{projectsPath.String()},
			assert.Error,
			[]Type{},
		},
		{
			"Multiple target folders",
			[]string{"all"},
			[]string{"true"},
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
	}

	for _, testTable := range testTables {
		for _, superprojectTypeFilter := range testTable.superprojectTypeFilter {
			for _, recursive := range testTable.recursive {
				flags := test.ConfigurationFlags()
				flags.Set("project-type", superprojectTypeFilter)
				if recursive != "" {
					flags.Set("recursive", recursive)
				}
				configuration.Initialize(flags, testTable.projectPaths)
				foundProjects, err := FindProjects()
				testTable.errorAssertion(t, err)
				if err == nil {
					assert.True(
						t,
						reflect.DeepEqual(foundProjects, testTable.expectedProjects),
						fmt.Sprintf(
							"%s (%s project-type=%s recursive=%s)",
							testTable.testName,
							testTable.projectPaths,
							superprojectTypeFilter, recursive,
						),
					)
				}
			}
		}
	}
}
