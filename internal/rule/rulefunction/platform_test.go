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

func TestProgrammersTxtProgrammerIDNameMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-programmers.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-programmers.txt", ruleresult.NotRun, ""},
		{"No programmers", "no-programmers-programmers.txt", ruleresult.Skip, ""},
		{"Property missing", "programmerID-name-missing-programmers.txt", ruleresult.Fail, "foo, bar"},
		{"Valid", "valid-programmers.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(ProgrammersTxtProgrammerIDNameMissing, testTables, t)
}

func TestProgrammersTxtProgrammerIDNameLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-programmers.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-programmers.txt", ruleresult.NotRun, ""},
		{"No programmers", "no-programmers-programmers.txt", ruleresult.Skip, ""},
		{"Property LT min", "programmerID-name-LT-programmers.txt", ruleresult.Fail, "foo, bar"},
		{"Valid", "valid-programmers.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(ProgrammersTxtProgrammerIDNameLTMinLength, testTables, t)
}

func TestProgrammersTxtProgrammerIDProgramToolMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-programmers.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-programmers.txt", ruleresult.NotRun, ""},
		{"No programmers", "no-programmers-programmers.txt", ruleresult.Skip, ""},
		{"Property missing", "programmerID-program-tool-missing-programmers.txt", ruleresult.Fail, "foo, bar"},
		{"Valid", "valid-programmers.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(ProgrammersTxtProgrammerIDProgramToolMissing, testTables, t)
}

func TestProgrammersTxtProgrammerIDProgramToolLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-programmers.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-programmers.txt", ruleresult.NotRun, ""},
		{"No programmers", "no-programmers-programmers.txt", ruleresult.Skip, ""},
		{"Property LT min", "programmerID-program-tool-LT-programmers.txt", ruleresult.Fail, "foo, bar"},
		{"Valid", "valid-programmers.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(ProgrammersTxtProgrammerIDProgramToolLTMinLength, testTables, t)
}

func TestPlatformTxtFormat(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.Fail, ""},
	}

	checkPlatformRuleFunction(PlatformTxtFormat, testTables, t)
}

func TestPlatformTxtNameMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "name-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtNameMissing, testTables, t)
}

func TestPlatformTxtNameLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "name-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtNameLTMinLength, testTables, t)
}

func TestPlatformTxtVersionMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "version-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtVersionMissing, testTables, t)
}

func TestPlatformTxtVersionNonRelaxedSemver(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property invalid", "version-invalid-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtVersionNonRelaxedSemver, testTables, t)
}

func TestPlatformTxtVersionNonSemver(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property invalid", "version-non-semver-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtVersionNonSemver, testTables, t)
}

func TestPlatformTxtCompilerWarningFlagsNoneMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-warning_flags-none-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerWarningFlagsNoneMissing, testTables, t)
}

func TestPlatformTxtCompilerWarningFlagsDefaultMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-warning_flags-default-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerWarningFlagsDefaultMissing, testTables, t)
}

func TestPlatformTxtCompilerWarningFlagsMoreMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-warning_flags-more-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerWarningFlagsMoreMissing, testTables, t)
}

func TestPlatformTxtCompilerWarningFlagsAllMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-warning_flags-all-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerWarningFlagsAllMissing, testTables, t)
}

func TestPlatformTxtCompilerOptimizationFlagsDebugMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Dependent property not present", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "compiler-optimization_flags-debug-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerOptimizationFlagsDebugMissing, testTables, t)
}

func TestPlatformTxtCompilerOptimizationFlagsReleaseMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Dependent property not present", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "compiler-optimization_flags-release-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerOptimizationFlagsReleaseMissing, testTables, t)
}

func TestPlatformTxtCompilerCExtraFlagsMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-c-extra_flags-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCExtraFlagsMissing, testTables, t)
}

func TestPlatformTxtCompilerCExtraFlagsNotEmpty(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Not empty", "compiler-c-extra_flags-not-empty-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCExtraFlagsNotEmpty, testTables, t)
}

func TestPlatformTxtCompilerCppExtraFlagsMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-cpp-extra_flags-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCppExtraFlagsMissing, testTables, t)
}

func TestPlatformTxtCompilerCppExtraFlagsNotEmpty(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Not empty", "compiler-cpp-extra_flags-not-empty-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCppExtraFlagsNotEmpty, testTables, t)
}

func TestPlatformTxtCompilerSExtraFlagsMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-S-extra_flags-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerSExtraFlagsMissing, testTables, t)
}

func TestPlatformTxtCompilerSExtraFlagsNotEmpty(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Not empty", "compiler-S-extra_flags-not-empty-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerSExtraFlagsNotEmpty, testTables, t)
}

func TestPlatformTxtCompilerArExtraFlagsMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-ar-extra_flags-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerArExtraFlagsMissing, testTables, t)
}

func TestPlatformTxtCompilerArExtraFlagsNotEmpty(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Not empty", "compiler-ar-extra_flags-not-empty-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerArExtraFlagsNotEmpty, testTables, t)
}

