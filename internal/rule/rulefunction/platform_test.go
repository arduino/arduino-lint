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

package rulefunction

import (
	"regexp"
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
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

type platformRuleFunctionTestTable struct {
	testName            string
	platformFolderName  string
	expectedRuleResult  ruleresult.Type
	expectedOutputQuery string
}

func checkPlatformRuleFunction(ruleFunction Type, testTables []platformRuleFunctionTestTable, t *testing.T) {
	for _, testTable := range testTables {
		expectedOutputRegexp := regexp.MustCompile(testTable.expectedOutputQuery)

		testProject := project.Type{
			Path:             platformTestDataPath.Join(testTable.platformFolderName),
			ProjectType:      projecttype.Platform,
			SuperprojectType: projecttype.Platform,
		}

		projectdata.Initialize(testProject)

		result, output := ruleFunction()
		assert.Equal(t, testTable.expectedRuleResult, result, testTable.testName)
		assert.True(t, expectedOutputRegexp.MatchString(output), testTable.testName)
	}
}

func TestBoardsTxtMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Present", "valid-boards.txt", ruleresult.Pass, ""},
		{"Missing", "missing-boards.txt", ruleresult.Fail, ""},
	}

	checkPlatformRuleFunction(BoardsTxtMissing, testTables, t)
}

func TestBoardsTxtFormat(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.Fail, ""},
	}

	checkPlatformRuleFunction(BoardsTxtFormat, testTables, t)
}

func TestBoardsTxtBoardIDNameMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-name-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDNameMissing, testTables, t)
}

func TestBoardsTxtBoardIDNameLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-name-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDNameLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDBuildBoardMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-build-board-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDBuildBoardMissing, testTables, t)
}

func TestBoardsTxtBoardIDBuildBoardLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-build-board-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDBuildBoardLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDBuildCoreMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-build-core-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDBuildCoreMissing, testTables, t)
}

func TestBoardsTxtBoardIDBuildCoreLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-build-core-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDBuildCoreLTMinLength, testTables, t)
}

func TestBoardsTxtUserExtraFlagsUsage(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Use of user extra flags", "boardID-compiler-x-extra_flags-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtUserExtraFlagsUsage, testTables, t)
}

func TestBoardsTxtBoardIDDebugToolLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-debug-tool-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDDebugToolLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDHideInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-hide-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDHideInvalid, testTables, t)
}

func TestBoardsTxtMenuMenuIDLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No menus", "no-menus-boards.txt", ruleresult.Skip, ""},
		{"Menu title too short", "menu-menuID-LT-boards.txt", ruleresult.Fail, "foo, baz"},
		{"Menu title valid", "menu-menuID-valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtMenuMenuIDLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDMenuMenuIDOptionIDLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-menu-menuID-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDMenuMenuIDOptionIDLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDSerialDisableDTRInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-serial-disableDTR-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDSerialDisableDTRInvalid, testTables, t)
}

func TestBoardsTxtBoardIDSerialDisableRTSInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-serial-disableRTS-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDSerialDisableRTSInvalid, testTables, t)
}

func TestBoardsTxtBoardIDUploadToolMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-upload-tool-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadToolMissing, testTables, t)
}

func TestBoardsTxtBoardIDUploadToolLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property LT min", "boardID-upload-tool-LT-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadToolLTMinLength, testTables, t)
}

func TestBoardsTxtBoardIDUploadMaximumSizeMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-upload-maximum_size-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadMaximumSizeMissing, testTables, t)
}

func TestBoardsTxtBoardIDUploadMaximumSizeInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-upload-maximum_size-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadMaximumSizeInvalid, testTables, t)
}

func TestBoardsTxtBoardIDUploadMaximumDataSizeMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property missing", "boardID-upload-maximum_data_size-missing-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadMaximumDataSizeMissing, testTables, t)
}

func TestBoardsTxtBoardIDUploadMaximumDataSizeInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-upload-maximum_data_size-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadMaximumDataSizeInvalid, testTables, t)
}

func TestBoardsTxtBoardIDUploadUse1200bpsTouchInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-upload-use_1200bps_touch-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadUse1200bpsTouchInvalid, testTables, t)
}

func TestBoardsTxtBoardIDUploadWaitForUploadPortInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-upload-wait_for_upload_port-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDUploadWaitForUploadPortInvalid, testTables, t)
}

func TestBoardsTxtBoardIDVidNInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-vid-n-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDVidNInvalid, testTables, t)
}

func TestBoardsTxtBoardIDPidNInvalid(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-boards.txt", ruleresult.NotRun, ""},
		{"Invalid", "invalid-boards.txt", ruleresult.NotRun, ""},
		{"No boards", "no-boards-boards.txt", ruleresult.Skip, ""},
		{"Property invalid", "boardID-pid-n-invalid-boards.txt", ruleresult.Fail, "buno, funo"},
		{"Valid", "valid-boards.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(BoardsTxtBoardIDPidNInvalid, testTables, t)
}
