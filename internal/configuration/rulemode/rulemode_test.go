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

package rulemode

import (
	"reflect"
	"testing"

	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	for key := range Types {
		_, valid := Types[key]
		assert.True(t, valid)
	}

	_, valid := Types[Strict]
	assert.True(t, valid)
	_, valid = Types[42]
	assert.False(t, valid)
}

func TestMode(t *testing.T) {
	defaultRuleModes := map[projecttype.Type]map[Type]bool{
		projecttype.Sketch: {
			LibraryManagerSubmission: false,
			LibraryManagerIndexed:    false,
			Official:                 false,
		},
		projecttype.Library: {
			LibraryManagerSubmission: true,
			LibraryManagerIndexed:    false,
			Official:                 false,
		},
	}

	customRuleModes := make(map[Type]bool)

	testProjectType := projecttype.Library

	mergedRuleModes := Modes(defaultRuleModes, customRuleModes, testProjectType)

	assert.True(t, reflect.DeepEqual(defaultRuleModes[testProjectType], mergedRuleModes), "Default configuration should be used when no custom configuration was set.")

	testRuleMode := Official
	customRuleModes[testRuleMode] = !defaultRuleModes[testProjectType][testRuleMode]
	mergedRuleModes = Modes(defaultRuleModes, customRuleModes, testProjectType)
	assert.Equal(t, customRuleModes[testRuleMode], mergedRuleModes[testRuleMode], "Should be set to custom value")
}

func TestCompliance(t *testing.T) {
	ruleModes := map[Type]bool{
		Strict:        false,
		Specification: false,
		Permissive:    false,
	}

	assert.Panics(t, func() { Compliance(ruleModes) })
	ruleModes[Strict] = true
	assert.Equal(t, Strict.String(), Compliance(ruleModes))
	ruleModes[Strict] = false
	ruleModes[Specification] = true
	assert.Equal(t, Specification.String(), Compliance(ruleModes))
	ruleModes[Specification] = false
	ruleModes[Permissive] = true
	assert.Equal(t, Permissive.String(), Compliance(ruleModes))
}
