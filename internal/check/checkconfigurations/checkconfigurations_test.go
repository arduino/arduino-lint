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

package checkconfigurations_test

import (
	"fmt"
	"testing"

	"github.com/arduino/arduino-lint/internal/check"
	"github.com/arduino/arduino-lint/internal/check/checkconfigurations"
	"github.com/arduino/arduino-lint/internal/check/checklevel"
	"github.com/arduino/arduino-lint/internal/configuration/checkmode"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestIncorrectCheckIDPrefix(t *testing.T) {
	for checkIndex, checkConfiguration := range checkconfigurations.Configurations() {
		var IDPrefix byte
		switch checkConfiguration.ProjectType {
		case projecttype.Sketch:
			IDPrefix = 'S'
		case projecttype.Library:
			IDPrefix = 'L'
		case projecttype.Platform:
			IDPrefix = 'P'
		case projecttype.PackageIndex:
			IDPrefix = 'I'
		default:
			panic(fmt.Errorf("No prefix configured for project type %s", checkConfiguration.ProjectType))
		}
		require.NotEmptyf(t, checkConfiguration.ID, "No check ID defined for check configuration #%v", checkIndex)
		assert.Equalf(t, IDPrefix, checkConfiguration.ID[0], "Check ID %s has incorrect prefix for project type %s.", checkConfiguration.ID, checkConfiguration.ProjectType)
	}
}

func TestDuplicateCheckID(t *testing.T) {
	checkIDMap := make(map[string]bool)
	for checkIndex, checkConfiguration := range checkconfigurations.Configurations() {
		checkIDMap[checkConfiguration.ID] = true
		require.Equalf(t, checkIndex+1, len(checkIDMap), "ID %s of check #%v is a duplicate", checkConfiguration.ID, checkIndex)
	}
}
