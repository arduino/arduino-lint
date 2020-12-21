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

package checkfunctions

import (
	"github.com/arduino/arduino-lint/internal/check/checkresult"
	"github.com/arduino/arduino-lint/internal/project/checkdata"
)

// The check functions for package indexes.

// PackageIndexJSONFormat checks whether the package index file is a valid JSON document.
func PackageIndexJSONFormat() (result checkresult.Type, output string) {
	if isValidJSON(checkdata.ProjectPath()) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, ""
}

// PackageIndexFormat checks for invalid package index data format.
func PackageIndexFormat() (result checkresult.Type, output string) {
	if checkdata.PackageIndexLoadError() != nil {
		return checkresult.Fail, checkdata.PackageIndexLoadError().Error()
	}

	return checkresult.Pass, ""
}