func TestPlatformTxtCompilerCElfExtraFlagsMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "compiler-c-elf-extra_flags-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCElfExtraFlagsMissing, testTables, t)
}

func TestPlatformTxtCompilerCElfExtraFlagsNotEmpty(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Not empty", "compiler-c-elf-extra_flags-not-empty-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtCompilerCElfExtraFlagsNotEmpty, testTables, t)
}

func TestPlatformTxtRecipePreprocMacrosLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Property LT min", "recipe-preproc-macros-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipePreprocMacrosLTMinLength, testTables, t)
}

func TestPlatformTxtRecipePreprocMacrosExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-preproc-macros-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipePreprocMacrosExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeCOPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-c-o-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCOPatternMissing, testTables, t)
}

func TestPlatformTxtRecipeCOPatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-c-o-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCOPatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeCOPatternExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-c-o-pattern-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCOPatternExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeCppOPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-cpp-o-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCppOPatternMissing, testTables, t)
}

func TestPlatformTxtRecipeCppOPatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-cpp-o-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCppOPatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeCppOPatternExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-cpp-o-pattern-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCppOPatternExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeSOPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-S-o-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSOPatternMissing, testTables, t)
}

func TestPlatformTxtRecipeSOPatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-S-o-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSOPatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeSOPatternExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-S-o-pattern-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSOPatternExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeArPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-ar-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeArPatternMissing, testTables, t)
}

func TestPlatformTxtRecipeArPatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-ar-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeArPatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeArPatternExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-ar-pattern-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeArPatternExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeCCombinePatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-c-combine-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCCombinePatternMissing, testTables, t)
}

func TestPlatformTxtRecipeCCombinePatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-c-combine-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCCombinePatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeCCombinePatternExtraFlagsSupport(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"No extra flags support", "recipe-c-combine-pattern-no-extra-flags-support-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeCCombinePatternExtraFlagsSupport, testTables, t)
}

func TestPlatformTxtRecipeOutputTmpFileMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-output-tmp_file-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeOutputTmpFileMissing, testTables, t)
}

func TestPlatformTxtRecipeOutputTmpFileLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-output-tmp_file-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeOutputTmpFileLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeOutputSaveFileMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-output-save_file-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeOutputSaveFileMissing, testTables, t)
}

func TestPlatformTxtRecipeOutputSaveFileLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.NotRun, ""},
		{"Property LT min", "recipe-output-save_file-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeOutputSaveFileLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeSizePatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-size-pattern-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSizePatternMissing, testTables, t)
}

func TestPlatformTxtRecipeSizePatternLTMinLength(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "properties-missing-platform.txt", ruleresult.Skip, ""},
		{"Property LT min", "recipe-size-pattern-LT-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSizePatternLTMinLength, testTables, t)
}

func TestPlatformTxtRecipeSizeRegexMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-size-regex-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSizeRegexMissing, testTables, t)
}

func TestPlatformTxtRecipeSizeRegexDataMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"Property missing", "recipe-size-regex-data-missing-platform.txt", ruleresult.Fail, ""},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtRecipeSizeRegexDataMissing, testTables, t)
}

func TestPlatformTxtUploadPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "upload-pattern-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtUploadPatternMissing, testTables, t)
}

func TestPlatformTxtProgramParamsVerboseMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "program-params-verbose-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtProgramParamsVerboseMissing, testTables, t)
}

func TestPlatformTxtProgramParamsQuietMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "program-params-quiet-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtProgramParamsQuietMissing, testTables, t)
}

func TestPlatformTxtProgramPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "program-pattern-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtProgramPatternMissing, testTables, t)
}

func TestPlatformTxtEraseParamsVerboseMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "erase-params-verbose-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtEraseParamsVerboseMissing, testTables, t)
}

func TestPlatformTxtEraseParamsQuietMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "erase-params-quiet-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtEraseParamsQuietMissing, testTables, t)
}

func TestPlatformTxtErasePatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "erase-pattern-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtErasePatternMissing, testTables, t)
}

func TestPlatformTxtBootloaderParamsVerboseMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "bootloader-params-verbose-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtBootloaderParamsVerboseMissing, testTables, t)
}

func TestPlatformTxtBootloaderParamsQuietMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "bootloader-params-quiet-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtBootloaderParamsQuietMissing, testTables, t)
}

func TestPlatformTxtBootloaderPatternMissing(t *testing.T) {
	testTables := []platformRuleFunctionTestTable{
		{"Missing", "missing-platform.txt", ruleresult.Skip, ""},
		{"Invalid", "invalid-platform.txt", ruleresult.NotRun, ""},
		{"No tools", "no-tools-platform.txt", ruleresult.Skip, ""},
		{"Property missing", "bootloader-pattern-missing-platform.txt", ruleresult.Fail, "avrdude, bossac"},
		{"Valid", "valid-platform.txt", ruleresult.Pass, ""},
	}

	checkPlatformRuleFunction(PlatformTxtBootloaderPatternMissing, testTables, t)
}
