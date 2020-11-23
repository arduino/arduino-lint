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

// Package configuration handles the configuration of the arduino-check tool.
package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/arduino/go-paths-helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Initialize sets up the tool configuration according to defaults and user-specified options.
func Initialize(flags *pflag.FlagSet, projectPaths []string) error {
	var err error
	outputFormatString, _ := flags.GetString("format")
	outputFormat, err = outputformat.FromString(outputFormatString)
	if err != nil {
		return fmt.Errorf("--format flag value %s not valid", outputFormatString)
	}

	libraryManagerModeString, _ := flags.GetString("library-manager")
	if libraryManagerModeString != "" {
		customCheckModes[checkmode.LibraryManagerSubmission], customCheckModes[checkmode.LibraryManagerIndexed], err = checkmode.LibraryManagerModeFromString(libraryManagerModeString)
		if err != nil {
			return fmt.Errorf("--library-manager flag value %s not valid", libraryManagerModeString)
		}
	}

	logFormatString, _ := flags.GetString("log-format")
	logFormat, err := logFormatFromString(logFormatString)
	if err != nil {
		return fmt.Errorf("--log-format flag value %s not valid", logFormatString)
	}
	logrus.SetFormatter(logFormat)

	logLevelString, _ := flags.GetString("log-level")
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		return fmt.Errorf("--log-level flag value %s not valid", logLevelString)
	}
	logrus.SetLevel(logLevel)

	customCheckModes[checkmode.Permissive], _ = flags.GetBool("permissive")

	superprojectTypeFilterString, _ := flags.GetString("project-type")
	superprojectTypeFilter, err = projecttype.FromString(superprojectTypeFilterString)
	if err != nil {
		return fmt.Errorf("--project-type flag value %s not valid", superprojectTypeFilterString)
	}

	recursive, _ = flags.GetBool("recursive")

	reportFilePathString, _ := flags.GetString("report-file")
	reportFilePath = paths.New(reportFilePathString)

	// TODO validate target path value, exit if not found
	// TODO support multiple paths
	targetPath = paths.New(projectPaths[0])

	// TODO: set via environment variable
	// customCheckModes[checkmode.Official] = false

	logrus.WithFields(logrus.Fields{
		"output format":                   OutputFormat(),
		"Library Manager submission mode": customCheckModes[checkmode.LibraryManagerSubmission],
		"Library Manager update mode":     customCheckModes[checkmode.LibraryManagerIndexed],
		"log format":                      logFormatString,
		"log level":                       logrus.GetLevel().String(),
		"permissive":                      customCheckModes[checkmode.Permissive],
		"superproject type filter":        SuperprojectTypeFilter(),
		"recursive":                       Recursive(),
		"report file":                     ReportFilePath(),
		"projects path":                   TargetPath(),
	}).Debug("Configuration initialized")

	return nil
}

// logFormatFromString parses the --log-format flag value and returns the corresponding log formatter.
func logFormatFromString(logFormatString string) (logrus.Formatter, error) {
	switch strings.ToLower(logFormatString) {
	case "text":
		return &logrus.TextFormatter{}, nil
	case "json":
		return &logrus.JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("No matching log format for string %s", logFormatString)
	}
}

var customCheckModes = make(map[checkmode.Type]bool)

// CheckModes returns the check modes configuration for the given project type.
func CheckModes(superprojectType projecttype.Type) map[checkmode.Type]bool {
	return checkmode.Modes(defaultCheckModes, customCheckModes, superprojectType)
}

var superprojectTypeFilter projecttype.Type

// SuperprojectType returns the superproject type filter configuration.
func SuperprojectTypeFilter() projecttype.Type {
	return superprojectTypeFilter
}

var recursive bool

// Recursive returns the recursive project search configuration value.
func Recursive() bool {
	return recursive
}

var outputFormat outputformat.Type

// OutputFormat returns the tool output format configuration value.
func OutputFormat() outputformat.Type {
	return outputFormat
}

var reportFilePath *paths.Path

// ReportFilePath returns the path to save the report file at.
func ReportFilePath() *paths.Path {
	return reportFilePath
}

var targetPath *paths.Path

// TargetPath returns the projects search path.
func TargetPath() *paths.Path {
	return targetPath
}

// SchemasPath returns the path to the folder containing the JSON schemas.
func SchemasPath() *paths.Path {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return paths.New(executablePath).Parent().Join("etc", "schemas")
}
