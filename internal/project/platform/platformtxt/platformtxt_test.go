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

package platformtxt

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path

var validPlatformTxtMap map[string]string

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = workingDirectory.Join("testdata")

	validPlatformTxtMap = map[string]string{
		"name":                                "Arduino AVR Boards",
		"version":                             "1.8.3",
		"compiler.warning_flags.none":         "asdf",
		"compiler.warning_flags.default":      "asdf",
		"compiler.warning_flags.more":         "asdf",
		"compiler.warning_flags.all":          "asdf",
		"compiler.c.extra_flags":              "",
		"compiler.c.elf.extra_flags":          "",
		"compiler.S.extra_flags":              "",
		"compiler.cpp.extra_flags":            "",
		"compiler.ar.extra_flags":             "",
		"compiler.objcopy.eep.extra_flags":    "",
		"compiler.elf2hex.extra_flags":        "",
		"recipe.c.o.pattern":                  "asdf {compiler.c.extra_flags}",
		"recipe.cpp.o.pattern":                "asdf {compiler.cpp.extra_flags}",
		"recipe.S.o.pattern":                  "asdf {compiler.S.extra_flags}",
		"recipe.ar.pattern":                   "asdf {compiler.ar.extra_flags}",
		"recipe.c.combine.pattern":            "asdf {compiler.c.elf.extra_flags}",
		"recipe.objcopy.eep.pattern":          "asdf",
		"recipe.objcopy.hex.pattern":          "asdf",
		"recipe.output.tmp_file":              "asdf",
		"recipe.output.save_file":             "asdf",
		"recipe.size.pattern":                 "asdf",
		"recipe.size.regex":                   "asdf",
		"recipe.size.regex.data":              "asdf",
		"pluggable_discovery.required.0":      "builtin:serial-discovery",
		"pluggable_discovery.required.1":      "builtin:mdns-discovery",
		"tools.avrdude.upload.params.verbose": "-v",
		"tools.avrdude.upload.params.quiet":   "-q -q",
		"tools.avrdude.upload.pattern":        "asdf",
	}
}

func TestProperties(t *testing.T) {
	propertiesOutput, err := Properties(testDataPath.Join("valid"))
	require.Nil(t, err)

	assert.True(t, properties.NewFromHashmap(validPlatformTxtMap).Equals(propertiesOutput))
}

func TestValidate(t *testing.T) {
	platformTxt := properties.NewFromHashmap(validPlatformTxtMap)
	validationResult := Validate(platformTxt)

	assert.Nil(t, validationResult[compliancelevel.Permissive].Result, "Valid (permissive)")
	assert.Nil(t, validationResult[compliancelevel.Specification].Result, "Valid (specification)")
	assert.Nil(t, validationResult[compliancelevel.Strict].Result, "Valid (strict)")

	platformTxt.Remove("name") // Remove required property.
	validationResult = Validate(platformTxt)
	assert.NotNil(t, validationResult[compliancelevel.Permissive].Result, "Invalid (permissive)")
	assert.NotNil(t, validationResult[compliancelevel.Specification].Result, "Invalid (specification)")
	assert.NotNil(t, validationResult[compliancelevel.Strict].Result, "Invalid (strict)")
}

func TestPluggableDiscoveryNames(t *testing.T) {
	platformTxt := properties.NewFromHashmap(validPlatformTxtMap)

	assert.ElementsMatch(t, []string{}, PluggableDiscoveryNames(platformTxt), "No elements for pluggable_discovery.required properties.")

	platformTxt.Set("pluggable_discovery.foo_discovery.pattern", "asdf")
	platformTxt.Set("pluggable_discovery.bar_discovery.pattern", "zxcv")
	assert.ElementsMatch(t, []string{"foo_discovery", "bar_discovery"}, PluggableDiscoveryNames(platformTxt), "pluggable_discovery.DISCOVERY_ID properties add elements for each DISCOVERY_ID.")
}

func TestToolNames(t *testing.T) {
	platformTxt := properties.NewFromHashmap(validPlatformTxtMap)

	assert.ElementsMatch(t, []string{"avrdude"}, ToolNames(platformTxt))

	platformTxt.Set("tools.bossac.program.params.verbose", "asdf")
	platformTxt.Set("tools.bossac.program.params.quiet", "asdf")
	platformTxt.Set("tools.bossac.program.pattern", "asdf")
	assert.ElementsMatch(t, []string{"avrdude", "bossac"}, ToolNames(platformTxt))
}
