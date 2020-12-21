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

package packageindex

import (
	"os"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = paths.New(workingDirectory, "testdata")
}

func TestHasValidExtension(t *testing.T) {
	assert.True(t, HasValidExtension(paths.New("/foo", "bar.json")))
	assert.False(t, HasValidExtension(paths.New("/foo", "bar.baz")))
}

func TestHasValidFilename(t *testing.T) {
	testTables := []struct {
		testName         string
		filename         string
		officialRuleMode bool
		assertion        assert.BoolAssertionFunc
	}{
		{"Official, primary", "package_index.json", true, assert.True},
		{"Official, secondary", "package_foo_index.json", true, assert.True},
		{"Official, invalid", "packageindex.json", true, assert.False},
		{"Unofficial, valid", "package_foo_index.json", false, assert.True},
		{"Unofficial, official", "package_index.json", false, assert.False},
		{"Unofficial, invalid", "packageindex.json", false, assert.False},
	}

	for _, testTable := range testTables {
		testTable.assertion(t, HasValidFilename(paths.New("/foo", testTable.filename), testTable.officialRuleMode), testTable.testName)
	}
}

func TestFind(t *testing.T) {
	testTables := []struct {
		testName     string
		path         *paths.Path
		expectedPath *paths.Path
		errAssertion assert.ValueAssertionFunc
	}{
		{"Nonexistent", testDataPath.Join("nonexistent"), nil, assert.NotNil},
		{"File", testDataPath.Join("HasPackageIndex", "package_foo_index.json"), testDataPath.Join("HasPackageIndex", "package_foo_index.json"), assert.Nil},
		{"Single", testDataPath.Join("HasPackageIndex"), testDataPath.Join("HasPackageIndex", "package_foo_index.json"), assert.Nil},
		{"Multiple", testDataPath.Join("HasMultiple"), testDataPath.Join("HasMultiple", "package_foo_index.json"), assert.Nil},
		{"Valid extension fallback", testDataPath.Join("HasJSON"), testDataPath.Join("HasJSON", "foo.json"), assert.Nil},
		{"None", testDataPath.Join("HasNone"), nil, assert.NotNil},
	}

	for _, testTable := range testTables {
		path, err := Find(testTable.path)
		testTable.errAssertion(t, err, testTable.testName)
		if err == nil {
			assert.Equal(t, testTable.expectedPath, path, testTable.testName)
		}
	}
}
