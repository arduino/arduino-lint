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

func TestPackageIndexAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional root properties", "root-additional-properties", ruleresult.Fail, ""},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Packages key missing", "packages-missing", ruleresult.Fail, ""},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesMissing, testTables, t)
}

func TestPackageIndexPackagesIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages type", "packages-incorrect-type", ruleresult.Fail, ""},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesIncorrectType, testTables, t)
}

func TestPackageIndexPackagesAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages properties", "packages-additional-properties", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].name missing", "packages-name-missing", ruleresult.Fail, "^/packages/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesNameMissing, testTables, t)
}

func TestPackageIndexPackagesNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].name type", "packages-name-incorrect-type", ruleresult.Fail, "^/packages/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].name < min length", "packages-name-length-lt", ruleresult.Fail, "^/packages/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesNameIsArduino(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].name is arduino", "packages-name-is-arduino", ruleresult.Fail, "^/packages/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesNameIsArduino, testTables, t)
}

func TestPackageIndexPackagesMaintainerMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].maintainer missing", "packages-maintainer-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesMaintainerMissing, testTables, t)
}

func TestPackageIndexPackagesMaintainerIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].maintainer type", "packages-maintainer-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesMaintainerIncorrectType, testTables, t)
}

func TestPackageIndexPackagesMaintainerLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].maintainer < min length", "packages-maintainer-length-lt", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesMaintainerLTMinLength, testTables, t)
}

func TestPackageIndexPackagesMaintainerStartsWithArduino(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].maintainer starts with arduino", "packages-maintainer-starts-with-arduino", ruleresult.Fail, "^/packages/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesMaintainerStartsWithArduino, testTables, t)
}

func TestPackageIndexPackagesWebsiteURLMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].websiteURL missing", "packages-websiteurl-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesWebsiteURLMissing, testTables, t)
}

func TestPackageIndexPackagesWebsiteURLIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].websiteURL type", "packages-websiteurl-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesWebsiteURLIncorrectType, testTables, t)
}

func TestPackageIndexPackagesWebsiteURLInvalidFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].websiteURL format", "packages-websiteurl-invalid-format", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesWebsiteURLInvalidFormat, testTables, t)
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

func TestPackageIndexPackagesEmailMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].email missing", "packages-email-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesEmailMissing, testTables, t)
}

func TestPackageIndexPackagesEmailIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].email type", "packages-email-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesEmailIncorrectType, testTables, t)
}

func TestPackageIndexPackagesHelpIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].help type", "packages-help-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpIncorrectType, testTables, t)
}

func TestPackageIndexPackagesHelpAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].help properties", "packages-help-additional-properties", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesHelpOnlineMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].help.online missing", "packages-help-online-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpOnlineMissing, testTables, t)
}

func TestPackageIndexPackagesHelpOnlineIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].help.online type", "packages-help-online-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpOnlineIncorrectType, testTables, t)
}

func TestPackageIndexPackagesHelpOnlineInvalidFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].help.online format", "packages-help-online-invalid-format", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpOnlineInvalidFormat, testTables, t)
}

func TestPackageIndexPackagesHelpOnlineDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-help-online-dead", ruleresult.Fail, "^foopackager1, foopackager2$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesHelpOnlineDeadLink, testTables, t)
}

func TestPackageIndexPackagesPlatformsMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms missing", "packages-platforms-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms type", "packages-platforms-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].platforms[] properties", "packages-platforms-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesPlatformsNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].name missing", "packages-platforms-name-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsNameMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].name type", "packages-platforms-name-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].name < min length", "packages-platforms-name-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchitectureMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].architecture missing", "packages-platforms-architecture-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchitectureMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchitectureIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].architecture type", "packages-platforms-architecture-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchitectureIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchitectureLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].architecture < min length", "packages-platforms-architecture-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchitectureLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsVersionMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].architecture missing", "packages-platforms-version-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsVersionMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsVersionIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].version type", "packages-platforms-version-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsVersionIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsVersionNonRelaxedSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].version not relaxed semver", "packages-platforms-version-non-relaxed-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@foo$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsVersionNonRelaxedSemver, testTables, t)
}

func TestPackageIndexPackagesPlatformsVersionNonSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].version not semver", "packages-platforms-version-not-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsVersionNonSemver, testTables, t)
}

func TestPackageIndexPackagesPlatformsDeprecatedIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].deprecated type", "packages-platforms-deprecated-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDeprecatedIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsCategoryMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].category missing", "packages-platforms-category-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsCategoryMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsCategoryIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].category type", "packages-platforms-category-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsCategoryIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsCategoryThirdPartyInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].category not valid for 3rd party", "packages-platforms-category-non-third-party", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsCategoryThirdPartyInvalid, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].help missing", "packages-platforms-help-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].help type", "packages-platforms-help-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].platforms[].help properties", "packages-platforms-help-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpOnlineMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].help.online missing", "packages-platforms-help-online-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpOnlineMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpOnlineIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].help.online type", "packages-platforms-help-online-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpOnlineIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpOnlineInvalidFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].help.online format", "packages-platforms-help-online-invalid-format", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpOnlineInvalidFormat, testTables, t)
}

