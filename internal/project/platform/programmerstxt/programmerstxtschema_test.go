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

// This file contains tests for the programmers.txt JSON schema.
package programmerstxt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/platform/programmerstxt"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validProgrammersTxtRaw = []byte(`
	usbasp.name=USBasp
	usbasp.protocol=usbasp
	usbasp.program.tool=avrdude
`)

func TestSchemaValid(t *testing.T) {
	validProgrammersTxtProperties, err := properties.LoadFromBytes(validProgrammersTxtRaw)
	require.Nil(t, err)

	validationResult := programmerstxt.Validate(validProgrammersTxtProperties)

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
		{"foo.name", "foo/name", 1, compliancelevel.Permissive},
		{"foo.name", "foo/name", 1, compliancelevel.Specification},
		{"foo.name", "foo/name", 1, compliancelevel.Strict},

		{"foo.program.tool", "foo/program\\.tool", 1, compliancelevel.Permissive},
		{"foo.program.tool", "foo/program\\.tool", 1, compliancelevel.Specification},
		{"foo.program.tool", "foo/program\\.tool", 1, compliancelevel.Strict},
	}

	// Test schema validation results with value length < minimum.
	for _, testTable := range testTables {
		programmersTxt, err := properties.LoadFromBytes(validProgrammersTxtRaw)
		require.Nil(t, err)
		programmersTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength-1))

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyLessThanMinLength(testTable.propertyName, programmerstxt.Validate(programmersTxt)[testTable.complianceLevel]))
		})

		// Test schema validation results with minimum value length.
		programmersTxt, err = properties.LoadFromBytes(validProgrammersTxtRaw)
		require.Nil(t, err)
		programmersTxt.Set(testTable.propertyName, strings.Repeat("a", testTable.minLength))

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", testTable.propertyName, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(testTable.validationErrorPropertyName, programmerstxt.Validate(programmersTxt)[testTable.complianceLevel]))
		})
	}
}

func TestEmpty(t *testing.T) {
	// None of the root properties are required, so an empty programmers.txt is valid.
	programmersTxt, err := properties.LoadFromBytes([]byte{})
	require.Nil(t, err)

	validationResult := programmerstxt.Validate(programmersTxt)

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
		{"usbasp.name", "usbasp/name", compliancelevel.Permissive, assert.True},
		{"usbasp.name", "usbasp/name", compliancelevel.Specification, assert.True},
		{"usbasp.name", "usbasp/name", compliancelevel.Strict, assert.True},

		{"usbasp.program.tool", "usbasp/program\\.tool", compliancelevel.Permissive, assert.True},
		{"usbasp.program.tool", "usbasp/program\\.tool", compliancelevel.Specification, assert.True},
		{"usbasp.program.tool", "usbasp/program\\.tool", compliancelevel.Strict, assert.True},

		{"usbasp.foo.bar", "usbasp/foo\\.bar", compliancelevel.Permissive, assert.False},
		{"usbasp.foo.bar", "usbasp/foo\\.bar", compliancelevel.Specification, assert.False},
		{"usbasp.foo.bar", "usbasp/foo\\.bar", compliancelevel.Strict, assert.False},
	}

	for _, testTable := range testTables {
		programmersTxt, err := properties.LoadFromBytes(validProgrammersTxtRaw)
		require.Nil(t, err)
		programmersTxt.Remove(testTable.propertyName)

		validationResult := programmerstxt.Validate(programmersTxt)
		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyName, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.RequiredPropertyMissing(testTable.validationErrorPropertyName, validationResult[testTable.complianceLevel]))
		})
	}
}
