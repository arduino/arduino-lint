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

// Package general provides functions that apply to multiple project types.
package general

import (
	"github.com/arduino/go-properties-orderedmap"
)

/*
PropertiesToMap converts properties.Map data structures to map[string]interface with the specified number of key levels.
The Arduino project configuration fields have an odd usage of the properties.Map format. Dots may sometimes indicate
nested keys, but in other cases they are merely a character in the key string. There are cases where a completely
programmatic recursion of the properties into a fully nested structure would result in the impossibility of some keys
having bot a string and a map type, which is not supported. For this reason, it's necessary to manually configure the
recursion of key levels on a case-by-case basis.
In the event a full recursion of key levels is desired, set the levels argument to a value <1.
*/
func PropertiesToMap(flatProperties *properties.Map, levels int) map[string]interface{} {
	propertiesInterface := make(map[string]interface{})

	var keys []string
	if levels != 1 {
		keys = flatProperties.FirstLevelKeys()
	} else {
		keys = flatProperties.Keys()
	}

	for _, key := range keys {
		subTree := flatProperties.SubTree(key)
		if subTree.Size() > 0 {
			// This key contains a map.
			propertiesInterface[key] = PropertiesToMap(subTree, levels-1)
		} else {
			// This key contains a string, no more recursion is possible.
			propertiesInterface[key] = flatProperties.Get(key)
		}
	}

	return propertiesInterface
}