func TestPackageIndexPackagesPlatformsHelpOnlineDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-platforms-help-online-dead", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0\n" + brokenOutputListIndent + "foopackager:samd@1\\.0\\.0$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsHelpOnlineDeadLink, testTables, t)
}

func TestPackageIndexPackagesPlatformsUrlMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].url missing", "packages-platforms-url-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsURLMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsUrlIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].url type", "packages-platforms-url-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsURLIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsUrlInvalidFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].url format", "packages-platforms-url-invalid-format", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsURLInvalidFormat, testTables, t)
}

func TestPackageIndexPackagesPlatformsURLDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-platforms-url-dead", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0\n" + brokenOutputListIndent + "foopackager:samd@1\\.0\\.0$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsURLDeadLink, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchiveFileNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].archiveFileName missing", "packages-platforms-archivefilename-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchiveFileNameMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchiveFileNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].archiveFileName type", "packages-platforms-archivefilename-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchiveFileNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchiveFileNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].archiveFileName < min length", "packages-platforms-archivefilename-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchiveFileNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsArchiveFileNameInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid filename", "packages-platforms-archivefilename-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsArchiveFileNameInvalid, testTables, t)
}

func TestPackageIndexPackagesPlatformsChecksumMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].checksum missing", "packages-platforms-checksum-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsChecksumMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsChecksumIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].checksum type", "packages-platforms-checksum-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsChecksumIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsChecksumInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].platforms[].checksum format", "packages-platforms-checksum-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsChecksumInvalid, testTables, t)
}

func TestPackageIndexPackagesPlatformsChecksumDiscouragedAlgorithm(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].checksum uses discouraged algorithm", "packages-platforms-checksum-discouraged", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsChecksumDiscouragedAlgorithm, testTables, t)
}

func TestPackageIndexPackagesPlatformsSizeMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].size missing", "packages-platforms-size-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsSizeMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsSizeIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].size type", "packages-platforms-size-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsSizeIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsSizeInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].platforms[].size format", "packages-platforms-size-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsSizeInvalid, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].boards[] missing", "packages-platforms-boards-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].boards type", "packages-platforms-boards-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].platforms[].boards[] properties", "packages-platforms-boards-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0 >> My Board$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].boards[].name missing", "packages-platforms-boards-name-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/boards/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsNameMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].boards[].name type", "packages-platforms-boards-name-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/boards/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsBoardsNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].boards[].name < min length", "packages-platforms-boards-name-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/boards/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsBoardsNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies missing", "packages-platforms-toolsdependencies-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].toolsDependencies type", "packages-platforms-toolsdependencies-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].platforms[].toolsDependencies[] properties", "packages-platforms-toolsdependencies-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0 >> arduino:avr-gcc@7\\.3\\.0-atmel3\\.6\\.1-arduino7$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesPackagerMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].packager missing", "packages-platforms-toolsdependencies-packager-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesPackagerMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesPackagerIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].toolsDependencies[].packager type", "packages-platforms-toolsdependencies-packager-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesPackagerIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesPackagerLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].packager < min length", "packages-platforms-toolsdependencies-packager-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesPackagerLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].name missing", "packages-platforms-toolsdependencies-name-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesNameMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].toolsDependencies[].name type", "packages-platforms-toolsdependencies-name-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].name < min length", "packages-platforms-toolsdependencies-name-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesVersionMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].version missing", "packages-platforms-toolsdependencies-version-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesVersionMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesVersionIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].toolsDependencies[].version type", "packages-platforms-toolsdependencies-version-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/toolsDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesVersionIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesVersionNonRelaxedSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].version not relaxed semver", "packages-platforms-toolsdependencies-version-non-relaxed-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0 >> arduino:avr-gcc@foo$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesVersionNonRelaxedSemver, testTables, t)
}

func TestPackageIndexPackagesPlatformsToolsDependenciesVersionNonSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].toolsDependencies[].version not semver", "packages-platforms-toolsdependencies-version-not-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0 >> arduino:avr-gcc@7\\.3$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsToolsDependenciesVersionNonSemver, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].discoveryDependencies type", "packages-platforms-discoverydependencies-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].platforms[].discoveryDependencies[] properties", "packages-platforms-discoverydependencies-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:avr@1\\.0\\.0 >> arduino:ble-discovery$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesPackagerMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].discoveryDependencies[].packager missing", "packages-platforms-discoverydependencies-packager-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesPackagerIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].discoveryDependencies[].packager type", "packages-platforms-discoverydependencies-packager-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesPackagerLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].discoveryDependencies[].packager < min length", "packages-platforms-discoverydependencies-packager-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerLTMinLength, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].discoveryDependencies[].name missing", "packages-platforms-discoverydependencies-name-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesNameMissing, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].platforms[].discoveryDependencies[].name type", "packages-platforms-discoverydependencies-name-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesPlatformsDiscoveryDependenciesNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].platforms[].discoveryDependencies[].name < min length", "packages-platforms-discoverydependencies-name-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/platforms/0/discoveryDependencies/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesPlatformsDiscoveryDependenciesNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesToolsMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools missing", "packages-tools-missing", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsMissing, testTables, t)
}

func TestPackageIndexPackagesToolsIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools type", "packages-tools-incorrect-type", ruleresult.Fail, "^foopackager$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].tools[] properties", "packages-tools-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesToolsNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].name missing", "packages-tools-name-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsNameMissing, testTables, t)
}

func TestPackageIndexPackagesToolsNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools.name type", "packages-tools-name-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].name < min length", "packages-tools-name-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesToolsVersionMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].version missing", "packages-tools-version-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsVersionMissing, testTables, t)
}

func TestPackageIndexPackagesToolsVersionIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].version type", "packages-tools-version-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsVersionIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsVersionNonRelaxedSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].version not relaxed semver", "packages-tools-version-non-relaxed-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@foo$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsVersionNonRelaxedSemver, testTables, t)
}

func TestPackageIndexPackagesToolsVersionNonSemver(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].version not semver", "packages-tools-version-not-semver", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@1.0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsVersionNonSemver, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[] missing", "packages-tools-systems-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems type", "packages-tools-systems-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsAdditionalProperties(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Additional packages[].tools[].systems[] properties", "packages-tools-systems-additional-properties", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsAdditionalProperties, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsHostMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].host missing", "packages-tools-systems-host-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0/systems/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsHostMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsHostIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].host type", "packages-tools-systems-host-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "/packages/0/tools/0/systems/0$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsHostIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsHostInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].tools[].systems[].host format", "packages-tools-systems-host-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> foo$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsHostInvalid, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsUrlMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].url missing", "packages-tools-systems-url-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsURLMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsUrlIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].url type", "packages-tools-systems-url-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsURLIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsUrlInvalidFormat(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].url format", "packages-tools-systems-url-invalid-format", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsURLInvalidFormat, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsURLDeadLink(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Dead URLs", "packages-tools-systems-url-dead", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:CMSIS@4\\.0\\.0-atmel >> arm-linux-gnueabihf\n" + brokenOutputListIndent + "foopackager:CMSIS@4\\.0\\.0-atmel >> i686-mingw32$"},
		{"Valid URL", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsURLDeadLink, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsArchiveFileNameMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].archiveFileName missing", "packages-tools-systems-archivefilename-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsArchiveFileNameMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsArchiveFileNameIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].archiveFileName type", "packages-tools-systems-archivefilename-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsArchiveFileNameIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsArchiveFileNameLTMinLength(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].archiveFileName < min length", "packages-tools-systems-archivefilename-length-lt", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsArchiveFileNameLTMinLength, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsArchiveFileNameInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].tools[].systems[].archiveFileName format", "packages-tools-systems-archivefilename-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsArchiveFileNameInvalid, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsChecksumMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].checksum missing", "packages-tools-systems-checksum-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsChecksumMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsChecksumIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].checksum type", "packages-tools-systems-checksum-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsChecksumIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsChecksumInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].tools[].systems[].checksum format", "packages-tools-systems-checksum-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsChecksumInvalid, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsChecksumDiscouragedAlgorithm(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].checksum uses discouraged algorithm", "packages-tools-systems-checksum-discouraged", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsChecksumDiscouragedAlgorithm, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsSizeMissing(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"packages[].tools[].systems[].size missing", "packages-tools-systems-size-missing", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsSizeMissing, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsSizeIncorrectType(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Incorrect packages[].tools[].systems[].size type", "packages-tools-systems-size-incorrect-type", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsSizeIncorrectType, testTables, t)
}

func TestPackageIndexPackagesToolsSystemsSizeInvalid(t *testing.T) {
	testTables := []packageIndexRuleFunctionTestTable{
		{"Invalid JSON", "invalid-JSON", ruleresult.NotRun, ""},
		{"Invalid packages[].tools[].systems[].size format", "packages-tools-systems-size-invalid", ruleresult.Fail, "^" + brokenOutputListIndent + "foopackager:openocd@0\\.11\\.0-arduino2 >> aarch64-linux-gnu$"},
		{"Valid", "valid-package-index", ruleresult.Pass, ""},
	}

	checkPackageIndexRuleFunction(PackageIndexPackagesToolsSystemsSizeInvalid, testTables, t)
}
