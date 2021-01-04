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
