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
	"strings"

	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
)

// The rule functions for platforms.

// BoardsTxtMissing checks whether the platform contains a boards.txt
func BoardsTxtMissing() (result ruleresult.Type, output string) {
	boardsTxtPath := projectdata.ProjectPath().Join("boards.txt")
	exist, err := boardsTxtPath.ExistCheck()
	if err != nil {
		panic(err)
	}

	if exist {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, boardsTxtPath.String()
}

// BoardsTxtFormat checks for invalid boards.txt format.
func BoardsTxtFormat() (result ruleresult.Type, output string) {
	if !projectdata.ProjectPath().Join("boards.txt").Exist() {
		return ruleresult.NotRun, "boards.txt missing"
	}

	if projectdata.BoardsTxtLoadError() == nil {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.BoardsTxtLoadError().Error()
}

// BoardsTxtBoardIDNameMissing checks if any of the boards are missing name properties.
func BoardsTxtBoardIDNameMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("name", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDNameLTMinLength checks if any of the board names are less than the minimum length.
func BoardsTxtBoardIDNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueLTMinLength("name", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDBuildBoardMissing checks if any of the boards are missing build.board properties.
func BoardsTxtBoardIDBuildBoardMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("build\\.board", compliancelevel.Strict)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDBuildBoardLTMinLength checks if any of the board build.board values are less than the minimum length.
func BoardsTxtBoardIDBuildBoardLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueLTMinLength("build\\.board", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDBuildCoreMissing checks if any of the boards are missing build.core properties.
func BoardsTxtBoardIDBuildCoreMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("build\\.core", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDBuildCoreLTMinLength checks if any of the board build.core values are less than the minimum length.
func BoardsTxtBoardIDBuildCoreLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueLTMinLength("build\\.core", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtUserExtraFlagsUsage checks if the user's compiler.x.extra_flags properties are used in boards.txt.
func BoardsTxtUserExtraFlagsUsage() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := []string{}
	for _, boardID := range projectdata.BoardsTxtBoardIds() {
		if schema.ValidationErrorMatch("#/"+boardID, "/userExtraFlagsProperties/", "", "", projectdata.BoardsTxtSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantBoardIDs = append(nonCompliantBoardIDs, boardID)
		}
	}

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDHideInvalid checks if any of the board hide values are less than the minimum length.
func BoardsTxtBoardIDHideInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueEnumMismatch("hide", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtMenuMenuIDLTMinLength checks if any of the menu titles are less than the minimum length.
func BoardsTxtMenuMenuIDLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtMenuIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no menus"
	}

	nonCompliantMenuIDs := []string{}
	for _, menuID := range projectdata.BoardsTxtMenuIds() {
		if schema.PropertyLessThanMinLength("menu/"+menuID, projectdata.BoardsTxtSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantMenuIDs = append(nonCompliantMenuIDs, menuID)
		}
	}

	if len(nonCompliantMenuIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantMenuIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDMenuMenuIDOptionIDLTMinLength checks if any of the board menu.MENU_ID.OPTION_ID values are less than the minimum length.
func BoardsTxtBoardIDMenuMenuIDOptionIDLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueLTMinLength("menu\\..+\\..+", compliancelevel.Strict)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDSerialDisableDTRInvalid checks if any of the board serial.disableDTR values are invalid.
func BoardsTxtBoardIDSerialDisableDTRInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueEnumMismatch("serial\\.disableDTR", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDSerialDisableRTSInvalid checks if any of the board serial.disableRTS values are invalid.
func BoardsTxtBoardIDSerialDisableRTSInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueEnumMismatch("serial\\.disableRTS", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadToolMissing checks if any of the boards are missing upload.tool properties.
func BoardsTxtBoardIDUploadToolMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("upload\\.tool", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadToolLTMinLength checks if any of the board upload.tool values are less than the minimum length.
func BoardsTxtBoardIDUploadToolLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueLTMinLength("upload\\.tool", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadMaximumSizeMissing checks if any of the boards are missing upload.maximum_size properties.
func BoardsTxtBoardIDUploadMaximumSizeMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("upload\\.maximum_size", compliancelevel.Strict)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadMaximumSizeInvalid checks if any of the board upload.maximum_size values have an invalid format.
func BoardsTxtBoardIDUploadMaximumSizeInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValuePatternMismatch("upload\\.maximum_size", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadMaximumDataSizeMissing checks if any of the boards are missing upload.maximum_data_size properties.
func BoardsTxtBoardIDUploadMaximumDataSizeMissing() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDMissingRequiredProperty("upload\\.maximum_data_size", compliancelevel.Strict)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadMaximumDataSizeInvalid checks if any of the board upload.maximum_data_size values have an invalid format.
func BoardsTxtBoardIDUploadMaximumDataSizeInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValuePatternMismatch("upload\\.maximum_data_size", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadUse1200bpsTouchInvalid checks if any of the board upload.use_1200bps_touch values are invalid.
func BoardsTxtBoardIDUploadUse1200bpsTouchInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueEnumMismatch("upload\\.use_1200bps_touch", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDUploadWaitForUploadPortInvalid checks if any of the board upload.wait_for_upload_port values are invalid.
func BoardsTxtBoardIDUploadWaitForUploadPortInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValueEnumMismatch("upload\\.wait_for_upload_port", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDVidNInvalid checks if any of the board vid.n values have an invalid format.
func BoardsTxtBoardIDVidNInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValuePatternMismatch("vid\\.[0-9]+", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// BoardsTxtBoardIDPidNInvalid checks if any of the board pid.n values have an invalid format.
func BoardsTxtBoardIDPidNInvalid() (result ruleresult.Type, output string) {
	if projectdata.BoardsTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load boards.txt"
	}

	if len(projectdata.BoardsTxtBoardIds()) == 0 {
		return ruleresult.Skip, "boards.txt has no boards"
	}

	nonCompliantBoardIDs := boardIDValuePatternMismatch("pid\\.[0-9]+", compliancelevel.Specification)

	if len(nonCompliantBoardIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantBoardIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// ProgrammersTxtFormat checks for invalid programmers.txt format.
func ProgrammersTxtFormat() (result ruleresult.Type, output string) {
	if !projectdata.ProgrammersTxtExists() {
		return ruleresult.Skip, "Platform has no programmers.txt"
	}

	if projectdata.ProgrammersTxtLoadError() == nil {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.ProgrammersTxtLoadError().Error()
}

// ProgrammersTxtProgrammerIDNameMissing checks if any of the programmers are missing name properties.
func ProgrammersTxtProgrammerIDNameMissing() (result ruleresult.Type, output string) {
	if !projectdata.ProgrammersTxtExists() {
		return ruleresult.Skip, "Platform has no programmers.txt"
	}

	if projectdata.ProgrammersTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load programmers.txt"
	}

	if len(projectdata.ProgrammersTxtProgrammerIds()) == 0 {
		return ruleresult.Skip, "programmers.txt has no programmers"
	}

	nonCompliantProgrammerIDs := programmerIDMissingRequiredProperty("name", compliancelevel.Specification)

	if len(nonCompliantProgrammerIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantProgrammerIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// ProgrammersTxtProgrammerIDNameLTMinLength checks if any of the programmer names are less than the minimum length.
func ProgrammersTxtProgrammerIDNameLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.ProgrammersTxtExists() {
		return ruleresult.Skip, "Platform has no programmers.txt"
	}

	if projectdata.ProgrammersTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load programmers.txt"
	}

	if len(projectdata.ProgrammersTxtProgrammerIds()) == 0 {
		return ruleresult.Skip, "programmers.txt has no programmers"
	}

	nonCompliantProgrammerIDs := programmerIDValueLTMinLength("name", compliancelevel.Specification)

	if len(nonCompliantProgrammerIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantProgrammerIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// ProgrammersTxtProgrammerIDProgramToolMissing checks if any of the programmers are missing program.tool properties.
func ProgrammersTxtProgrammerIDProgramToolMissing() (result ruleresult.Type, output string) {
	if !projectdata.ProgrammersTxtExists() {
		return ruleresult.Skip, "Platform has no programmers.txt"
	}

	if projectdata.ProgrammersTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load programmers.txt"
	}

	if len(projectdata.ProgrammersTxtProgrammerIds()) == 0 {
		return ruleresult.Skip, "programmers.txt has no programmers"
	}

	nonCompliantProgrammerIDs := programmerIDMissingRequiredProperty("program\\.tool", compliancelevel.Specification)

	if len(nonCompliantProgrammerIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantProgrammerIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// ProgrammersTxtProgrammerIDProgramToolLTMinLength checks if any of the programmer program.tool properties are less than the minimum length.
func ProgrammersTxtProgrammerIDProgramToolLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.ProgrammersTxtExists() {
		return ruleresult.Skip, "Platform has no programmers.txt"
	}

	if projectdata.ProgrammersTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load programmers.txt"
	}

	if len(projectdata.ProgrammersTxtProgrammerIds()) == 0 {
		return ruleresult.Skip, "programmers.txt has no programmers"
	}

	nonCompliantProgrammerIDs := programmerIDValueLTMinLength("program\\.tool", compliancelevel.Specification)

	if len(nonCompliantProgrammerIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantProgrammerIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtFormat checks for invalid platform.txt format.
func PlatformTxtFormat() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() == nil {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.PlatformTxtLoadError().Error()
}

// PlatformTxtNameMissing checks for missing name property in platform.txt.
func PlatformTxtNameMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("name", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtNameLTMinLength checks if the platform.txt name property value is less than the minimum length.
func PlatformTxtNameLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("name") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("name", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtVersionMissing checks for missing version property in platform.txt.
func PlatformTxtVersionMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("version", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtVersionNonRelaxedSemver checks whether the platform.txt version property is "relaxed semver" compliant.
func PlatformTxtVersionNonRelaxedSemver() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	version, ok := projectdata.PlatformTxt().GetOk("version")
	if !ok {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("version", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, version
	}

	return ruleresult.Pass, ""
}

// PlatformTxtVersionNonSemver checks whether the platform.txt version property is semver compliant.
func PlatformTxtVersionNonSemver() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	version, ok := projectdata.PlatformTxt().GetOk("version")
	if !ok {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("version", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, version
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerWarningFlagsNoneMissing checks for missing compiler.warning_flags.none property in platform.txt.
func PlatformTxtCompilerWarningFlagsNoneMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.warning_flags\\.none", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerWarningFlagsDefaultMissing checks for missing compiler.warning_flags.default property in platform.txt.
func PlatformTxtCompilerWarningFlagsDefaultMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.warning_flags\\.default", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerWarningFlagsMoreMissing checks for missing compiler.warning_flags.more property in platform.txt.
func PlatformTxtCompilerWarningFlagsMoreMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.warning_flags\\.more", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerWarningFlagsAllMissing checks for missing compiler.warning_flags.all property in platform.txt.
func PlatformTxtCompilerWarningFlagsAllMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.warning_flags\\.all", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerOptimizationFlagsDebugMissing checks for missing compiler.optimization_flags.debug property in platform.txt.
func PlatformTxtCompilerOptimizationFlagsDebugMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.optimization_flags.release") {
		return ruleresult.Skip, "Dependent property not present"
	}

	if schema.PropertyDependenciesMissing("compiler\\.optimization_flags\\.release", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerOptimizationFlagsReleaseMissing checks for missing compiler.optimization_flags.release property in platform.txt.
func PlatformTxtCompilerOptimizationFlagsReleaseMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.optimization_flags.debug") {
		return ruleresult.Skip, "Dependent property not present"
	}

	if schema.PropertyDependenciesMissing("compiler\\.optimization_flags\\.debug", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCExtraFlagsMissing checks for missing compiler.c.extra_flags property in platform.txt.
func PlatformTxtCompilerCExtraFlagsMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.c\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCExtraFlagsNotEmpty checks for non-empty compiler.c.extra_flags property in platform.txt.
func PlatformTxtCompilerCExtraFlagsNotEmpty() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.c.extra_flags") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyEnumMismatch("compiler\\.c\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCppExtraFlagsMissing checks for missing compiler.cpp.extra_flags property in platform.txt.
func PlatformTxtCompilerCppExtraFlagsMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.cpp\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCppExtraFlagsNotEmpty checks for non-empty compiler.cpp.extra_flags property in platform.txt.
func PlatformTxtCompilerCppExtraFlagsNotEmpty() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.cpp.extra_flags") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyEnumMismatch("compiler\\.cpp\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerSExtraFlagsMissing checks for missing compiler.S.extra_flags property in platform.txt.
func PlatformTxtCompilerSExtraFlagsMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.S\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerSExtraFlagsNotEmpty checks for non-empty compiler.S.extra_flags property in platform.txt.
func PlatformTxtCompilerSExtraFlagsNotEmpty() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.S.extra_flags") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyEnumMismatch("compiler\\.S\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerArExtraFlagsMissing checks for missing compiler.ar.extra_flags property in platform.txt.
func PlatformTxtCompilerArExtraFlagsMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.ar\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerArExtraFlagsNotEmpty checks for non-empty compiler.ar.extra_flags property in platform.txt.
func PlatformTxtCompilerArExtraFlagsNotEmpty() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.ar.extra_flags") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyEnumMismatch("compiler\\.ar\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCElfExtraFlagsMissing checks for missing compiler.c.elf.extra_flags property in platform.txt.
func PlatformTxtCompilerCElfExtraFlagsMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("compiler\\.c\\.elf\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtCompilerCExtraFlagsNotEmpty checks for non-empty compiler.c.extra_flags property in platform.txt.
func PlatformTxtCompilerCElfExtraFlagsNotEmpty() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("compiler.c.elf.extra_flags") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyEnumMismatch("compiler\\.c\\.elf\\.extra_flags", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipePreprocMacrosLTMinLength checks if the platform.txt recipe.preproc.macros property value is less than the minimum length.
func PlatformTxtRecipePreprocMacrosLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.preproc.macros") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.preproc\\.macros", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipePreprocMacrosExtraFlagsSupport checks if platform.txt recipe.preproc.macros provides support for user extra flags.
func PlatformTxtRecipePreprocMacrosExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.preproc.macros") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.preproc\\.macros", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCOPatternMissing checks for missing recipe.c.o.pattern property in platform.txt.
func PlatformTxtRecipeCOPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.c\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCOPatternLTMinLength checks if the platform.txt recipe.c.o.pattern property value is less than the minimum length.
func PlatformTxtRecipeCOPatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.c.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.c\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCOPatternExtraFlagsSupport checks if platform.txt recipe.c.o.pattern provides support for user extra flags.
func PlatformTxtRecipeCOPatternExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.c.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.c\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCppOPatternMissing checks for missing recipe.cpp.o.pattern property in platform.txt.
func PlatformTxtRecipeCppOPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.cpp\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCppOPatternLTMinLength checks if the platform.txt recipe.cpp.o.pattern property value is less than the minimum length.
func PlatformTxtRecipeCppOPatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.cpp.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.cpp\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCppOPatternExtraFlagsSupport checks if platform.txt recipe.cpp.o.pattern provides support for user extra flags.
func PlatformTxtRecipeCppOPatternExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.cpp.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.cpp\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSOPatternMissing checks for missing recipe.S.o.pattern property in platform.txt.
func PlatformTxtRecipeSOPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.S\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSOPatternLTMinLength checks if the platform.txt recipe.S.o.pattern property value is less than the minimum length.
func PlatformTxtRecipeSOPatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.S.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.S\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSOPatternExtraFlagsSupport checks if platform.txt recipe.S.o.pattern provides support for user extra flags.
func PlatformTxtRecipeSOPatternExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.S.o.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.S\\.o\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeArPatternMissing checks for missing recipe.ar.o.pattern property in platform.txt.
func PlatformTxtRecipeArPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.ar\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeArPatternLTMinLength checks if the platform.txt recipe.ar.pattern property value is less than the minimum length.
func PlatformTxtRecipeArPatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.ar.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.ar\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeArPatternExtraFlagsSupport checks if platform.txt recipe.ar.o.pattern provides support for user extra flags.
func PlatformTxtRecipeArPatternExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.ar.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.ar\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCCombinePatternMissing checks for missing recipe.c.combine.pattern property in platform.txt.
func PlatformTxtRecipeCCombinePatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.c\\.combine\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCCombinePatternLTMinLength checks if the platform.txt recipe.c.combine.pattern property value is less than the minimum length.
func PlatformTxtRecipeCCombinePatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.c.combine.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.c\\.combine\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeCCombinePatternExtraFlagsSupport checks if platform.txt recipe.c.combine.pattern provides support for user extra flags.
func PlatformTxtRecipeCCombinePatternExtraFlagsSupport() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.c.combine.pattern") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyPatternMismatch("recipe\\.c\\.combine\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeOutputTmpFileMissing checks for missing recipe.output.tmp_file property in platform.txt.
func PlatformTxtRecipeOutputTmpFileMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.output\\.tmp_file", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeOutputTmpFileLTMinLength checks if the platform.txt recipe.output.tmp_file property value is less than the minimum length.
func PlatformTxtRecipeOutputTmpFileLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.output.tmp_file") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.output\\.tmp_file", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeOutputSaveFileMissing checks for missing recipe.output.save_file property in platform.txt.
func PlatformTxtRecipeOutputSaveFileMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.output\\.save_file", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeOutputSaveFileLTMinLength checks if the platform.txt recipe.output.save_file property value is less than the minimum length.
func PlatformTxtRecipeOutputSaveFileLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.output.save_file") {
		return ruleresult.NotRun, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.output\\.save_file", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSizePatternMissing checks for missing recipe.size.pattern property in platform.txt.
func PlatformTxtRecipeSizePatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.size\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSizePatternLTMinLength checks if the platform.txt recipe.size.pattern property value is less than the minimum length.
func PlatformTxtRecipeSizePatternLTMinLength() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if !projectdata.PlatformTxt().ContainsKey("recipe.size.pattern") {
		return ruleresult.Skip, "Property not present"
	}

	if schema.PropertyLessThanMinLength("recipe\\.size\\.pattern", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSizeRegexMissing checks for missing recipe.size.regex property in platform.txt.
func PlatformTxtRecipeSizeRegexMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.size\\.regex", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtRecipeSizeRegexDataMissing checks for missing recipe.size.regex.data property in platform.txt.
func PlatformTxtRecipeSizeRegexDataMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if schema.RequiredPropertyMissing("recipe\\.size\\.regex\\.data", projectdata.PlatformTxtSchemaValidationResult()[compliancelevel.Strict]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PlatformTxtUploadParamsVerboseMissing checks if any of the tools are missing upload.params.verbose properties.
func PlatformTxtUploadParamsVerboseMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("upload/params\\.verbose", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtUploadParamsQuietMissing checks if any of the programmers are missing upload.params.quiet properties.
func PlatformTxtUploadParamsQuietMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("upload/params\\.quiet", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtUploadPatternMissing checks if any of the programmers are missing upload.pattern properties.
func PlatformTxtUploadPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("upload/pattern", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtProgramParamsVerboseMissing checks if any of the tools are missing program.params.verbose properties.
func PlatformTxtProgramParamsVerboseMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("program/params\\.verbose", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtProgramParamsQuietMissing checks if any of the programmers are missing program.params.quiet properties.
func PlatformTxtProgramParamsQuietMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("program/params\\.quiet", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtProgramPatternMissing checks if any of the programmers are missing program.pattern properties.
func PlatformTxtProgramPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("program/pattern", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtEraseParamsVerboseMissing checks if any of the tools are missing erase.params.verbos properties.
func PlatformTxtEraseParamsVerboseMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("erase/params\\.verbose", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtEraseParamsQuietMissing checks if any of the programmers are missing erase.params.quiet properties.
func PlatformTxtEraseParamsQuietMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("erase/params\\.quiet", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtErasePatternMissing checks if any of the programmers are missing erase.pattern properties.
func PlatformTxtErasePatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("erase/pattern", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtBootloaderParamsVerboseMissing checks if any of the tools are missing bootloader.params.verbos properties.
func PlatformTxtBootloaderParamsVerboseMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("bootloader/params\\.verbose", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtBootloaderParamsQuietMissing checks if any of the programmers are missing bootloader.params.quiet properties.
func PlatformTxtBootloaderParamsQuietMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("bootloader/params\\.quiet", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// PlatformTxtBootloaderPatternMissing checks if any of the programmers are missing bootloader.pattern properties.
func PlatformTxtBootloaderPatternMissing() (result ruleresult.Type, output string) {
	if !projectdata.PlatformTxtExists() {
		return ruleresult.Skip, "Platform has no platform.txt"
	}

	if projectdata.PlatformTxtLoadError() != nil {
		return ruleresult.NotRun, "Couldn't load platform.txt"
	}

	if len(projectdata.PlatformTxtToolNames()) == 0 {
		return ruleresult.Skip, "platform.txt has no tools"
	}

	nonCompliantTools := toolNameMissingRequiredProperty("bootloader/pattern", compliancelevel.Specification)

	if len(nonCompliantTools) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantTools, ", ")
	}

	return ruleresult.Pass, ""
}

// boardIDMissingRequiredProperty returns the list of board IDs missing the given required property.
func boardIDMissingRequiredProperty(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDMissingRequiredProperty(projectdata.BoardsTxtBoardIds(), propertyNameQuery, projectdata.BoardsTxtSchemaValidationResult()[complianceLevel])
}

// boardIDValueLTMinLength returns the list of board IDs with value of the given property less than the minimum length.
func boardIDValueLTMinLength(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValueLTMinLength(projectdata.BoardsTxtBoardIds(), propertyNameQuery, projectdata.BoardsTxtSchemaValidationResult()[complianceLevel])
}

// boardIDValueEnumMismatch returns the list of board IDs with value of the given property not matching the JSON schema enum.
func boardIDValueEnumMismatch(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValueEnumMismatch(projectdata.BoardsTxtBoardIds(), propertyNameQuery, projectdata.BoardsTxtSchemaValidationResult()[complianceLevel])
}

// boardIDValueEnumMismatch returns the list of board IDs with value of the given property not matching the JSON schema pattern.
func boardIDValuePatternMismatch(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValuePatternMismatch(projectdata.BoardsTxtBoardIds(), propertyNameQuery, projectdata.BoardsTxtSchemaValidationResult()[complianceLevel])
}

// programmerIDMissingRequiredProperty returns the list of programmer IDs missing the given required property.
func programmerIDMissingRequiredProperty(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDMissingRequiredProperty(projectdata.ProgrammersTxtProgrammerIds(), propertyNameQuery, projectdata.ProgrammersTxtSchemaValidationResult()[complianceLevel])
}

// programmerIDValueLTMinLength returns the list of programmer IDs with value of the given property less than the minimum length.
func programmerIDValueLTMinLength(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValueLTMinLength(projectdata.ProgrammersTxtProgrammerIds(), propertyNameQuery, projectdata.ProgrammersTxtSchemaValidationResult()[complianceLevel])
}

// programmerIDValueEnumMismatch returns the list of programmer IDs with value of the given property not matching the JSON schema enum.
func programmerIDValueEnumMismatch(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValueEnumMismatch(projectdata.ProgrammersTxtProgrammerIds(), propertyNameQuery, projectdata.ProgrammersTxtSchemaValidationResult()[complianceLevel])
}

// programmerIDValueEnumMismatch returns the list of programmer IDs with value of the given property not matching the JSON schema pattern.
func programmerIDValuePatternMismatch(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	return iDValuePatternMismatch(projectdata.ProgrammersTxtProgrammerIds(), propertyNameQuery, projectdata.ProgrammersTxtSchemaValidationResult()[complianceLevel])
}

// toolNameMissingRequiredProperty returns the list of tool names missing the given required property.
func toolNameMissingRequiredProperty(propertyNameQuery string, complianceLevel compliancelevel.Type) []string {
	nonCompliantTools := []string{}
	for _, tool := range projectdata.PlatformTxtToolNames() {
		if schema.RequiredPropertyMissing("tools/"+tool+"/"+propertyNameQuery, projectdata.PlatformTxtSchemaValidationResult()[complianceLevel]) {
			nonCompliantTools = append(nonCompliantTools, tool)
		}
	}

	return nonCompliantTools
}

// iDMissingRequiredProperty returns the list of first level keys missing the given required property.
func iDMissingRequiredProperty(iDs []string, propertyNameQuery string, validationResult schema.ValidationResult) []string {
	nonCompliantIDs := []string{}
	for _, iD := range iDs {
		if schema.RequiredPropertyMissing(iD+"/"+propertyNameQuery, validationResult) {
			nonCompliantIDs = append(nonCompliantIDs, iD)
		}
	}

	return nonCompliantIDs
}

// iDValueLTMinLength returns the list of first level keys with value of the given property less than the minimum length.
func iDValueLTMinLength(iDs []string, propertyNameQuery string, validationResult schema.ValidationResult) []string {
	nonCompliantIDs := []string{}
	for _, iD := range iDs {
		if schema.PropertyLessThanMinLength(iD+"/"+propertyNameQuery, validationResult) {
			nonCompliantIDs = append(nonCompliantIDs, iD)
		}
	}

	return nonCompliantIDs
}

// iDValueEnumMismatch returns the list of first level keys with value of the given property not matching the JSON schema enum.
func iDValueEnumMismatch(iDs []string, propertyNameQuery string, validationResult schema.ValidationResult) []string {
	nonCompliantIDs := []string{}
	for _, iD := range iDs {
		if schema.PropertyEnumMismatch(iD+"/"+propertyNameQuery, validationResult) {
			nonCompliantIDs = append(nonCompliantIDs, iD)
		}
	}

	return nonCompliantIDs
}

// iDValueEnumMismatch returns the list of first level keys with value of the given property not matching the JSON schema pattern.
func iDValuePatternMismatch(iDs []string, propertyNameQuery string, validationResult schema.ValidationResult) []string {
	nonCompliantIDs := []string{}
	for _, iD := range iDs {
		if schema.PropertyPatternMismatch(iD+"/"+propertyNameQuery, validationResult) {
			nonCompliantIDs = append(nonCompliantIDs, iD)
		}
	}

	return nonCompliantIDs
}
