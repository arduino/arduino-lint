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

// Package libraryproperties provides functions for working with the library.properties Arduino library metadata file.
package libraryproperties

import (
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/arduino-lint/internal/rule/schema/schemadata"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// Properties parses the library.properties from the given path and returns the data.
func Properties(libraryPath *paths.Path) (*properties.Map, error) {
	return properties.SafeLoadFromPath(libraryPath.Join("library.properties"))
}

var schemaObject = make(map[compliancelevel.Type]schema.Schema)

// Validate validates library.properties data against the JSON schema and returns a map of the result for each compliance level.
func Validate(libraryProperties *properties.Map) map[compliancelevel.Type]schema.ValidationResult {
	referencedSchemaFilenames := []string{
		"general-definitions-schema.json",
		"arduino-library-properties-definitions-schema.json",
	}

	var validationResults = make(map[compliancelevel.Type]schema.ValidationResult)

	if schemaObject[compliancelevel.Permissive].Compiled == nil { // Only compile the schemas once.
		schemaObject[compliancelevel.Permissive] = schema.Compile("arduino-library-properties-permissive-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Specification] = schema.Compile("arduino-library-properties-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Strict] = schema.Compile("arduino-library-properties-strict-schema.json", referencedSchemaFilenames, schemadata.Asset)
	}

	validationResults[compliancelevel.Permissive] = schema.Validate(libraryProperties, schemaObject[compliancelevel.Permissive])
	validationResults[compliancelevel.Specification] = schema.Validate(libraryProperties, schemaObject[compliancelevel.Specification])
	validationResults[compliancelevel.Strict] = schema.Validate(libraryProperties, schemaObject[compliancelevel.Strict])

	return validationResults
}
