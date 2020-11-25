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

// Package cli defines the arduino-check command line interface.
package cli

import (
	"github.com/arduino/arduino-check/command"
	"github.com/spf13/cobra"
)

// Root creates a new arduino-check command root.
func Root() *cobra.Command {
	rootCommand := &cobra.Command{
		Short:                 "Linter for Arduino projects.",
		Long:                  "arduino-check checks for specification compliance and other common problems with Arduino projects",
		DisableFlagsInUseLine: true,
		Use:                   "arduino-check [FLAG]... [PROJECT_PATH]...\n\nRun checks on PROJECT_PATH or current path if no PROJECT_PATH argument provided.",
		Run:                   command.ArduinoCheck,
	}

	rootCommand.PersistentFlags().String("format", "text", "The output format can be {text|json}.")
	rootCommand.PersistentFlags().String("library-manager", "", "Configure the checks for libraries in the Arduino Library Manager index. Can be {submit|update|false}.\nsubmit: Also run additional checks required to pass before a library is accepted for inclusion in the index.\nupdate: Also run additional checks required to pass before new releases of a library already in the index are accepted.\nfalse: Don't run any Library Manager-specific checks.")
	rootCommand.PersistentFlags().String("log-format", "text", "The output format for the logs, can be {text|json}.")
	rootCommand.PersistentFlags().String("log-level", "panic", "Messages with this level and above will be logged. Valid levels are: trace, debug, info, warn, error, fatal, panic")
	rootCommand.PersistentFlags().Bool("permissive", false, "Only fail when critical issues are detected.")
	rootCommand.PersistentFlags().String("project-type", "all", "Only check projects of the specified type and their subprojects. Can be {sketch|library|all}.")
	rootCommand.PersistentFlags().Bool("recursive", true, "Search path recursively for Arduino projects to check. Can be {true|false}.")
	rootCommand.PersistentFlags().String("report-file", "", "Save a report on the checks to this file.")

	return rootCommand
}
