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

/*
Package boardstxt provides functions specific to linting the boards.txt configuration files of Arduino boards platforms.
See: https://arduino.github.io/arduino-cli/latest/platform-specification/#boardstxt
*/
package boardstxt

import (
	"strings"

	"github.com/arduino/arduino-lint/internal/project/general"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/arduino-lint/internal/rule/schema/schemadata"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// Properties parses the boards.txt from the given path and returns the data.
func Properties(platformPath *paths.Path) (*properties.Map, error) {
	return properties.LoadFromPath(platformPath.Join("boards.txt"))
}

var schemaObject = make(map[compliancelevel.Type]schema.Schema)

// Validate validates boards.txt data against the JSON schema and returns a map of the result for each compliance level.
func Validate(boardsTxt *properties.Map) map[compliancelevel.Type]schema.ValidationResult {
	referencedSchemaFilenames := []string{
		"general-definitions-schema.json",
		"arduino-boards-txt-definitions-schema.json",
	}

	var validationResults = make(map[compliancelevel.Type]schema.ValidationResult)

	if schemaObject[compliancelevel.Permissive].Compiled == nil { // Only compile the schemas once.
		schemaObject[compliancelevel.Permissive] = schema.Compile("arduino-boards-txt-permissive-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Specification] = schema.Compile("arduino-boards-txt-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Strict] = schema.Compile("arduino-boards-txt-strict-schema.json", referencedSchemaFilenames, schemadata.Asset)
	}

	//Convert the boards.txt data from the native properties.Map type to the interface type required by the schema validation package.
	boardsTxtInterface := make(map[string]interface{})
	keys := boardsTxt.FirstLevelKeys()
	for _, key := range keys {
		if key == "menu" {
			// Menu title subproperties are flat.
			boardsTxtInterface[key] = general.PropertiesToMap(boardsTxt.SubTree(key), 1)
		} else {
			boardIDInterface := make(map[string]interface{})
			boardIDProperties := boardsTxt.SubTree(key)
			boardIDKeys := boardIDProperties.Keys()

			// Add the standard properties for the board.
			for _, boardIDKey := range boardIDKeys {
				if !strings.HasPrefix(boardIDKey, "menu.") {
					boardIDInterface[boardIDKey] = boardIDProperties.Get(boardIDKey)
				}
			}

			// Add the custom board option properties for the board, nested down to OPTION_ID.
			boardIDInterface["menu"] = general.PropertiesToMap(boardIDProperties.SubTree("menu"), 3)

			boardsTxtInterface[key] = boardIDInterface
		}
	}

	validationResults[compliancelevel.Permissive] = schema.Validate(boardsTxtInterface, schemaObject[compliancelevel.Permissive])
	validationResults[compliancelevel.Specification] = schema.Validate(boardsTxtInterface, schemaObject[compliancelevel.Specification])
	validationResults[compliancelevel.Strict] = schema.Validate(boardsTxtInterface, schemaObject[compliancelevel.Strict])

	return validationResults
}

// MenuIDs returns the list of menu IDs from the given boards.txt properties.
func MenuIDs(boardsTxt *properties.Map) []string {
	// Each menu must have a property defining its title with the format `menu.MENU_ID=MENU_TITLE`.
	return boardsTxt.SubTree("menu").FirstLevelKeys()
}

// BoardIDs returns the list of board IDs from the given boards.txt properties.
func BoardIDs(boardsTxt *properties.Map) []string {
	boardIDs := boardsTxt.FirstLevelKeys()
	boardIDCount := 0
	for _, boardID := range boardIDs {
		if boardID != "menu" {
			// This element is a board ID, retain it in the section of the array that will be returned.
			boardIDs[boardIDCount] = boardID
			boardIDCount++
		}
	}

	return boardIDs[:boardIDCount]
}

// VisibleBoardIDs returns the list of IDs for non-hidden boards from the given boards.txt properties.
func VisibleBoardIDs(boardsTxt *properties.Map) []string {
	boardIDs := BoardIDs(boardsTxt)
	boardIDCount := 0
	for _, boardID := range boardIDs {
		if !boardsTxt.ContainsKey(boardID + ".hide") {
			// This element is a visible board, retain it in the section of the array that will be returned.
			boardIDs[boardIDCount] = boardID
			boardIDCount++
		}
	}

	return boardIDs[:boardIDCount]
}
