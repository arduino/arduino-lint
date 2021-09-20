// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License, either
// version 3 of the License, or (at your option) any later version.
// This license covers the main part of Arduino Lint.
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
	"reflect"
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
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
		testName                                    string
		platformFolderName                          string
		boardsTxtAssertion                          assert.ValueAssertionFunc
		boardsTxtLoadErrorAssertion                 assert.ValueAssertionFunc
		platformTxtExistsAssertion                  assert.BoolAssertionFunc
		platformTxtAssertion                        assert.ValueAssertionFunc
		platformTxtLoadErrorAssertion               assert.ValueAssertionFunc
		platformTxtSchemaValidationResultAssertion  assert.ValueAssertionFunc
		platformTxtPluggableDiscoveryNamesAssertion []string
		platformTxtUserProvidedFieldNamesAssertion  map[string][]string
		platformTxtToolNamesAssertion               []string
	}{
		{"Valid boards.txt", "valid-boards.txt", assert.NotNil, assert.Nil, assert.False, assert.Nil, assert.NotNil, assert.Nil, nil, nil, nil},
		{"Invalid boards.txt", "invalid-boards.txt", assert.Nil, assert.NotNil, assert.False, assert.Nil, assert.NotNil, assert.Nil, nil, nil, nil},
		{"Missing boards.txt", "missing-boards.txt", assert.Nil, assert.NotNil, assert.False, assert.Nil, assert.NotNil, assert.Nil, nil, nil, nil},
		{"Valid platform.txt", "valid-platform.txt", assert.NotNil, assert.Nil, assert.True, assert.NotNil, assert.Nil, assert.NotNil, []string{"foo_discovery", "bar_discovery"}, map[string][]string{"avrdude": {"foo_field_name"}}, []string{"avrdude", "bossac"}},
		{"Invalid platform.txt", "invalid-platform.txt", assert.NotNil, assert.Nil, assert.True, assert.Nil, assert.NotNil, assert.Nil, nil, nil, nil},
		{"Missing platform.txt", "missing-platform.txt", assert.NotNil, assert.Nil, assert.False, assert.Nil, assert.NotNil, assert.Nil, nil, nil, nil},
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

		testTable.platformTxtExistsAssertion(t, PlatformTxtExists(), testTable.testName)
		testTable.platformTxtAssertion(t, PlatformTxt(), testTable.testName)
		testTable.platformTxtLoadErrorAssertion(t, PlatformTxtLoadError(), testTable.testName)
		testTable.platformTxtSchemaValidationResultAssertion(t, PlatformTxtSchemaValidationResult(), testTable.testName)
		assert.Equal(t, testTable.platformTxtPluggableDiscoveryNamesAssertion, PlatformTxtPluggableDiscoveryNames(), testTable.testName)
		assert.True(t, reflect.DeepEqual(testTable.platformTxtUserProvidedFieldNamesAssertion, PlatformTxtUserProvidedFieldNames()), testTable.testName)
		assert.Equal(t, testTable.platformTxtToolNamesAssertion, PlatformTxtToolNames(), testTable.testName)
	}
}
