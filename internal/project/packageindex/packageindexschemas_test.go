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

// This file contains tests for the package index JSON schemas.
package packageindex_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/packageindex"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonpointer"
)

var validIndexRaw = []byte(`
{
	"packages": [
		{
			"name": "notarduino",
			"maintainer": "NotArduino",
			"websiteURL": "http://www.arduino.cc/",
			"email": "packages@arduino.cc",
			"help": {
				"online": "http://www.arduino.cc/en/Reference/HomePage"
			},
			"platforms": [
				{
					"name": "Arduino AVR Boards",
					"architecture": "avr",
					"version": "1.8.3",
					"category": "Contributed",
					"help": {
						"online": "http://www.arduino.cc/en/Reference/HomePage"
					},
					"url": "http://downloads.arduino.cc/cores/avr-1.8.3.tar.bz2",
					"archiveFileName": "avr-1.8.3.tar.bz2",
					"checksum": "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14",
					"size": "4941548",
					"boards": [
						{"name": "Arduino Uno"}
					],
					"toolsDependencies": [
						{
							"packager": "arduino",
							"name": "avr-gcc",
							"version": "7.3.0-atmel3.6.1-arduino7"
						}
					]
				}
			],
			"tools": [
				{
					"name": "avr-gcc",
					"version": "7.3.0-atmel3.6.1-arduino7",
					"systems": [
						{
							"size": "34683056",
							"checksum": "SHA-256:3903553d035da59e33cff9941b857c3cb379cb0638105dfdf69c97f0acc8e7b5",
							"host": "arm-linux-gnueabihf",
							"archiveFileName": "avr-gcc-7.3.0-atmel3.6.1-arduino7-arm-linux-gnueabihf.tar.bz2",
							"url": "http://downloads.arduino.cc/tools/avr-gcc-7.3.0-atmel3.6.1-arduino7-arm-linux-gnueabihf.tar.bz2"
						}
					]
				}
			]
		}
	]
}
`)

func TestSchemaValid(t *testing.T) {
	var validIndex map[string]interface{}
	err := json.Unmarshal(validIndexRaw, &validIndex)
	require.NoError(t, err)

	validationResult := packageindex.Validate(validIndex)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)
}

