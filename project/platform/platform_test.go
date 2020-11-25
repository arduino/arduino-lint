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

package platform

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

func TestIsConfigurationFile(t *testing.T) {
	testTables := []struct {
		filename  string
		assertion assert.BoolAssertionFunc
	}{
		{"boards.txt", assert.True},
		{"boards.local.txt", assert.True},
		{"platform.txt", assert.True},
		{"platform.local.txt", assert.True},
		{"programmers.txt", assert.True},
		{"foo.txt", assert.False},
	}

	for _, testTable := range testTables {
		testTable.assertion(t, IsConfigurationFile(paths.New("/foo", testTable.filename)), testTable.filename)
	}
}

func TestIsRequiredConfigurationFile(t *testing.T) {
	assert.True(t, IsRequiredConfigurationFile(paths.New("/foo", "boards.txt")))
	assert.False(t, IsRequiredConfigurationFile(paths.New("/foo", "platform.txt")))
}
