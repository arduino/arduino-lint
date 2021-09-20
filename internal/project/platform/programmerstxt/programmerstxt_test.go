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

package programmerstxt

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path

var validProgrammersTxtMap map[string]string

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = workingDirectory.Join("testdata")

	validProgrammersTxtMap = map[string]string{
		"usbasp.name":                 "USBasp",
		"usbasp.program.tool":         "avrdude",
		"usbasp.program.extra_params": "-Pusb",
		"arduinoasisp.name":           "Arduino as ISP",
		"arduinoasisp.program.tool":   "avrdude",
	}
}

func TestProperties(t *testing.T) {
	propertiesOutput, err := Properties(testDataPath.Join("valid"))
	require.Nil(t, err)

	assert.True(t, properties.NewFromHashmap(validProgrammersTxtMap).Equals(propertiesOutput))
}

func TestValidate(t *testing.T) {
	programmersTxt := properties.NewFromHashmap(validProgrammersTxtMap)
	validationResult := Validate(programmersTxt)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)

	programmersTxt.Remove("usbasp.name") // Remove required property.
	validationResult = Validate(programmersTxt)
	assert.NotNil(t, validationResult[compliancelevel.Permissive].Result)
	assert.NotNil(t, validationResult[compliancelevel.Specification].Result)
	assert.NotNil(t, validationResult[compliancelevel.Strict].Result)
}

func TestProgrammerIDs(t *testing.T) {
	programmersTxt := properties.NewFromHashmap(validProgrammersTxtMap)

	assert.ElementsMatch(t, []string{"usbasp", "arduinoasisp"}, ProgrammerIDs(programmersTxt))
}
