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

package checkdata

import (
	"testing"

	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var platformTestDataPath *paths.Path

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	platformTestDataPath = workingDirectory.Join("testdata", "platforms")
}

func TestInitializeForPlatform(t *testing.T) {
	testTables := []struct {
		testName                    string
		platformFolderName          string
		boardsTxtAssertion          assert.ValueAssertionFunc
		boardsTxtLoadErrorAssertion assert.ValueAssertionFunc
	}{
		{"Valid", "valid-boards.txt", assert.NotNil, assert.Nil},
		{"Invalid", "invalid-boards.txt", assert.Nil, assert.NotNil},
		{"Missing", "missing-boards.txt", assert.NotNil, assert.Nil},
	}

	for _, testTable := range testTables {

		testProject := project.Type{
			Path:             platformTestDataPath.Join(testTable.platformFolderName),
			ProjectType:      projecttype.Platform,
			SuperprojectType: projecttype.Platform,
		}
		Initialize(testProject)

		testTable.boardsTxtLoadErrorAssertion(t, BoardsTxtLoadError(), testTable.testName)
		if BoardsTxtLoadError() == nil {
			testTable.boardsTxtAssertion(t, BoardsTxt(), testTable.testName)
		}
	}
}
