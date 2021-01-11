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

package boardstxt

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path

var validBoardsTxtMap map[string]string

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = workingDirectory.Join("testdata")

	validBoardsTxtMap = map[string]string{
		"uno.name":                     "Arduino Uno",
		"uno.build.board":              "AVR_UNO",
		"uno.build.core":               "arduino",
		"uno.upload.tool":              "avrdude",
		"uno.upload.maximum_size":      "123",
		"uno.upload.maximum_data_size": "123",
	}
}

func TestProperties(t *testing.T) {
	propertiesOutput, err := Properties(testDataPath.Join("valid"))
	require.Nil(t, err)

	assert.True(t, properties.NewFromHashmap(validBoardsTxtMap).Equals(propertiesOutput))
}

func TestValidate(t *testing.T) {
	boardsTxt := properties.NewFromHashmap(validBoardsTxtMap)
	validationResult := Validate(boardsTxt)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result)
	assert.Nil(t, validationResult[compliancelevel.Specification].Result)
	assert.Nil(t, validationResult[compliancelevel.Strict].Result)

	boardsTxt.Remove("uno.name") // Remove required property.
	validationResult = Validate(boardsTxt)
	assert.NotNil(t, validationResult[compliancelevel.Permissive].Result)
	assert.NotNil(t, validationResult[compliancelevel.Specification].Result)
	assert.NotNil(t, validationResult[compliancelevel.Strict].Result)
}

func TestMenuIDs(t *testing.T) {
	boardsTxt := properties.NewFromHashmap(validBoardsTxtMap)

	assert.ElementsMatch(t, []string{}, MenuIDs(boardsTxt), "No menu IDs")

	boardsTxt.Set("menu", "noooo")
	assert.ElementsMatch(t, []string{}, MenuIDs(boardsTxt), "Some silly defined a menu property without a subproperty")

	boardsTxt.Set("menu.foo", "asdf")
	boardsTxt.Set("menu.bar", "zxcv")
	boardsTxt.Set("baz.name", "qwer")
	assert.ElementsMatch(t, []string{"foo", "bar"}, MenuIDs(boardsTxt), "Has menu IDs")
}

func TestBoardIDs(t *testing.T) {
	boardsTxt := properties.NewFromHashmap(validBoardsTxtMap)

	assert.ElementsMatch(t, []string{"uno"}, BoardIDs(boardsTxt))

	boardsTxt.Set("menu.foo", "asdf")
	boardsTxt.Set("menu.bar", "zxcv")
	boardsTxt.Set("baz.name", "qwer")
	assert.ElementsMatch(t, []string{"uno", "baz"}, BoardIDs(boardsTxt))
}

func TestVisibleBoardIDs(t *testing.T) {
	boardsTxt := properties.NewFromHashmap(validBoardsTxtMap)

	assert.ElementsMatch(t, []string{"uno"}, VisibleBoardIDs(boardsTxt))

	boardsTxt.Set("menu.foo", "asdf")
	boardsTxt.Set("menu.bar", "zxcv")
	boardsTxt.Set("baz.name", "qwer")
	boardsTxt.Set("bat.name", "sdfg")
	boardsTxt.Set("bat.hide", "")
	assert.ElementsMatch(t, []string{"uno", "baz"}, VisibleBoardIDs(boardsTxt))
}
