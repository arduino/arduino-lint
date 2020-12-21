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
	"github.com/arduino/arduino-lint/internal/project/projectdata"
)

// The check functions for platforms.

// BoardsTxtMissing checks whether the platform contains a boards.txt
func BoardsTxtMissing() (result checkresult.Type, output string) {
	boardsTxtPath := projectdata.ProjectPath().Join("boards.txt")
	exist, err := boardsTxtPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if exist {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, boardsTxtPath.String()
}

// BoardsTxtFormat checks for invalid boards.txt format.
func BoardsTxtFormat() (result checkresult.Type, output string) {
	if !projectdata.ProjectPath().Join("boards.txt").Exist() {
		return checkresult.NotRun, "boards.txt missing"
	}

	if projectdata.BoardsTxtLoadError() == nil {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, projectdata.BoardsTxtLoadError().Error()
}
