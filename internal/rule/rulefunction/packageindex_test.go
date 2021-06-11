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

package rulefunction

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var packageIndexesTestDataPath *paths.Path

func init() {
	workingDirectory, _ := paths.Getwd()
	packageIndexesTestDataPath = workingDirectory.Join("testdata", "packageindexes")
}

type packageIndexRuleFunctionTestTable struct {
	testName               string
	packageIndexFolderName string
	expectedRuleResult     ruleresult.Type
	expectedOutputQuery    string
}

func checkPackageIndexRuleFunction(ruleFunction Type, testTables []packageIndexRuleFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             packageIndexesTestDataPath.Join(testTable.packageIndexFolderName),
			ProjectType:      projecttype.PackageIndex,
			SuperprojectType: projecttype.PackageIndex,
		}

		projectdata.Initialize(testProject)

		result, output := ruleFunction()
		assert.Equal(t, testTable.expectedRuleResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), fmt.Sprintf("%s (output: %s, assertion regex: %s)", testTable.testName, output, testTable.expectedOutputQuery))
	}
}

func TestPackageIndexMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Missing", "missing", ruleresult.Fail, ""},
		{"Present", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexMissing, testTables, t)
}

func TestPackageIndexFilenameInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Missing", "missing", ruleresult.NotRun, ""},
		{"Valid 3rd party", "3rd-party-filename", ruleresult.Pass, ""},
		{"Valid official", "official-filename", ruleresult.Fail, "^package_index.json$"},
		{"Invalid", "invalid-filename", ruleresult.Fail, "^invalid-filename.json$"},
	}

	checkPackageIndexRuleFunction(PackageIndexFilenameInvalid, testTables, t)
}

func TestPackageIndexOfficialFilenameInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Missing", "missing", ruleresult.NotRun, ""},
		{"Valid 3rd party", "3rd-party-filename", ruleresult.Pass, ""},
		{"Valid official", "official-filename", ruleresult.Pass, ""},
		{"Invalid", "invalid-filename", ruleresult.Fail, "^invalid-filename.json$"},
	}

	checkPackageIndexRuleFunction(PackageIndexOfficialFilenameInvalid, testTables, t)
}

func TestPackageIndexJSONFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.Fail, ""},
		{"Not valid package index", "invalid-package-index", ruleresult.Pass, ""},
		{"Valid package index", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexJSONFormat, testTables, t)
}

func TestPackageIndexFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.Fail, ""},
		{"Not valid package index", "invalid-package-index", ruleresult.Fail, ""},
		{"Valid package index", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexFormat, testTables, t)
}

func TestPackageIndexPackagesWebsiteURLDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-websiteurl-dead", ruleresult.Fail, "^foopackager1, foopackager2$"},
		{"Invalid URL", "packages-websiteurl-invalid", ruleresult.Fail, "^foopackager$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesWebsiteURLDeadLink, testTables, t)
}

func TestPackageIndexPackagesHelpOnlineDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-help-online-dead", ruleresult.Fail, "^foopackager1, foopackager2$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpOnlineDeadLink, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpOnlineDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-platforms-help-online-dead", ruleresult.Fail, "^foopackager:avr@1\\.0\\.0, foopackager:samd@1\\.0\\.0$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpOnlineDeadLink, testTables, t)
}

func TestPackageIndexPackagesPlatformsURLDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-platforms-url-dead", ruleresult.Fail, "^foopackager:avr@1\\.0\\.0, foopackager:samd@1\\.0\\.0$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsURLDeadLink, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsURLDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-tools-systems-url-dead", ruleresult.Fail, "^foopackager:CMSIS@4\\.0\\.0-atmel - arm-linux-gnueabihf, foopackager:CMSIS@4\\.0\\.0-atmel - i686-mingw32$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsURLDeadLink, testTables, t)
}
