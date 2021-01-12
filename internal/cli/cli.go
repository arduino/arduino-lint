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

// Package cli defines the Arduino Lint command line interface.
package cli

import (
	"github.com/arduino/arduino-lint/internal/command"
	"github.com/spf13/cobra"
)

// Root creates a new arduino-lint command root.
func Root() *cobra.Command {
	rootCommand := &cobra.Command{
		Short:                 "Linter for Arduino projects.",
		Long:                  "Arduino Lint checks for specification compliance and other common problems with Arduino projects",
		DisableFlagsInUseLine: true,
		Use:                   "arduino-lint [FLAG]... [PROJECT_PATH]...\n\nLint project in PROJECT_PATH or current path if no PROJECT_PATH argument provided.",
		Run:                   command.ArduinoLint,
	}

	rootCommand.PersistentFlags().String("compliance", "specification", "Configure how strict the tool is. Can be {strict|specification|permissive}")
	rootCommand.PersistentFlags().String("format", "text", "The output format can be {text|json}.")
	rootCommand.PersistentFlags().String("library-manager", "", "Configure the rules for libraries in the Arduino Library Manager index. Can be {submit|update|false}.\nsubmit: Also run additional rules required to pass before a library is accepted for inclusion in the index.\nupdate: Also run additional rules required to pass before new releases of a library already in the index are accepted.\nfalse: Don't run any Library Manager-specific rules.")
	rootCommand.PersistentFlags().String("project-type", "all", "Only lint projects of the specified type and their subprojects. Can be {sketch|library|platform|all}.")
	rootCommand.PersistentFlags().Bool("recursive", false, "Search path recursively for Arduino projects to lint. Can be {true|false}.")
	rootCommand.PersistentFlags().String("report-file", "", "Save a report on the rules to this file.")
	rootCommand.PersistentFlags().BoolP("verbose", "v", false, "Show more information while running rules.")
	rootCommand.PersistentFlags().Bool("version", false, "Print version and timestamp of the build.")

	return rootCommand
}
