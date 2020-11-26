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

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/sirupsen/logrus"
)

// Printf behaves like fmt.Printf but only prints when output format is set to `text`.
func Printf(format string, v ...interface{}) {
	Print(fmt.Sprintf(format, v...))
}

// Print behaves like fmt.Print but only prints when output format is set to `text`.
func Print(message string) {
	if configuration.OutputFormat() == outputformat.Text {
		fmt.Printf(message)
	}
}

// Errorf behaves like fmt.Printf but also logs the error.
func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

// Error behaves like fmt.Print but also logs the error.
func Error(errorMessage string) {
	fmt.Printf(errorMessage)
	logrus.Error(fmt.Sprint(errorMessage))
}
