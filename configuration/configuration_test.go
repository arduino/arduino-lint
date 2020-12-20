// This file is part of arduino-lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-lint.
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

	"github.com/arduino/arduino-lint/configuration/checkmode"
	"github.com/arduino/arduino-lint/project/projecttype"
	"github.com/arduino/arduino-lint/result/outputformat"
	"github.com/arduino/arduino-lint/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var projectPaths []string

func init() {
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectPaths = []string{projectPath}
}

func TestInitializeCompliance(t *testing.T) {
	flags := test.ConfigurationFlags()

	flags.Set("compliance", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("compliance", "strict")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, customCheckModes[checkmode.Strict])
	assert.False(t, customCheckModes[checkmode.Specification])
	assert.False(t, customCheckModes[checkmode.Permissive])

	flags.Set("compliance", "specification")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, customCheckModes[checkmode.Strict])
	assert.True(t, customCheckModes[checkmode.Specification])
	assert.False(t, customCheckModes[checkmode.Permissive])

	flags.Set("compliance", "permissive")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, customCheckModes[checkmode.Strict])
	assert.False(t, customCheckModes[checkmode.Specification])
	assert.True(t, customCheckModes[checkmode.Permissive])
}

func TestInitializeFormat(t *testing.T) {
	flags := test.ConfigurationFlags()
	flags.Set("format", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("format", "text")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, outputformat.Text, OutputFormat())

	flags.Set("format", "json")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, outputformat.JSON, OutputFormat())
}

func TestInitializeLibraryManager(t *testing.T) {
	flags := test.ConfigurationFlags()
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
}

func TestInitializeLogFormat(t *testing.T) {
	os.Setenv("ARDUINO_LINT_LOG_FORMAT", "foo")
	assert.Error(t, Initialize(test.ConfigurationFlags(), projectPaths), "Invalid format")

	os.Setenv("ARDUINO_LINT_LOG_FORMAT", "text")
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths), "text format")

	os.Setenv("ARDUINO_LINT_LOG_FORMAT", "json")
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths), "json format")
}

func TestInitializeLogLevel(t *testing.T) {
	require.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths))

	os.Setenv("ARDUINO_LINT_LOG_LEVEL", "foo")
	assert.Error(t, Initialize(test.ConfigurationFlags(), projectPaths), "Invalid level")

	os.Setenv("ARDUINO_LINT_LOG_LEVEL", "info")
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths), "Valid level")
	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
}

func TestInitializeProjectType(t *testing.T) {
	flags := test.ConfigurationFlags()

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
}

func TestInitializeRecursive(t *testing.T) {
	flags := test.ConfigurationFlags()

	flags.Set("recursive", "foo")
	assert.Error(t, Initialize(flags, projectPaths))

	flags.Set("recursive", "true")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, Recursive())

	flags.Set("recursive", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, Recursive())
}

func TestInitializeReportFile(t *testing.T) {
	flags := test.ConfigurationFlags()

	flags.Set("report-file", "")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Nil(t, ReportFilePath())

	reportFilePath := paths.New("/bar")
	flags.Set("report-file", reportFilePath.String())
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.Equal(t, reportFilePath, ReportFilePath())
}

func TestInitializeVersion(t *testing.T) {
	flags := test.ConfigurationFlags()

	flags.Set("version", "true")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, VersionMode())

	flags.Set("version", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, VersionMode())
}

func TestInitializeVerbose(t *testing.T) {
	flags := test.ConfigurationFlags()

	flags.Set("verbose", "true")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.True(t, Verbose())

	flags.Set("verbose", "false")
	assert.Nil(t, Initialize(flags, projectPaths))
	assert.False(t, Verbose())
}

func TestInitializeProjectPath(t *testing.T) {
	assert.Nil(t, Initialize(test.ConfigurationFlags(), []string{}))
	workingDirectoryPath, err := os.Getwd()
	require.Nil(t, err)
	assert.Equal(t, paths.NewPathList(workingDirectoryPath), TargetPaths(), "Default PROJECT_PATH to current working directory")

	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths))
	assert.Equal(t, paths.NewPathList(projectPaths[0]), TargetPaths())

	assert.Error(t, Initialize(test.ConfigurationFlags(), []string{"/nonexistent"}))
}

func TestInitializeOfficial(t *testing.T) {
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths))
	assert.False(t, customCheckModes[checkmode.Official], "Default official check mode")

	os.Setenv("ARDUINO_LINT_OFFICIAL", "true")
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths))
	assert.True(t, customCheckModes[checkmode.Official])

	os.Setenv("ARDUINO_LINT_OFFICIAL", "false")
	assert.Nil(t, Initialize(test.ConfigurationFlags(), projectPaths))
	assert.False(t, customCheckModes[checkmode.Official])

	os.Setenv("ARDUINO_LINT_OFFICIAL", "invalid value")
	assert.Error(t, Initialize(test.ConfigurationFlags(), projectPaths))
}
