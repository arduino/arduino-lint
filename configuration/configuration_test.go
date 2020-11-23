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

package configuration

import (
	"os"
	"testing"

	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/arduino/arduino-check/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	flags := test.ConfigurationFlags()

	projectPath, err := os.Getwd()
	require.Nil(t, err)
	projectPaths := []string{projectPath}

	flags.Set("format", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("format", "text")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, outputformat.Text, OutputFormat())

	flags.Set("format", "json")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, outputformat.JSON, OutputFormat())

	flags.Set("library-manager", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	customCheckModes = make(map[checkmode.Type]bool)
	flags.Set("library-manager", "")
	assert.Nil(t, Initialize(flags, projectPaths))
	_, ok := customCheckModes[checkmode.LibraryManagerSubmission]
	assert.False(t, ok)
	_, ok = customCheckModes[checkmode.LibraryManagerIndexed]
	assert.False(t, ok)

	flags.Set("library-manager", "submit")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, customCheckModes[checkmode.LibraryManagerSubmission])
	assert.False(t, customCheckModes[checkmode.LibraryManagerIndexed])

	flags.Set("library-manager", "update")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, customCheckModes[checkmode.LibraryManagerSubmission])
	assert.True(t, customCheckModes[checkmode.LibraryManagerIndexed])

	flags.Set("library-manager", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, customCheckModes[checkmode.LibraryManagerSubmission])
	assert.False(t, customCheckModes[checkmode.LibraryManagerIndexed])

	flags.Set("log-format", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("log-format", "text")
	assert.Nil(t, Initialize(flags, projectPaths))

	flags.Set("log-format", "json")
	assert.Nil(t, Initialize(flags, projectPaths))

	flags.Set("permissive", "true")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, customCheckModes[checkmode.Permissive])

	flags.Set("permissive", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, customCheckModes[checkmode.Permissive])

	flags.Set("project-type", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("project-type", "sketch")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, projecttype.Sketch, SuperprojectTypeFilter())

	flags.Set("project-type", "library")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, projecttype.Library, SuperprojectTypeFilter())

	flags.Set("project-type", "platform")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, projecttype.Platform, SuperprojectTypeFilter())

	flags.Set("project-type", "package-index")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, projecttype.PackageIndex, SuperprojectTypeFilter())

	flags.Set("project-type", "all")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, projecttype.All, SuperprojectTypeFilter())

	flags.Set("recursive", "true")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, Recursive())

	flags.Set("recursive", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, Recursive())

	flags.Set("report-file", "")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Nil(t, ReportFilePath())

	reportFilePath := paths.New("/bar")
	flags.Set("report-file", reportFilePath.String())
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, reportFilePath, ReportFilePath())

	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, paths.New(projectPaths[0]), TargetPath())

	assert.Error(t, Initialize(flags, []string{"/nonexistent"}))
}
