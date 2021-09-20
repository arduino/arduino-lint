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

/*
Package platformtxt provides functions specific to linting the platform.txt configuration files of Arduino platform platforms.
See: https://arduino.github.io/arduino-cli/latest/platform-specification/#platformtxt
*/
package platformtxt

import (
	"fmt"
	"strings"

	"github.com/arduino/arduino-lint/internal/project/general"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/arduino-lint/internal/rule/schema/schemadata"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// Properties parses the platform.txt from the given path and returns the data.
func Properties(platformPath *paths.Path) (*properties.Map, error) {
	return properties.LoadFromPath(platformPath.Join("platform.txt"))
}

var schemaObject = make(map[compliancelevel.Type]schema.Schema)

// Validate validates platform.txt data against the JSON schema and returns a map of the result for each compliance level.
func Validate(platformTxt *properties.Map) map[compliancelevel.Type]schema.ValidationResult {
	referencedSchemaFilenames := []string{
		"general-definitions-schema.json",
		"arduino-platform-txt-definitions-schema.json",
	}

	var validationResults = make(map[compliancelevel.Type]schema.ValidationResult)

	if schemaObject[compliancelevel.Permissive].Compiled == nil { // Only compile the schemas once.
		schemaObject[compliancelevel.Permissive] = schema.Compile("arduino-platform-txt-permissive-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Specification] = schema.Compile("arduino-platform-txt-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Strict] = schema.Compile("arduino-platform-txt-strict-schema.json", referencedSchemaFilenames, schemadata.Asset)
	}

	/*
		Convert the platform.txt data from the native properties.Map type to the interface type required by the schema
		validation package.
		Even though platform.txt has a multi-level nested data structure, the format has the odd characteristic of allowing
		a key to be both an object and a string simultaneously, which is not compatible with Golang maps or JSON. So the
		data structure used is a selective map, using a flat map except for the tools and pluggable_discovery keys, which
		can contain any number of arbitrary subproperties which must be linted.
	*/
	platformTxtInterface := make(map[string]interface{})
	keys := platformTxt.Keys()
	for _, key := range keys {
		if strings.HasPrefix(key, "pluggable_discovery.") {
			if key == "pluggable_discovery.required" || strings.HasPrefix(key, "pluggable_discovery.required.") {
				platformTxtInterface["pluggable_discovery"] = general.PropertiesToList(platformTxt.SubTree("pluggable_discovery"), "required")
			} else {
				// It is a pluggable_discovery.DISCOVERY_ID property.
				platformTxtInterface["pluggable_discovery"] = general.PropertiesToMap(platformTxt.SubTree("pluggable_discovery"), 2)
			}
		} else if strings.HasPrefix(key, "tools.") {
			platformTxtInterface["tools"] = general.PropertiesToMap(platformTxt.SubTree("tools"), 4)
		} else {
			platformTxtInterface[key] = platformTxt.Get(key)
		}
	}

	validationResults[compliancelevel.Permissive] = schema.Validate(platformTxtInterface, schemaObject[compliancelevel.Permissive])
	validationResults[compliancelevel.Specification] = schema.Validate(platformTxtInterface, schemaObject[compliancelevel.Specification])
	validationResults[compliancelevel.Strict] = schema.Validate(platformTxtInterface, schemaObject[compliancelevel.Strict])

	return validationResults
}

// PluggableDiscoveryNames returns the list of pluggable discovery names from the given platform.txt properties.
func PluggableDiscoveryNames(platformTxt *properties.Map) []string {
	names := platformTxt.SubTree("pluggable_discovery").FirstLevelKeys()
	for i := range names {
		for i < len(names) && names[i] == "required" {
			names = append(names[:i], names[i+1:]...)
		}
	}

	return names
}

// UserProvidedFieldNames returns the list of user provided field names platform.txt properties, mapped by tool name.
func UserProvidedFieldNames(platformTxt *properties.Map) map[string][]string {
	fieldNames := make(map[string][]string)
	toolsProps := platformTxt.SubTree("tools")
	for _, tool := range toolsProps.FirstLevelKeys() {
		fieldProps := toolsProps.SubTree(fmt.Sprintf("%s.upload.field", tool))
		for _, fieldName := range fieldProps.FirstLevelKeys() {
			fieldNames[tool] = append(fieldNames[tool], fieldName)
		}
	}

	return fieldNames
}

// ToolNames returns the list of tool names from the given platform.txt properties.
func ToolNames(platformTxt *properties.Map) []string {
	return platformTxt.SubTree("tools").FirstLevelKeys()
}
