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

package sketch

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

func TestHasMainFileValidExtension(t *testing.T) {
	assert.True(t, HasMainFileValidExtension(paths.New("/foo/bar.ino")))
	assert.False(t, HasMainFileValidExtension(paths.New("/foo/bar.h")))
}

func TestContainsMainSketchFile(t *testing.T) {
	assert.True(t, ContainsMainSketchFile(testDataPath.Join("Valid")))
	assert.False(t, ContainsMainSketchFile(testDataPath.Join("ContainsNoMainSketchFile")))
}

func TestHasSupportedExtension(t *testing.T) {
	assert.True(t, HasSupportedExtension(paths.New("/foo/bar.ino")))
	assert.True(t, HasSupportedExtension(paths.New("/foo/bar.h")))
	assert.False(t, HasSupportedExtension(paths.New("/foo/bar.baz")))
}
