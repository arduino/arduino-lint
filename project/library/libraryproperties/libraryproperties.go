// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
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
	"github.com/arduino/arduino-check/check/checkdata/schema"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

// Properties parses the library.properties from the given path and returns the data.
func Properties(libraryPath *paths.Path) (*properties.Map, error) {
	libraryProperties, err := properties.Load(libraryPath.Join("library.properties").String())
	if err != nil {
		return nil, err
	}
	return libraryProperties, nil
}

// Validate validates library.properties data against the JSON schema.
func Validate(libraryProperties *properties.Map) *gojsonschema.Result {
	referencedSchemaFilenames := []string{}
	schemaObject := schema.Compile("arduino-library-properties-schema.json", referencedSchemaFilenames, configuration.SchemasPath())

	return schema.Validate(libraryProperties, schemaObject)
}
