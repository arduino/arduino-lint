// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package checkconfigurations_test

import (
	"fmt"
	"testing"

	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checklevel"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationResolution(t *testing.T) {
	for _, checkConfiguration := range checkconfigurations.Configurations() {
		for checkMode := range checkmode.Types {
			enabled, err := check.IsEnabled(checkConfiguration, map[checkmode.Type]bool{checkMode: true})
			assert.Nil(t, err, fmt.Sprintf("Enable configuration of check %s doesn't resolve for check mode %s", checkConfiguration.ID, checkMode))
			if err == nil && enabled {
				_, err := checklevel.FailCheckLevel(checkConfiguration, map[checkmode.Type]bool{checkMode: true})
				assert.Nil(t, err, fmt.Sprintf("Level configuration of check %s doesn't resolve for check mode %s", checkConfiguration.ID, checkMode))
			}
		}
	}
}

func TestConfigurationCheckModeConflict(t *testing.T) {
	// Having the same check mode in multiple configurations results in the configuration behavior being dependent on which order the fields are processed in, which may change.
	for _, checkConfiguration := range checkconfigurations.Configurations() {
		conflict, checkMode := checkModeConflict(checkConfiguration.DisableModes, checkConfiguration.EnableModes)
		assert.False(t, conflict, fmt.Sprintf("Duplicated check mode %s in enable configuration of check %s", checkMode, checkConfiguration.ID))

		conflict, checkMode = checkModeConflict(checkConfiguration.InfoModes, checkConfiguration.WarningModes, checkConfiguration.ErrorModes)
		assert.False(t, conflict, fmt.Sprintf("Duplicated check mode %s in level configuration of check %s", checkMode, checkConfiguration.ID))
	}
}

// checkModeConflict checks whether the same check mode is present in multiple configuration fields.
func checkModeConflict(configurations ...[]checkmode.Type) (bool, checkmode.Type) {
	modeCounter := 0
	checkModeMap := make(map[checkmode.Type]bool)
	for _, configuration := range configurations {
		for _, checkMode := range configuration {
			modeCounter += 1
			checkModeMap[checkMode] = true
			if len(checkModeMap) < modeCounter {
				return true, checkMode
			}
		}
	}
	return false, checkmode.Default
}
