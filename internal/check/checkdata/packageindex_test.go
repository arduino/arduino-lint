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

package checkdata

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var packageIndexTestDataPath *paths.Path

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	packageIndexTestDataPath = workingDirectory.Join("testdata", "packageindexes")
}

func TestInitializeForPackageIndex(t *testing.T) {
	testTables := []struct {
		testName                       string
		path                           *paths.Path
		packageIndexAssertion          assert.ValueAssertionFunc
		packageIndexLoadErrorAssertion assert.ValueAssertionFunc
	}{
		{"Valid", packageIndexTestDataPath.Join("valid-package-index", "package_foo_index.json"), assert.NotNil, assert.Nil},
		{"Invalid package index", packageIndexTestDataPath.Join("invalid-package-index", "package_foo_index.json"), assert.Nil, assert.NotNil},
		{"Invalid JSON", packageIndexTestDataPath.Join("invalid-JSON", "package_foo_index.json"), assert.Nil, assert.NotNil},
	}

	for _, testTable := range testTables {

		testProject := project.Type{
			Path:             testTable.path,
			ProjectType:      projecttype.PackageIndex,
			SuperprojectType: projecttype.PackageIndex,
		}
		Initialize(testProject)

		testTable.packageIndexLoadErrorAssertion(t, PackageIndexLoadError(), testTable.testName)
		if PackageIndexLoadError() == nil {
			testTable.packageIndexAssertion(t, PackageIndex(), testTable.testName)
		}
	}
}
