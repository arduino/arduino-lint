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

// Package general provides functions that apply to multiple project types.
package general

import (
	"github.com/arduino/go-properties-orderedmap"
)

/*
PropertiesToFirstLevelExpandedMap converts properties.Map data structures to map[string]interface that maps between .
Even though boards/properties.txt have a multi-level nested data structure, the format has the odd characteristic of
allowing a key to be both an object and a string simultaneously, which is not compatible with Golang maps or JSON. So
the data structure used is a map of the first level keys (necessary to accommodate the board/prograrmmer IDs) to the
full remainder of the keys (rather than recursing through each key level individually), to string values.
*/
func PropertiesToFirstLevelExpandedMap(flatProperties *properties.Map) map[string]interface{} {
	propertiesInterface := make(map[string]interface{})
	keys := flatProperties.FirstLevelKeys()
	for _, key := range keys {
		subtreeMap := flatProperties.SubTree(key).AsMap()
		// This level also must be converted to map[string]interface{}.
		subtreeInterface := make(map[string]interface{})
		for subtreeKey, subtreeValue := range subtreeMap {
			subtreeInterface[subtreeKey] = subtreeValue
		}
		propertiesInterface[key] = subtreeInterface
	}

	return propertiesInterface
}