func TestMinLength(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		minLength             int
		complianceLevel       compliancelevel.Type
	}{
		{"/packages/0/name", 1, compliancelevel.Permissive},
		{"/packages/0/name", 1, compliancelevel.Specification},
		{"/packages/0/name", 1, compliancelevel.Strict},

		{"/packages/0/maintainer", 1, compliancelevel.Permissive},
		{"/packages/0/maintainer", 1, compliancelevel.Specification},
		{"/packages/0/maintainer", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/name", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/name", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/name", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/architecture", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/architecture", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/architecture", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/archiveFileName", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/archiveFileName", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/archiveFileName", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/boards/0/name", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/boards/0/name", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/boards/0/name", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/toolsDependencies/0/packager", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/toolsDependencies/0/packager", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/toolsDependencies/0/packager", 1, compliancelevel.Strict},

		{"/packages/0/platforms/0/toolsDependencies/0/name", 1, compliancelevel.Permissive},
		{"/packages/0/platforms/0/toolsDependencies/0/name", 1, compliancelevel.Specification},
		{"/packages/0/platforms/0/toolsDependencies/0/name", 1, compliancelevel.Strict},

		{"/packages/0/tools/0/systems/0/archiveFileName", 1, compliancelevel.Permissive},
		{"/packages/0/tools/0/systems/0/archiveFileName", 1, compliancelevel.Specification},
		{"/packages/0/tools/0/systems/0/archiveFileName", 1, compliancelevel.Strict},

		{"/packages/0/tools/0/name", 1, compliancelevel.Permissive},
		{"/packages/0/tools/0/name", 1, compliancelevel.Specification},
		{"/packages/0/tools/0/name", 1, compliancelevel.Strict},
	}

	// Test schema validation results with value length < minimum.
	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
		require.NoError(t, err)
		_, err = propertyPointer.Set(packageIndex, strings.Repeat("a", testTable.minLength-1))
		require.NoError(t, err)

		t.Run(fmt.Sprintf("%s less than minimum length of %d (%s)", testTable.propertyPointerString, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.True(t, schema.PropertyLessThanMinLength(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})

		// Test schema validation results with minimum value length.
		propertyPointer.Set(packageIndex, strings.Repeat("a", testTable.minLength))

		t.Run(fmt.Sprintf("%s at minimum length of %d (%s)", testTable.propertyPointerString, testTable.minLength, testTable.complianceLevel), func(t *testing.T) {
			assert.False(t, schema.PropertyLessThanMinLength(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})
	}
}

func TestRequired(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		complianceLevel       compliancelevel.Type
		assertion             assert.BoolAssertionFunc
	}{
		{"/packages", compliancelevel.Permissive, assert.True},
		{"/packages", compliancelevel.Specification, assert.True},
		{"/packages", compliancelevel.Strict, assert.True},

		{"/packages/0/name", compliancelevel.Permissive, assert.True},
		{"/packages/0/name", compliancelevel.Specification, assert.True},
		{"/packages/0/name", compliancelevel.Strict, assert.True},

		{"/packages/0/maintainer", compliancelevel.Permissive, assert.True},
		{"/packages/0/maintainer", compliancelevel.Specification, assert.True},
		{"/packages/0/maintainer", compliancelevel.Strict, assert.True},

		{"/packages/0/websiteURL", compliancelevel.Permissive, assert.True},
		{"/packages/0/websiteURL", compliancelevel.Specification, assert.True},
		{"/packages/0/websiteURL", compliancelevel.Strict, assert.True},

		{"/packages/0/email", compliancelevel.Permissive, assert.True},
		{"/packages/0/email", compliancelevel.Specification, assert.True},
		{"/packages/0/email", compliancelevel.Strict, assert.True},

		{"/packages/0/help", compliancelevel.Permissive, assert.False},
		{"/packages/0/help", compliancelevel.Specification, assert.False},
		{"/packages/0/help", compliancelevel.Strict, assert.False},

		{"/packages/0/help/online", compliancelevel.Permissive, assert.True},
		{"/packages/0/help/online", compliancelevel.Specification, assert.True},
		{"/packages/0/help/online", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms", compliancelevel.Strict, assert.True},

		{"/packages/0/tools", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools", compliancelevel.Specification, assert.True},
		{"/packages/0/tools", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/name", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/name", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/name", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/architecture", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/architecture", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/architecture", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/version", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/version", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/version", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/category", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/category", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/category", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/help", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/help", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/help", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/help/online", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/help/online", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/help/online", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/url", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/url", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/url", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/archiveFileName", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/archiveFileName", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/archiveFileName", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/checksum", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/checksum", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/checksum", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/size", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/size", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/size", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/boards", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/boards", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/boards", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/toolsDependencies", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/toolsDependencies", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/toolsDependencies", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/toolsDependencies/0/packager", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/packager", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/packager", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/toolsDependencies/0/name", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/name", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/name", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/toolsDependencies/0/version", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/version", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/version", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/name", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/name", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/name", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/version", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/version", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/version", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/host", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/host", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/host", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/url", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/url", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/url", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/archiveFileName", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/archiveFileName", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/archiveFileName", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/size", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/size", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/size", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/checksum", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/checksum", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/checksum", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
		require.NoError(t, err)
		_, err = propertyPointer.Delete(packageIndex)
		require.NoError(t, err)

		validationResult := packageindex.Validate(packageIndex)
		t.Run(fmt.Sprintf("%s (%s)", testTable.propertyPointerString, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.RequiredPropertyMissing(testTable.propertyPointerString, validationResult[testTable.complianceLevel]))
		})
	}
}

func TestEnum(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		propertyValue         string
		complianceLevel       compliancelevel.Type
		assertion             assert.BoolAssertionFunc
	}{
		{"/packages/0/platforms/0/category", "Contributed", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/category", "Contributed", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/category", "Contributed", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/category", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/category", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/category", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
		require.NoError(t, err)
		_, err = propertyPointer.Set(packageIndex, testTable.propertyValue)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyPointerString, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyEnumMismatch(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})
	}
}

func TestPattern(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		propertyValue         string
		complianceLevel       compliancelevel.Type
		assertion             assert.BoolAssertionFunc
	}{
		{"/packages/0/name", "foo", compliancelevel.Permissive, assert.False},
		{"/packages/0/name", "foo", compliancelevel.Specification, assert.False},
		{"/packages/0/name", "foo", compliancelevel.Strict, assert.False},

		{"/packages/0/name", "arduino", compliancelevel.Permissive, assert.False},
		{"/packages/0/name", "arduino", compliancelevel.Specification, assert.True},
		{"/packages/0/name", "arduino", compliancelevel.Strict, assert.True},

		{"/packages/0/name", "Arduino", compliancelevel.Permissive, assert.False},
		{"/packages/0/name", "Arduino", compliancelevel.Specification, assert.True},
		{"/packages/0/name", "Arduino", compliancelevel.Strict, assert.True},

		{"/packages/0/name", "ARDUINO", compliancelevel.Permissive, assert.False},
		{"/packages/0/name", "ARDUINO", compliancelevel.Specification, assert.True},
		{"/packages/0/name", "ARDUINO", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/archiveFileName", "foo.tar.bz2", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.tar.bz2", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.tar.bz2", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/archiveFileName", "foo.tar.gz", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.tar.gz", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.tar.gz", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/archiveFileName", "foo.zip", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.zip", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/archiveFileName", "foo.zip", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/archiveFileName", "foo.bar", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/archiveFileName", "foo.bar", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/archiveFileName", "foo.bar", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/checksum", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/checksum", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/checksum", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/size", "42", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/size", "42", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/size", "42", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/size", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/size", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/size", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.bz2", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.bz2", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.bz2", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.gz", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.gz", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.tar.gz", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.zip", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.zip", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.zip", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.bar", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.bar", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/archiveFileName", "foo.bar", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "SHA-256:de8a9b982477762d3d3e52fc2b682cdd8ff194dc3f1d46f4debdea6a01b33c14", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "SHA-1:f89bb8563bf86eb097679dce9d2b29b86d06bf66", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/checksum", "MD5:6c0f556759894aa1a45e8af423a08ce8", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/checksum", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/checksum", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/checksum", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/size", "42", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/size", "42", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/size", "42", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/size", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/size", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/size", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/host", "arm-linux-gnueabihf", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "arm-linux-gnueabihf", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "arm-linux-gnueabihf", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "aarch64-linux-gnu", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "aarch64-linux-gnu", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "aarch64-linux-gnu", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "arm64-linux-gnu", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "arm64-linux-gnu", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "arm64-linux-gnu", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "x86_64-linux-gnu", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "x86_64-linux-gnu", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "x86_64-linux-gnu", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "i686-mingw32", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i686-mingw32", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i686-mingw32", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "i686-cygwin", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i686-cygwin", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i686-cygwin", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "x86_64-apple-darwin", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "x86_64-apple-darwin", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "x86_64-apple-darwin", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "i386-apple-darwin11", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i386-apple-darwin11", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i386-apple-darwin11", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "i386-freebsd11", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i386-freebsd11", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "i386-freebsd11", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "amd64-freebsd11", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/host", "amd64-freebsd11", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/host", "amd64-freebsd11", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/host", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/host", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/host", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
		require.NoError(t, err)
		_, err = propertyPointer.Set(packageIndex, testTable.propertyValue)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyPointerString, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyPatternMismatch(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})
	}
}

func TestType(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		propertyValue         interface{}
		assertion             assert.BoolAssertionFunc
	}{
		{"/packages", 42, assert.True},
		{"/packages/0/name", 42, assert.True},
		{"/packages/0/maintainer", 42, assert.True},
		{"/packages/0/websiteURL", 42, assert.True},
		{"/packages/0/email", 42, assert.True},
		{"/packages/0/help", 42, assert.True},
		{"/packages/0/help/online", 42, assert.True},
		{"/packages/0/platforms", 42, assert.True},
		{"/packages/0/platforms/0/name", 42, assert.True},
		{"/packages/0/platforms/0/architecture", 42, assert.True},
		{"/packages/0/platforms/0/version", 42, assert.True},
		{"/packages/0/platforms/0/help", 42, assert.True},
		{"/packages/0/platforms/0/help/online", 42, assert.True},
		{"/packages/0/platforms/0/category", 42, assert.True},
		{"/packages/0/platforms/0/url", 42, assert.True},
		{"/packages/0/platforms/0/archiveFileName", 42, assert.True},
		{"/packages/0/platforms/0/checksum", 42, assert.True},
		{"/packages/0/platforms/0/size", 42, assert.True},
		{"/packages/0/platforms/0/boards", 42, assert.True},
		{"/packages/0/platforms/0/boards/0/name", 42, assert.True},
		{"/packages/0/platforms/0/toolsDependencies", 42, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0/packager", 42, assert.True},
		{"/packages/0/tools", 42, assert.True},
		{"/packages/0/tools/0/name", 42, assert.True},
		{"/packages/0/tools/0/version", 42, assert.True},
		{"/packages/0/tools/0/systems", 42, assert.True},
		{"/packages/0/tools/0/systems/0/host", 42, assert.True},
		{"/packages/0/tools/0/systems/0/url", 42, assert.True},
		{"/packages/0/tools/0/systems/0/archiveFileName", 42, assert.True},
		{"/packages/0/tools/0/systems/0/checksum", 42, assert.True},
		{"/packages/0/tools/0/systems/0/size", 42, assert.True},
	}

	for _, testTable := range testTables {
		for _, complianceLevel := range []compliancelevel.Type{compliancelevel.Permissive, compliancelevel.Specification, compliancelevel.Strict} {
			var packageIndex map[string]interface{}
			err := json.Unmarshal(validIndexRaw, &packageIndex)
			require.NoError(t, err)

			propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
			require.NoError(t, err)
			_, err = propertyPointer.Set(packageIndex, testTable.propertyValue)

			t.Run(fmt.Sprintf("%s: %v (%s)", testTable.propertyPointerString, testTable.propertyValue, complianceLevel), func(t *testing.T) {
				testTable.assertion(t, schema.PropertyTypeMismatch(testTable.propertyPointerString, packageindex.Validate(packageIndex)[complianceLevel]))
			})
		}
	}
}

func TestFormat(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		propertyValue         string
		complianceLevel       compliancelevel.Type
		assertion             assert.BoolAssertionFunc
	}{
		{"/packages/0/websiteURL", "http://example.com", compliancelevel.Permissive, assert.False},
		{"/packages/0/websiteURL", "http://example.com", compliancelevel.Specification, assert.False},
		{"/packages/0/websiteURL", "http://example.com", compliancelevel.Strict, assert.False},

		{"/packages/0/websiteURL", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/websiteURL", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/websiteURL", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/help/online", "http://example.com", compliancelevel.Permissive, assert.False},
		{"/packages/0/help/online", "http://example.com", compliancelevel.Specification, assert.False},
		{"/packages/0/help/online", "http://example.com", compliancelevel.Strict, assert.False},

		{"/packages/0/help/online", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/help/online", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/help/online", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/help/online", "http://example.com", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/help/online", "http://example.com", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/help/online", "http://example.com", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/help/online", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/help/online", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/help/online", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Permissive, assert.False},
		{"/packages/0/platforms/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Specification, assert.False},
		{"/packages/0/platforms/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Strict, assert.False},

		{"/packages/0/platforms/0/url", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/url", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/url", "foo", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Permissive, assert.False},
		{"/packages/0/tools/0/systems/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Specification, assert.False},
		{"/packages/0/tools/0/systems/0/url", "http://example.com/foo.tar.bz2", compliancelevel.Strict, assert.False},

		{"/packages/0/tools/0/systems/0/url", "foo", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0/url", "foo", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0/url", "foo", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString)
		require.NoError(t, err)
		_, err = propertyPointer.Set(packageIndex, testTable.propertyValue)
		require.NoError(t, err)

		t.Run(fmt.Sprintf("%s: %s (%s)", testTable.propertyPointerString, testTable.propertyValue, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.PropertyFormatMismatch(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})
	}
}

func TestAdditionalProperties(t *testing.T) {
	testTables := []struct {
		propertyPointerString string
		complianceLevel       compliancelevel.Type
		assertion             assert.BoolAssertionFunc
	}{
		// Root
		{"", compliancelevel.Permissive, assert.True},
		{"", compliancelevel.Specification, assert.True},
		{"", compliancelevel.Strict, assert.True},

		{"/packages/0", compliancelevel.Permissive, assert.True},
		{"/packages/0", compliancelevel.Specification, assert.True},
		{"/packages/0", compliancelevel.Strict, assert.True},

		{"/packages/0/help", compliancelevel.Permissive, assert.True},
		{"/packages/0/help", compliancelevel.Specification, assert.True},
		{"/packages/0/help", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/help", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/help", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/help", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/boards/0", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/boards/0", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/boards/0", compliancelevel.Strict, assert.True},

		{"/packages/0/platforms/0/toolsDependencies/0", compliancelevel.Permissive, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0", compliancelevel.Specification, assert.True},
		{"/packages/0/platforms/0/toolsDependencies/0", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0", compliancelevel.Strict, assert.True},

		{"/packages/0/tools/0/systems/0", compliancelevel.Permissive, assert.True},
		{"/packages/0/tools/0/systems/0", compliancelevel.Specification, assert.True},
		{"/packages/0/tools/0/systems/0", compliancelevel.Strict, assert.True},
	}

	for _, testTable := range testTables {
		var packageIndex map[string]interface{}
		err := json.Unmarshal(validIndexRaw, &packageIndex)
		require.NoError(t, err)

		// Add an additional property to the object.
		propertyPointer, err := gojsonpointer.NewJsonPointer(testTable.propertyPointerString + "/fooAdditionalProperty")
		require.NoError(t, err)
		_, err = propertyPointer.Set(packageIndex, "bar")
		require.NoError(t, err)

		t.Run(fmt.Sprintf("Additional property in the %s object (%s)", testTable.propertyPointerString, testTable.complianceLevel), func(t *testing.T) {
			testTable.assertion(t, schema.ProhibitedAdditionalProperties(testTable.propertyPointerString, packageindex.Validate(packageIndex)[testTable.complianceLevel]))
		})
	}
}
