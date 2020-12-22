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

// Package test provides resources for testing arduino-lint.
package test

import "github.com/spf13/pflag"

// ConfigurationFlags returns a set of the flags used for command line configuration of arduino-lint.
func ConfigurationFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	flags.String("compliance", "specification", "")
	flags.String("format", "text", "")
	flags.String("library-manager", "", "")
	flags.String("log-format", "text", "")
	flags.String("log-level", "panic", "")
	flags.String("project-type", "all", "")
	flags.Bool("recursive", true, "")
	flags.String("report-file", "", "")
	flags.Bool("verbose", false, "")
	flags.Bool("version", false, "")

	return flags
}
