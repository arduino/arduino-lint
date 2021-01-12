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

// This file contains tests for the platform.txt JSON schema.
package platformtxt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/platform/platformtxt"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validPlatformTxtRaw = []byte(`
	name=Arduino AVR Boards
	version=1.8.3
	compiler.warning_flags.none=asdf
	compiler.warning_flags.default=asdf
	compiler.warning_flags.more=asdf
	compiler.warning_flags.all=asdf
	compiler.c.extra_flags=
	compiler.c.elf.extra_flags=
	compiler.S.extra_flags=
	compiler.cpp.extra_flags=
	compiler.ar.extra_flags=
	compiler.objcopy.eep.extra_flags=
	compiler.elf2hex.extra_flags=
	recipe.c.o.pattern=asdf {compiler.c.extra_flags}
	recipe.cpp.o.pattern=asdf {compiler.cpp.extra_flags}
	recipe.S.o.pattern=asdf {compiler.S.extra_flags}
	recipe.ar.pattern=asdf {compiler.ar.extra_flags}
	recipe.c.combine.pattern=asdf {compiler.c.elf.extra_flags}
	recipe.objcopy.eep.pattern=asdf
	recipe.objcopy.hex.pattern=asdf
	recipe.output.tmp_file=asdf
	recipe.output.save_file=asdf
	recipe.size.pattern=asdf
	recipe.size.regex=asdf
	recipe.size.regex.data=asdf
	tools.avrdude.upload.params.verbose=-v
	tools.avrdude.upload.params.quiet=-q -q
	tools.avrdude.upload.pattern=asdf
	tools.avrdude.program.params.verbose=-v
	tools.avrdude.program.params.quiet=-q -q
	tools.avrdude.program.pattern=asdf
	tools.bossac.upload.params.verbose=-v
	tools.bossac.upload.params.quiet=-q -q
	tools.bossac.upload.pattern=asdf
