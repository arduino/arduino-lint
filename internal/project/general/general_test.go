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

package general

import (
	"reflect"
	"testing"

	"github.com/arduino/go-properties-orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertiesToMap(t *testing.T) {
	rawProperties := []byte(`
		hello=world
		goodbye=
		foo.bar=asdf
		foo.baz=zxcv
		bar.bat.bam=123
		qux.a=x
		qux.a.b=y
		fuz.a.b=y
		fuz.a=x
	`)
	propertiesInput, err := properties.LoadFromBytes(rawProperties)
	require.Nil(t, err)

	expectedMapOutput := map[string]interface{}{
		"hello":       "world",
		"goodbye":     "",
		"foo.bar":     "asdf",
		"foo.baz":     "zxcv",
		"bar.bat.bam": "123",
		"qux.a":       "x",
		"qux.a.b":     "y",
		"fuz.a.b":     "y",
		"fuz.a":       "x",
	}

	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToMap(propertiesInput, 1)))

	expectedMapOutput = map[string]interface{}{
		"hello":   "world",
		"goodbye": "",
		"foo": map[string]interface{}{
			"bar": "asdf",
			"baz": "zxcv",
		},
		"bar": map[string]interface{}{
			"bat.bam": "123",
		},
		"qux": map[string]interface{}{
			"a":   "x",
			"a.b": "y",
		},
		"fuz": map[string]interface{}{
			"a.b": "y",
			"a":   "x",
		},
	}

	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToMap(propertiesInput, 2)))

	expectedMapOutput = map[string]interface{}{
		"hello":   "world",
		"goodbye": "",
		"foo": map[string]interface{}{
			"bar": "asdf",
			"baz": "zxcv",
		},
		"bar": map[string]interface{}{
			"bat": map[string]interface{}{
				"bam": "123",
			},
		},
		"qux": map[string]interface{}{
			"a": map[string]interface{}{
				"b": "y", // It is impossible to represent the complete "properties" data structure recursed to this depth.
			},
		},
		"fuz": map[string]interface{}{
			"a": map[string]interface{}{
				"b": "y",
			},
		},
	}

	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToMap(propertiesInput, 3)))
	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToMap(propertiesInput, 0)))
}

func TestPropertiesToList(t *testing.T) {
	rawProperties := []byte(`
		hello=world
		foo.1=asdf
		foo.2=zxcv
	`)
	propertiesInput, err := properties.LoadFromBytes(rawProperties)
	require.Nil(t, err)

	expectedMapOutput := map[string]interface{}{
		"hello": []interface{}{"world"},
	}
	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToList(propertiesInput, "hello")))

	expectedMapOutput = map[string]interface{}{
		"foo": []interface{}{"asdf", "zxcv"},
	}
	assert.True(t, reflect.DeepEqual(expectedMapOutput, PropertiesToList(propertiesInput, "foo")))
}
