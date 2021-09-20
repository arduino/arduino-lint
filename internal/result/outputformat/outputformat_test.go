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

package outputformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	testTables := []struct {
		formatString   string
		expectedFormat Type
		errorAssertion assert.ErrorAssertionFunc
	}{
		{"text", Text, assert.NoError},
		{"json", JSON, assert.NoError},
		{"TEXT", Text, assert.NoError},
		{"foo", 0, assert.Error},
	}

	for _, testTable := range testTables {
		outputFormat, err := FromString(testTable.formatString)
		testTable.errorAssertion(t, err, testTable.formatString)
		if err == nil {
			assert.Equal(t, testTable.expectedFormat, outputFormat, testTable.formatString)
		}
	}
}