`)

func TestSchemaValid(t *testing.T) {
	validPlatformTxtProperties, err := properties.LoadFromBytes(validPlatformTxtRaw)
	require.Nil(t, err)

	validationResult := platformtxt.Validate(validPlatformTxtProperties)

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
		{"name", "name", 1, compliancelevel.Permissive},
		{"name", "name", 1, compliancelevel.Specification},
		{"name", "name", 1, compliancelevel.Strict},

		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", 1, compliancelevel.Permissive},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", 1, compliancelevel.Specification},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", 1, compliancelevel.Strict},

		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", 1, compliancelevel.Permissive},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", 1, compliancelevel.Specification},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", 1, compliancelevel.Strict},

		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", 1, compliancelevel.Permissive},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", 1, compliancelevel.Specification},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", 1, compliancelevel.Strict},

		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", 1, compliancelevel.Permissive},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", 1, compliancelevel.Specification},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", 1, compliancelevel.Strict},

		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", 1, compliancelevel.Permissive},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", 1, compliancelevel.Specification},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", 1, compliancelevel.Strict},

		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", 1, compliancelevel.Permissive},
		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", 1, compliancelevel.Specification},
		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", 1, compliancelevel.Strict},

		{"recipe.output.save_file", "recipe\\.output\\.save_file", 1, compliancelevel.Permissive},
		{"recipe.output.save_file", "recipe\\.output\\.save_file", 1, compliancelevel.Specification},
		{"recipe.output.save_file", "recipe\\.output\\.save_file", 1, compliancelevel.Strict},

		{"recipe.size.pattern", "recipe\\.size\\.pattern", 1, compliancelevel.Strict},

		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", 1, compliancelevel.Permissive},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", 1, compliancelevel.Specification},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", 1, compliancelevel.Strict},
	}

	// Test schema validation results with value length < minimum.
	for _, testTable := range testTables {
		platformTxt, err := properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength-1))

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyLessThanMinLength(testTable.propertyName, platformtxt.Validate(platformTxt)[testTable.complianceLevel]))
		})

		// Test schema validation results with minimum value length.
		platformTxt, err = properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength))

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(testTable.validationErrorPropertyName, platformtxt.Validate(platformTxt)[testTable.complianceLevel]))
		})
	}
}

func TestRequired(t *testing.T) {
	testTables := []struct {
		propertyName                string
		validationErrorPropertyName string
		complianceLevel             compliancelevel.Type
		assertion                   assert.BoolAssertionFunc
	}{
		{"name", "name", compliancelevel.Permissive, assert.True},
		{"name", "name", compliancelevel.Specification, assert.True},
		{"name", "name", compliancelevel.Strict, assert.True},

		{"version", "version", compliancelevel.Permissive, assert.True},
		{"version", "version", compliancelevel.Specification, assert.True},
		{"version", "version", compliancelevel.Strict, assert.True},

		{"compiler.warning_flags.none", "compiler\\.warning_flags\\.none", compliancelevel.Permissive, assert.False},
		{"compiler.warning_flags.none", "compiler\\.warning_flags\\.none", compliancelevel.Specification, assert.False},
		{"compiler.warning_flags.none", "compiler\\.warning_flags\\.none", compliancelevel.Strict, assert.True},

		{"compiler.warning_flags.default", "compiler\\.warning_flags\\.default", compliancelevel.Permissive, assert.False},
		{"compiler.warning_flags.default", "compiler\\.warning_flags\\.default", compliancelevel.Specification, assert.False},
		{"compiler.warning_flags.default", "compiler\\.warning_flags\\.default", compliancelevel.Strict, assert.True},

		{"compiler.warning_flags.more", "compiler\\.warning_flags\\.more", compliancelevel.Permissive, assert.False},
		{"compiler.warning_flags.more", "compiler\\.warning_flags\\.more", compliancelevel.Specification, assert.False},
		{"compiler.warning_flags.more", "compiler\\.warning_flags\\.more", compliancelevel.Strict, assert.True},

		{"compiler.warning_flags.all", "compiler\\.warning_flags\\.all", compliancelevel.Permissive, assert.False},
		{"compiler.warning_flags.all", "compiler\\.warning_flags\\.all", compliancelevel.Specification, assert.False},
		{"compiler.warning_flags.all", "compiler\\.warning_flags\\.all", compliancelevel.Strict, assert.True},

		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", compliancelevel.Permissive, assert.True},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", compliancelevel.Specification, assert.True},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", compliancelevel.Strict, assert.True},

		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", compliancelevel.Permissive, assert.True},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", compliancelevel.Specification, assert.True},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", compliancelevel.Strict, assert.True},

		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", compliancelevel.Permissive, assert.True},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", compliancelevel.Specification, assert.True},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", compliancelevel.Strict, assert.True},

		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", compliancelevel.Permissive, assert.True},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", compliancelevel.Specification, assert.True},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", compliancelevel.Strict, assert.True},

		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", compliancelevel.Permissive, assert.True},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", compliancelevel.Specification, assert.True},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", compliancelevel.Strict, assert.True},

		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", compliancelevel.Permissive, assert.True},
		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", compliancelevel.Specification, assert.True},
		{"recipe.output.tmp_file", "recipe\\.output\\.tmp_file", compliancelevel.Strict, assert.True},

		{"tools.avrdude.upload.pattern", "tools/avrdude/upload/pattern", compliancelevel.Permissive, assert.True},
		{"tools.avrdude.upload.pattern", "tools/avrdude/upload/pattern", compliancelevel.Specification, assert.True},
		{"tools.avrdude.upload.pattern", "tools/avrdude/upload/pattern", compliancelevel.Strict, assert.True},

		{"tools.avrdude.program.params.verbose", "tools/avrdude/program/params\\.verbose", compliancelevel.Permissive, assert.True},
		{"tools.avrdude.program.params.verbose", "tools/avrdude/program/params\\.verbose", compliancelevel.Specification, assert.True},
		{"tools.avrdude.program.params.verbose", "tools/avrdude/program/params\\.verbose", compliancelevel.Strict, assert.True},

		{"tools.avrdude.program.params.quiet", "tools/avrdude/program/params\\.quiet", compliancelevel.Permissive, assert.True},
		{"tools.avrdude.program.params.quiet", "tools/avrdude/program/params\\.quiet", compliancelevel.Specification, assert.True},
		{"tools.avrdude.program.params.quiet", "tools/avrdude/program/params\\.quiet", compliancelevel.Strict, assert.True},

		{"tools.avrdude.program.pattern", "tools/avrdude/program/pattern", compliancelevel.Permissive, assert.True},
		{"tools.avrdude.program.pattern", "tools/avrdude/program/pattern", compliancelevel.Specification, assert.True},
		{"tools.avrdude.program.pattern", "tools/avrdude/program/pattern", compliancelevel.Strict, assert.True},

		{"tools.bossac.upload.pattern", "tools/bossac/upload/pattern", compliancelevel.Permissive, assert.True},
		{"tools.bossac.upload.pattern", "tools/bossac/upload/pattern", compliancelevel.Specification, assert.True},
		{"tools.bossac.upload.pattern", "tools/bossac/upload/pattern", compliancelevel.Strict, assert.True},

		{"compiler.c.extra_flags", "compiler.c.extra_flags", compliancelevel.Permissive, assert.False},
		{"compiler.c.extra_flags", "compiler.c.extra_flags", compliancelevel.Specification, assert.False},
		{"compiler.c.extra_flags", "compiler.c.extra_flags", compliancelevel.Strict, assert.True},

		{"compiler.c.elf.extra_flags", "compiler.c.elf.extra_flags", compliancelevel.Permissive, assert.False},
		{"compiler.c.elf.extra_flags", "compiler.c.elf.extra_flags", compliancelevel.Specification, assert.False},
		{"compiler.c.elf.extra_flags", "compiler.c.elf.extra_flags", compliancelevel.Strict, assert.True},

		{"compiler.S.extra_flags", "compiler.S.extra_flags", compliancelevel.Permissive, assert.False},
		{"compiler.S.extra_flags", "compiler.S.extra_flags", compliancelevel.Specification, assert.False},
		{"compiler.S.extra_flags", "compiler.S.extra_flags", compliancelevel.Strict, assert.True},

		{"compiler.cpp.extra_flags", "compiler.cpp.extra_flags", compliancelevel.Permissive, assert.False},
		{"compiler.cpp.extra_flags", "compiler.cpp.extra_flags", compliancelevel.Specification, assert.False},
		{"compiler.cpp.extra_flags", "compiler.cpp.extra_flags", compliancelevel.Strict, assert.True},

		{"compiler.ar.extra_flags", "compiler.ar.extra_flags", compliancelevel.Permissive, assert.False},
		{"compiler.ar.extra_flags", "compiler.ar.extra_flags", compliancelevel.Specification, assert.False},
		{"compiler.ar.extra_flags", "compiler.ar.extra_flags", compliancelevel.Strict, assert.True},

		{"recipe.size.pattern", "recipe.size.pattern", compliancelevel.Permissive, assert.False},
		{"recipe.size.pattern", "recipe.size.pattern", compliancelevel.Specification, assert.False},
		{"recipe.size.pattern", "recipe.size.pattern", compliancelevel.Strict, assert.True},

		{"recipe.size.regex", "recipe.size.regex", compliancelevel.Permissive, assert.False},
		{"recipe.size.regex", "recipe.size.regex", compliancelevel.Specification, assert.False},
		{"recipe.size.regex", "recipe.size.regex", compliancelevel.Strict, assert.True},

		{"recipe.size.regex.data", "recipe.size.regex.data", compliancelevel.Permissive, assert.False},
		{"recipe.size.regex.data", "recipe.size.regex.data", compliancelevel.Specification, assert.False},
		{"recipe.size.regex.data", "recipe.size.regex.data", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		platformTxt, err := properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Remove(testTable.propertyName)

		validationResult := platformtxt.Validate(platformTxt)
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
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "", compliancelevel.Permissive, assert.False},
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "", compliancelevel.Specification, assert.False},
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "", compliancelevel.Strict, assert.False},
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "foo", compliancelevel.Permissive, assert.False},
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "foo", compliancelevel.Specification, assert.False},
		{"compiler.c.extra_flags", "compiler\\.c\\.extra_flags", "foo", compliancelevel.Strict, assert.True},

		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "", compliancelevel.Permissive, assert.False},
		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "", compliancelevel.Specification, assert.False},
		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "", compliancelevel.Strict, assert.False},
		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "foo", compliancelevel.Permissive, assert.False},
		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "foo", compliancelevel.Specification, assert.False},
		{"compiler.c.elf.extra_flags", "compiler\\.c\\.elf\\.extra_flags", "foo", compliancelevel.Strict, assert.True},

		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "", compliancelevel.Permissive, assert.False},
		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "", compliancelevel.Specification, assert.False},
		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "", compliancelevel.Strict, assert.False},
		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "foo", compliancelevel.Permissive, assert.False},
		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "foo", compliancelevel.Specification, assert.False},
		{"compiler.S.extra_flags", "compiler\\.S\\.extra_flags", "foo", compliancelevel.Strict, assert.True},

		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "", compliancelevel.Permissive, assert.False},
		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "", compliancelevel.Specification, assert.False},
		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "", compliancelevel.Strict, assert.False},
		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "foo", compliancelevel.Permissive, assert.False},
		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "foo", compliancelevel.Specification, assert.False},
		{"compiler.cpp.extra_flags", "compiler\\.cpp\\.extra_flags", "foo", compliancelevel.Strict, assert.True},

		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "", compliancelevel.Permissive, assert.False},
		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "", compliancelevel.Specification, assert.False},
		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "", compliancelevel.Strict, assert.False},
		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "foo", compliancelevel.Permissive, assert.False},
		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "foo", compliancelevel.Specification, assert.False},
		{"compiler.ar.extra_flags", "compiler\\.ar\\.extra_flags", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		platformTxt, err := properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Set(testTable.propertyName, testTable.propertyValue)

		validationResult := platformtxt.Validate(platformTxt)

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
		{"version", "version", "1.0.0", compliancelevel.Permissive, assert.False},
		{"version", "version", "1.0.0", compliancelevel.Specification, assert.False},
		{"version", "version", "1.0.0", compliancelevel.Strict, assert.False},
		{"version", "version", "1.0", compliancelevel.Permissive, assert.False},
		{"version", "version", "1.0", compliancelevel.Specification, assert.True},
		{"version", "version", "1.0", compliancelevel.Strict, assert.True},
		{"version", "version", "{foo}", compliancelevel.Permissive, assert.False},
		{"version", "version", "{foo}", compliancelevel.Specification, assert.False},
		{"version", "version", "{foo}", compliancelevel.Strict, assert.False},
		{"version", "version", "foo", compliancelevel.Permissive, assert.True},
		{"version", "version", "foo", compliancelevel.Specification, assert.True},
		{"version", "version", "foo", compliancelevel.Strict, assert.True},

		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo {compiler.c.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo {compiler.c.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo {compiler.c.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo", compliancelevel.Specification, assert.False},
		{"recipe.c.o.pattern", "recipe\\.c\\.o\\.pattern", "foo", compliancelevel.Strict, assert.True},

		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo", compliancelevel.Specification, assert.False},
		{"recipe.cpp.o.pattern", "recipe\\.cpp\\.o\\.pattern", "foo", compliancelevel.Strict, assert.True},

		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo {compiler.S.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo {compiler.S.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo {compiler.S.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo", compliancelevel.Specification, assert.False},
		{"recipe.S.o.pattern", "recipe\\.S\\.o\\.pattern", "foo", compliancelevel.Strict, assert.True},

		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo {compiler.ar.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo {compiler.ar.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo {compiler.ar.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo", compliancelevel.Specification, assert.False},
		{"recipe.ar.pattern", "recipe\\.ar\\.pattern", "foo", compliancelevel.Strict, assert.True},

		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo {compiler.c.elf.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo {compiler.c.elf.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo {compiler.c.elf.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo", compliancelevel.Specification, assert.False},
		{"recipe.c.combine.pattern", "recipe\\.c\\.combine\\.pattern", "foo", compliancelevel.Strict, assert.True},

		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Permissive, assert.False},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Specification, assert.False},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo {compiler.cpp.extra_flags} bar", compliancelevel.Strict, assert.False},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo", compliancelevel.Permissive, assert.False},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo", compliancelevel.Specification, assert.False},
		{"recipe.preproc.macros", "recipe\\.preproc\\.macros", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		platformTxt, err := properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Set(testTable.propertyName, testTable.propertyValue)

		validationResult := platformtxt.Validate(platformTxt)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyName, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyPatternMismatch(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestDependencies(t *testing.T) {
	testTables := []struct {
		dependentPropertyName       string
		dependencyPropertyName      string
		validationErrorPropertyName string
		complianceLevel             compliancelevel.Type
	}{
		{"compiler.optimization_flags.debug", "compiler.optimization_flags.release", "compiler\\.optimization_flags\\.debug", compliancelevel.Permissive},
		{"compiler.optimization_flags.debug", "compiler.optimization_flags.release", "compiler\\.optimization_flags\\.debug", compliancelevel.Specification},
		{"compiler.optimization_flags.debug", "compiler.optimization_flags.release", "compiler\\.optimization_flags\\.debug", compliancelevel.Strict},

		{"compiler.optimization_flags.release", "compiler.optimization_flags.debug", "compiler\\.optimization_flags\\.release", compliancelevel.Permissive},    // This is a bidirectional dependency.
		{"compiler.optimization_flags.release", "compiler.optimization_flags.debug", "compiler\\.optimization_flags\\.release", compliancelevel.Specification}, // This is a bidirectional dependency.
		{"compiler.optimization_flags.release", "compiler.optimization_flags.debug", "compiler\\.optimization_flags\\.release", compliancelevel.Strict},        // This is a bidirectional dependency.
	}

	for _, testTable := range testTables {
		platformTxt, err := properties.LoadFromBytes(validPlatformTxtRaw)
		require.Nil(t, err)
		platformTxt.Set(testTable.dependentPropertyName, "foo")
		platformTxt.Set(testTable.dependencyPropertyName, "foo")

		validationResult := platformtxt.Validate(platformTxt)
		t.Run(fmt.Sprintf("Dependency of %s present (%s)", testTable.dependentPropertyName, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyDependenciesMissing(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})

		platformTxt.Remove(testTable.dependencyPropertyName)

		validationResult = platformtxt.Validate(platformTxt)
		t.Run(fmt.Sprintf("Dependency of %s missing (%s)", testTable.dependentPropertyName, testTable.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyDependenciesMissing(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})

		platformTxt.Remove(testTable.dependentPropertyName)

		validationResult = platformtxt.Validate(platformTxt)
		t.Run(fmt.Sprintf("Dependent %s not present (%s)", testTable.dependentPropertyName, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyDependenciesMissing(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}
