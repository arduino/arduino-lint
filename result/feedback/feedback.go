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

// Package feedback provides feedback to the user.
package feedback

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/sirupsen/logrus"
)

// VerbosePrintln behaves like Println but only prints when verbosity is enabled.
func VerbosePrintln(v ...interface{}) {
	VerbosePrint(v...)
	VerbosePrint("\n")
}

// VerbosePrintf behaves like Printf but only prints when verbosity is enabled.
func VerbosePrintf(format string, v ...interface{}) {
	VerbosePrint(fmt.Sprintf(format, v...))
}

// VerbosePrint behaves like Print but only prints when verbosity is enabled.
func VerbosePrint(v ...interface{}) {
	if configuration.Verbose() && (configuration.OutputFormat() == outputformat.Text) {
		Print(v...)
	}
}

// Println behaves like fmt.Println but only prints when output format is set to `text`.
func Println(v ...interface{}) {
	Print(v...)
	Print("\n")
}

// Printf behaves like fmt.Printf but only prints when output format is set to `text`.
func Printf(format string, v ...interface{}) {
	Print(fmt.Sprintf(format, v...))
}

// Print behaves like fmt.Print but only prints when output format is set to `text`.
func Print(v ...interface{}) {
	if configuration.OutputFormat() == outputformat.Text {
		fmt.Print(v...)
	}
}

// Errorf behaves like fmt.Printf but adds a newline and also logs the error.
func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

// Error behaves like fmt.Print but adds a newline and also logs the error.
func Error(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	logrus.Error(fmt.Sprint(v...))
}
