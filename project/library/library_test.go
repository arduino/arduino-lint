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

package library

import (
	"os"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, _ := os.Getwd()
	testDataPath = paths.New(workingDirectory, "testdata")
}

func TestContainsHeaderFile(t *testing.T) {
	assert.True(t, ContainsHeaderFile(testDataPath.Join("ContainsHeaderFile")))
	assert.False(t, ContainsHeaderFile(testDataPath.Join("ContainsNoHeaderFile")))
}

func TestContainsMetadataFile(t *testing.T) {
	assert.True(t, ContainsMetadataFile(testDataPath.Join("ContainsMetadataFile")))
	assert.False(t, ContainsMetadataFile(testDataPath.Join("ContainsNoMetadataFile")))
}
