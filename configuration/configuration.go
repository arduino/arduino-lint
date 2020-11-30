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
	"strconv"
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

	complianceString, _ := flags.GetString("compliance")
	if complianceString != "" {
		customCheckModes[checkmode.Strict], customCheckModes[checkmode.Specification], customCheckModes[checkmode.Permissive], err = checkmode.ComplianceModeFromString(complianceString)
		if err != nil {
			return fmt.Errorf("--compliance flag value %s not valid", complianceString)
		}
	}

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

	superprojectTypeFilterString, _ := flags.GetString("project-type")
	superprojectTypeFilter, err = projecttype.FromString(superprojectTypeFilterString)
	if err != nil {
		return fmt.Errorf("--project-type flag value %s not valid", superprojectTypeFilterString)
	}

	recursive, _ = flags.GetBool("recursive")

	reportFilePathString, _ := flags.GetString("report-file")
	reportFilePath = paths.New(reportFilePathString)

	versionMode, _ = flags.GetBool("version")

	targetPaths = nil
	if len(projectPaths) == 0 {
		// Default to using current working directory.
		workingDirectoryPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Error when setting default PROJECT_PATH argument: %s", err)
		}
		targetPaths.Add(paths.New(workingDirectoryPath))
	} else {
		for _, projectPath := range projectPaths {
			targetPath := paths.New(projectPath)
			targetPathExists, err := targetPath.ExistCheck()
			if err != nil {
				return fmt.Errorf("Problem processing PROJECT_PATH argument value %v: %v", targetPath, err)
			}
			if !targetPathExists {
				return fmt.Errorf("PROJECT_PATH argument %v does not exist", targetPath)
			}
			targetPaths.AddIfMissing(targetPath)
		}
	}

	if officialModeString, ok := os.LookupEnv("ARDUINO_CHECK_OFFICIAL"); ok {
		customCheckModes[checkmode.Official], err = strconv.ParseBool(officialModeString)
		if err != nil {
			return fmt.Errorf("ARDUINO_CHECK_OFFICIAL environment variable value %s not valid", officialModeString)
		}
	}

	logrus.WithFields(logrus.Fields{
		"compliance strict mode":          customCheckModes[checkmode.Strict],
		"compliance specification mode":   customCheckModes[checkmode.Specification],
		"compliance permissive mode":      customCheckModes[checkmode.Permissive],
		"output format":                   OutputFormat(),
		"Library Manager submission mode": customCheckModes[checkmode.LibraryManagerSubmission],
		"Library Manager update mode":     customCheckModes[checkmode.LibraryManagerIndexed],
		"log format":                      logFormatString,
		"log level":                       logrus.GetLevel().String(),
		"superproject type filter":        SuperprojectTypeFilter(),
		"recursive":                       Recursive(),
		"report file":                     ReportFilePath(),
		"projects path":                   TargetPaths(),
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

// SuperprojectTypeFilter returns the superproject type filter configuration.
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

var versionMode bool

func VersionMode() bool {
	return versionMode
}

var version string
var commit string

func Version() string {
	if version == "" {
		return "0.0.0+" + commit
	}

	return version
}

var buildTimestamp string

func BuildTimestamp() string {
	return buildTimestamp
}

var targetPaths paths.PathList

// TargetPaths returns the projects search paths.
func TargetPaths() paths.PathList {
	return targetPaths
}

// SchemasPath returns the path to the folder containing the JSON schemas.
func SchemasPath() *paths.Path {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return paths.New(executablePath).Parent().Join("etc", "schemas")
}
