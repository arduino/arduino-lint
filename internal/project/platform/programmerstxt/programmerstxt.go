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

/*
Package programmerstxt provides functions specific to linting the programmers.txt configuration files of Arduino programmers platforms.
See: https://arduino.github.io/arduino-cli/latest/platform-specification/#programmerstxt
*/
package programmerstxt

import (
	"github.com/arduino/arduino-lint/internal/project/general"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/arduino-lint/internal/rule/schema/schemadata"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// Properties parses the programmers.txt from the given path and returns the data.
func Properties(platformPath *paths.Path) (*properties.Map, error) {
	return properties.LoadFromPath(platformPath.Join("programmers.txt"))
}

var schemaObject = make(map[compliancelevel.Type]schema.Schema)

// Validate validates programmers.txt data against the JSON schema and returns the result.
func Validate(programmersTxt *properties.Map) map[compliancelevel.Type]schema.ValidationResult {
	referencedSchemaFilenames := []string{
		"arduino-programmers-txt-definitions-schema.json",
	}

	var validationResults = make(map[compliancelevel.Type]schema.ValidationResult)

	if schemaObject[compliancelevel.Permissive].Compiled == nil { // Only compile the schemas once.
		schemaObject[compliancelevel.Permissive] = schema.Compile("arduino-programmers-txt-permissive-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Specification] = schema.Compile("arduino-programmers-txt-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Strict] = schema.Compile("arduino-programmers-txt-strict-schema.json", referencedSchemaFilenames, schemadata.Asset)
	}

	//Convert the programmers.txt data from the native properties.Map type to the interface type required by the schema validation package.
	programmersTxtInterface := general.PropertiesToMap(programmersTxt, 2)

	validationResults[compliancelevel.Permissive] = schema.Validate(programmersTxtInterface, schemaObject[compliancelevel.Permissive])
	validationResults[compliancelevel.Specification] = schema.Validate(programmersTxtInterface, schemaObject[compliancelevel.Specification])
	validationResults[compliancelevel.Strict] = schema.Validate(programmersTxtInterface, schemaObject[compliancelevel.Strict])

	return validationResults
}

// ProgrammerIDs returns the list of programmer IDs from the given programmers.txt properties.
func ProgrammerIDs(programmersTxt *properties.Map) []string {
	return programmersTxt.FirstLevelKeys()
}
