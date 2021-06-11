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

package projectdata

import (
	clipackageindex "github.com/arduino/arduino-cli/arduino/cores/packageindex"
	"github.com/arduino/arduino-lint/internal/project/packageindex"
)

// InitializeForPackageIndex gathers the package index rule data for the specified project.
func InitializeForPackageIndex() {
	packageIndex, packageIndexLoadError = packageindex.Properties(ProjectPath())
	if ProjectPath() != nil {
		_, packageIndexCLILoadError = clipackageindex.LoadIndex(ProjectPath())
	}
}

var packageIndex map[string]interface{}

// PackageIndex returns the package index data.
func PackageIndex() map[string]interface{} {
	return packageIndex
}

var packageIndexLoadError error

// PackageIndexLoadError returns the error from loading the package index.
func PackageIndexLoadError() error {
	return packageIndexLoadError
}

var packageIndexCLILoadError error

// PackageIndexCLILoadError returns the error return of Arduino CLI's packageindex.LoadIndex().
func PackageIndexCLILoadError() error {
	return packageIndexCLILoadError
}
