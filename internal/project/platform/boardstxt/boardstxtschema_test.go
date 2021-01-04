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

// This file contains tests for the boards.txt JSON schema.
package boardstxt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/platform/boardstxt"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validBoardsTxtRaw = []byte(`
	menu.cpu=Processor
	nano.name=Arduino Nano
	nano.upload.tool=avrdude
	nano.upload.protocol=arduino
	nano.upload.maximum_size=123
	nano.upload.maximum_data_size=123
	nano.build.board=AVR_NANO
	nano.build.core=arduino
	nano.menu.cpu.atmega328=ATmega328P
`)

func TestSchemaValid(t *testing.T) {
	validBoardsTxtProperties, err := properties.LoadFromBytes(validBoardsTxtRaw)
	require.Nil(t, err)

	validationResult := boardstxt.Validate(validBoardsTxtProperties)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)
}

func TestMinLength(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		minLength                   int
		complianceLevel             compliancelevel.Type
	}{
		{"menu.foo", "menu/foo", 1, compliancelevel.Permissive},
		{"menu.foo", "menu/foo", 1, compliancelevel.Specification},
		{"menu.foo", "menu/foo", 1, compliancelevel.Strict},

		{"foo.name", "foo/name", 1, compliancelevel.Permissive},
		{"foo.name", "foo/name", 1, compliancelevel.Specification},
		{"foo.name", "foo/name", 1, compliancelevel.Strict},

		{"foo.build.board", "foo/build\\.board", 1, compliancelevel.Permissive},
		{"foo.build.board", "foo/build\\.board", 1, compliancelevel.Specification},
		{"foo.build.board", "foo/build\\.board", 1, compliancelevel.Strict},

		{"foo.build.core", "foo/build\\.core", 1, compliancelevel.Permissive},
		{"foo.build.core", "foo/build\\.core", 1, compliancelevel.Specification},
		{"foo.build.core", "foo/build\\.core", 1, compliancelevel.Strict},

		{"foo.menu.bar.baz", "foo/menu\\.bar\\.baz", 1, compliancelevel.Permissive},
		{"foo.menu.bar.baz", "foo/menu\\.bar\\.baz", 1, compliancelevel.Specification},
		{"foo.menu.bar.baz", "foo/menu\\.bar\\.baz", 1, compliancelevel.Strict},

		{"foo.upload.tool", "foo/upload\\.tool", 1, compliancelevel.Permissive},
		{"foo.upload.tool", "foo/upload\\.tool", 1, compliancelevel.Specification},
		{"foo.upload.tool", "foo/upload\\.tool", 1, compliancelevel.Strict},
	}

	// Test schema validation results with value length < minimum.
	for _, testTable := range testTables {
		boardsTxt, err := properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength-1))

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyLessThanMinLength(testTable.propertyName, boardstxt.Validate(boardsTxt)[testTable.complianceLevel]))
		})

		// Test schema validation results with minimum value length.
		boardsTxt, err = properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength))

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(testTable.validationErrorPropertyName, boardstxt.Validate(boardsTxt)[testTable.complianceLevel]))
		})
	}
}

func TestEmpty(t *testing.T) {
	// None of the root properties are required, so an empty boards.txt is valid.
	validBoardsTxtProperties, err := properties.LoadFromBytes([]byte{})
	require.Nil(t, err)

	validationResult := boardstxt.Validate(validBoardsTxtProperties)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)
}

func TestRequired(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		complianceLevel             compliancelevel.Type
		assertion                   assert.BoolAssertionFunc
	}{
		{"menu.cpu", "menu", compliancelevel.Permissive, assert.False},
		{"menu.cpu", "menu", compliancelevel.Specification, assert.False},
		{"menu.cpu", "menu", compliancelevel.Strict, assert.False},

		{"nano.name", "nano/name", compliancelevel.Permissive, assert.True},
		{"nano.name", "nano/name", compliancelevel.Specification, assert.True},
		{"nano.name", "nano/name", compliancelevel.Strict, assert.True},

		{"nano.upload.tool", "nano/upload\\.tool", compliancelevel.Permissive, assert.True},
		{"nano.upload.tool", "nano/upload\\.tool", compliancelevel.Specification, assert.True},
		{"nano.upload.tool", "nano/upload\\.tool", compliancelevel.Strict, assert.True},

		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", compliancelevel.Permissive, assert.False},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", compliancelevel.Specification, assert.False},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", compliancelevel.Strict, assert.True},

		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", compliancelevel.Permissive, assert.False},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", compliancelevel.Specification, assert.False},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", compliancelevel.Strict, assert.True},

		{"nano.upload.protocol", "nano/upload\\.protocol", compliancelevel.Permissive, assert.False},
		{"nano.upload.protocol", "nano/upload\\.protocol", compliancelevel.Specification, assert.False},
		{"nano.upload.protocol", "nano/upload\\.protocol", compliancelevel.Strict, assert.False},

		{"nano.build.board", "nano/build\\.board", compliancelevel.Permissive, assert.False},
		{"nano.build.board", "nano/build\\.board", compliancelevel.Specification, assert.False},
		{"nano.build.board", "nano/build\\.board", compliancelevel.Strict, assert.True},

		{"nano.build.core", "nano/build\\.core", compliancelevel.Permissive, assert.True},
		{"nano.build.core", "nano/build\\.core", compliancelevel.Specification, assert.True},
		{"nano.build.core", "nano/build\\.core", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		boardsTxt, err := properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Remove(testTable.propertyName)

		validationResult := boardstxt.Validate(boardsTxt)
		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.RequiredPropertyMissing(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestEnum(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		propertyValue               string
		complianceLevel             compliancelevel.Type
		assertion                   assert.BoolAssertionFunc
	}{
		{"nano.hide", "nano/hide", "true", compliancelevel.Permissive, assert.False},
		{"nano.hide", "nano/hide", "true", compliancelevel.Specification, assert.True},
		{"nano.hide", "nano/hide", "true", compliancelevel.Strict, assert.True},
		{"nano.hide", "nano/hide", "false", compliancelevel.Permissive, assert.False},
		{"nano.hide", "nano/hide", "false", compliancelevel.Specification, assert.True},
		{"nano.hide", "nano/hide", "false", compliancelevel.Strict, assert.True},
		{"nano.hide", "nano/hide", "", compliancelevel.Permissive, assert.False},
		{"nano.hide", "nano/hide", "", compliancelevel.Specification, assert.False},
		{"nano.hide", "nano/hide", "", compliancelevel.Strict, assert.False},

		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "true", compliancelevel.Permissive, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "true", compliancelevel.Specification, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "true", compliancelevel.Strict, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "false", compliancelevel.Permissive, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "false", compliancelevel.Specification, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "false", compliancelevel.Strict, assert.False},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "foo", compliancelevel.Permissive, assert.True},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "foo", compliancelevel.Specification, assert.True},
		{"nano.serial.disableDTR", "nano/serial\\.disableDTR", "foo", compliancelevel.Strict, assert.True},

		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "true", compliancelevel.Permissive, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "true", compliancelevel.Specification, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "true", compliancelevel.Strict, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "false", compliancelevel.Permissive, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "false", compliancelevel.Specification, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "false", compliancelevel.Strict, assert.False},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "foo", compliancelevel.Permissive, assert.True},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "foo", compliancelevel.Specification, assert.True},
		{"nano.serial.disableRTS", "nano/serial\\.disableRTS", "foo", compliancelevel.Strict, assert.True},

		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "true", compliancelevel.Permissive, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "true", compliancelevel.Specification, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "true", compliancelevel.Strict, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "false", compliancelevel.Permissive, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "false", compliancelevel.Specification, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "false", compliancelevel.Strict, assert.False},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "foo", compliancelevel.Permissive, assert.True},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "foo", compliancelevel.Specification, assert.True},
		{"nano.upload.use_1200bps_touch", "nano/upload\\.use_1200bps_touch", "foo", compliancelevel.Strict, assert.True},

		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "true", compliancelevel.Permissive, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "true", compliancelevel.Specification, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "true", compliancelevel.Strict, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "false", compliancelevel.Permissive, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "false", compliancelevel.Specification, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "false", compliancelevel.Strict, assert.False},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "foo", compliancelevel.Permissive, assert.True},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "foo", compliancelevel.Specification, assert.True},
		{"nano.upload.wait_for_upload_port", "nano/upload\\.wait_for_upload_port", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		boardsTxt, err := properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Set(testTable.propertyName, testTable.propertyValue)

		validationResult := boardstxt.Validate(boardsTxt)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyName, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyEnumMismatch(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestPattern(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		propertyValue               string
		complianceLevel             compliancelevel.Type
		assertion                   assert.BoolAssertionFunc
	}{
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "123", compliancelevel.Permissive, assert.False},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "123", compliancelevel.Specification, assert.False},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "123", compliancelevel.Strict, assert.False},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "foo", compliancelevel.Permissive, assert.True},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "foo", compliancelevel.Specification, assert.True},
		{"nano.upload.maximum_size", "nano/upload\\.maximum_size", "foo", compliancelevel.Strict, assert.True},

		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "123", compliancelevel.Permissive, assert.False},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "123", compliancelevel.Specification, assert.False},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "123", compliancelevel.Strict, assert.False},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "foo", compliancelevel.Permissive, assert.True},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "foo", compliancelevel.Specification, assert.True},
		{"nano.upload.maximum_data_size", "nano/upload\\.maximum_data_size", "foo", compliancelevel.Strict, assert.True},

		{"nano.vid.0", "nano/vid\\.0", "0xABCD", compliancelevel.Permissive, assert.False},
		{"nano.vid.0", "nano/vid\\.0", "0xABCD", compliancelevel.Specification, assert.False},
		{"nano.vid.0", "nano/vid\\.0", "0xABCD", compliancelevel.Strict, assert.False},
		{"nano.vid.0", "nano/vid\\.0", "foo", compliancelevel.Permissive, assert.True},
		{"nano.vid.0", "nano/vid\\.0", "foo", compliancelevel.Specification, assert.True},
		{"nano.vid.0", "nano/vid\\.0", "foo", compliancelevel.Strict, assert.True},

		{"nano.pid.0", "nano/pid\\.0", "0xABCD", compliancelevel.Permissive, assert.False},
		{"nano.pid.0", "nano/pid\\.0", "0xABCD", compliancelevel.Specification, assert.False},
		{"nano.pid.0", "nano/pid\\.0", "0xABCD", compliancelevel.Strict, assert.False},
		{"nano.pid.0", "nano/pid\\.0", "foo", compliancelevel.Permissive, assert.True},
		{"nano.pid.0", "nano/pid\\.0", "foo", compliancelevel.Specification, assert.True},
		{"nano.pid.0", "nano/pid\\.0", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		boardsTxt, err := properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Set(testTable.propertyName, testTable.propertyValue)

		validationResult := boardstxt.Validate(boardsTxt)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyName, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyPatternMismatch(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestPropertyNames(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		complianceLevel             compliancelevel.Type
		assertion                   assert.BoolAssertionFunc
	}{
		{"nano.compiler.c.extra_flags", "nano/compiler\\.c\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.c.extra_flags", "nano/compiler\\.c\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.c.extra_flags", "nano/compiler\\.c\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.c.elf.extra_flags", "nano/compiler\\.c\\.elf\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.c.elf.extra_flags", "nano/compiler\\.c\\.elf\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.c.elf.extra_flags", "nano/compiler\\.c\\.elf\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.S.extra_flags", "nano/compiler\\.S\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.S.extra_flags", "nano/compiler\\.S\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.S.extra_flags", "nano/compiler\\.S\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.cpp.extra_flags", "nano/compiler\\.cpp\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.cpp.extra_flags", "nano/compiler\\.cpp\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.cpp.extra_flags", "nano/compiler\\.cpp\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.ar.extra_flags", "nano/compiler\\.ar\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.ar.extra_flags", "nano/compiler\\.ar\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.ar.extra_flags", "nano/compiler\\.ar\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.objcopy.eep.extra_flags", "nano/compiler\\.objcopy\\.eep\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.objcopy.eep.extra_flags", "nano/compiler\\.objcopy\\.eep\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.objcopy.eep.extra_flags", "nano/compiler\\.objcopy\\.eep\\.extra_flags", compliancelevel.Strict, assert.True},

		{"nano.compiler.elf2hex.extra_flags", "nano/compiler\\.elf2hex\\.extra_flags", compliancelevel.Permissive, assert.False},
		{"nano.compiler.elf2hex.extra_flags", "nano/compiler\\.elf2hex\\.extra_flags", compliancelevel.Specification, assert.False},
		{"nano.compiler.elf2hex.extra_flags", "nano/compiler\\.elf2hex\\.extra_flags", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		boardsTxt, err := properties.LoadFromBytes(validBoardsTxtRaw)
		require.Nil(t, err)
		boardsTxt.Set(testTable.propertyName, "foo")

		validationResult := boardstxt.Validate(boardsTxt)

		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.ValidationErrorMatch("#/"+testTable.validationErrorPropertyName, "/userExtraFlagsProperties/", "", "", validationResult[testTable.complianceLevel]))
		})
	}
}
