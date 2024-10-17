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

package ruleconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/rulelevel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurationResolution(t *testing.T) {
	for _, ruleConfiguration := range ruleconfiguration.Configurations() {
		for ruleMode := range rulemode.Types {
			enabled, err := rule.IsEnabled(ruleConfiguration, map[rulemode.Type]bool{ruleMode: true})
			assert.NoError(t, err, fmt.Sprintf("Enable configuration of rule %s doesn't resolve for rule mode %s", ruleConfiguration.ID, ruleMode))
			if err == nil && enabled {
				_, err := rulelevel.FailRuleLevel(ruleConfiguration, map[rulemode.Type]bool{ruleMode: true})
				assert.Nil(t, err, fmt.Sprintf("Level configuration of rule %s doesn't resolve for rule mode %s", ruleConfiguration.ID, ruleMode))
			}
		}
	}
}

func TestConfigurationRuleModeConflict(t *testing.T) {
	// Having the same rule mode in multiple configurations results in the configuration behavior being dependent on which order the fields are processed in, which may change.
	for _, ruleConfiguration := range ruleconfiguration.Configurations() {
		conflict, ruleMode := ruleModeConflict(ruleConfiguration.DisableModes, ruleConfiguration.EnableModes)
		assert.False(t, conflict, fmt.Sprintf("Duplicated rule mode %s in enable configuration of rule %s", ruleMode, ruleConfiguration.ID))

		conflict, ruleMode = ruleModeConflict(ruleConfiguration.InfoModes, ruleConfiguration.WarningModes, ruleConfiguration.ErrorModes)
		assert.False(t, conflict, fmt.Sprintf("Duplicated rule mode %s in level configuration of rule %s", ruleMode, ruleConfiguration.ID))
	}
}

// ruleModeConflict rules whether the same rule mode is present in multiple configuration fields.
func ruleModeConflict(configurations ...[]rulemode.Type) (bool, rulemode.Type) {
	modeCounter := 0
	ruleModeMap := make(map[rulemode.Type]bool)
	for _, configuration := range configurations {
		for _, ruleMode := range configuration {
			modeCounter += 1
			ruleModeMap[ruleMode] = true
			if len(ruleModeMap) < modeCounter {
				return true, ruleMode
			}
		}
	}
	return false, rulemode.Default
}

func TestIncorrectRuleIDPrefix(t *testing.T) {
	for ruleIndex, ruleConfiguration := range ruleconfiguration.Configurations() {
		var IDPrefix byte
		switch ruleConfiguration.ProjectType {
		case projecttype.Sketch:
			IDPrefix = 'S'
		case projecttype.Library:
			IDPrefix = 'L'
		case projecttype.Platform:
			IDPrefix = 'P'
		case projecttype.PackageIndex:
			IDPrefix = 'I'
		default:
			panic(fmt.Errorf("No prefix configured for project type %s", ruleConfiguration.ProjectType))
		}
		require.NotEmptyf(t, ruleConfiguration.ID, "No rule ID defined for rule configuration #%v", ruleIndex)
		assert.Equalf(t, IDPrefix, ruleConfiguration.ID[0], "Rule ID %s has incorrect prefix for project type %s.", ruleConfiguration.ID, ruleConfiguration.ProjectType)
	}
}

func TestDuplicateRuleID(t *testing.T) {
	ruleIDMap := make(map[string]bool)
	for ruleIndex, ruleConfiguration := range ruleconfiguration.Configurations() {
		ruleIDMap[ruleConfiguration.ID] = true
		require.Equalf(t, ruleIndex+1, len(ruleIDMap), "ID %s of rule #%v is a duplicate", ruleConfiguration.ID, ruleIndex)
	}
}

func TestRequiredConfigField(t *testing.T) {
	for _, ruleConfiguration := range ruleconfiguration.Configurations() {
		assert.NotEmptyf(t, ruleConfiguration.Category, "No category defined for rule %s", ruleConfiguration.ID)
		assert.NotEmptyf(t, ruleConfiguration.Subcategory, "No subcategory defined for rule %s", ruleConfiguration.ID)
		assert.NotEmptyf(t, ruleConfiguration.Brief, "No brief defined for rule %s", ruleConfiguration.ID)
		assert.NotEmptyf(t, ruleConfiguration.Description, "No description defined for rule %s", ruleConfiguration.ID)
		assert.NotEmptyf(t, ruleConfiguration.MessageTemplate, "No message template defined for rule %s", ruleConfiguration.ID)
	}
}
