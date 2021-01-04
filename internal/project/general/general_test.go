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

package general

import (
	"reflect"
	"testing"

	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertiesToFirstLevelExpandedMap(t *testing.T) {
	rawProperties := []byte(`
		foo.bar=asdf
		foo.baz=zxcv
		bar.bat.bam=123
	`)
	propertiesInput, err := properties.LoadFromBytes(rawProperties)
	require.Nil(t, err)

	expectedMapOutput := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "asdf",
			"baz": "zxcv",
		},
		"bar": map[string]interface{}{
			"bat.bam": "123",
		},
	}

	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToFirstLevelExpandedMap(propertiesInput)))
}
