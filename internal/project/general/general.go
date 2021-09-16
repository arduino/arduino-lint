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
having both a string and a map type, which is not supported. For this reason, it's necessary to manually configure the
recursion of key levels on a case-by-case basis.
In the event a full recursion of key levels is desired, set the levels argument to a value <1.
*/
func PropertiesToMap(flatProperties *properties.Map, levels int) map[string]interface{} {
	propertiesInterface := make(map[string]interface{})

	if levels != 1 {
		for _, key := range flatProperties.FirstLevelKeys() {
			subTree := flatProperties.SubTree(key)
			if subTree.Size() > 0 {
				// This key contains a map.
				propertiesInterface[key] = PropertiesToMap(subTree, levels-1)
			} else {
				// This key contains a string, no more recursion is possible.
				propertiesInterface[key] = flatProperties.Get(key)
			}
		}
	} else {
		for _, key := range flatProperties.Keys() {
			propertiesInterface[key] = flatProperties.Get(key)
		}
	}

	return propertiesInterface
}

// PropertiesToList parses a property that has a list data type and returns it in the map[string]interface{} type
// consumed by the JSON schema parser.
func PropertiesToList(flatProperties *properties.Map, key string) map[string]interface{} {
	list := flatProperties.ExtractSubIndexLists(key)
	// Convert the slice to the required interface type
	listInterface := make([]interface{}, len(list))
	for i, v := range list {
		listInterface[i] = v
	}
	mapInterface := make(map[string]interface{})
	mapInterface[key] = listInterface
	return mapInterface
}
