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

package checkdata

import (
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/ory/jsonschema/v3"
)

// Initialize gathers the library check data for the specified project.
func InitializeForLibrary(project project.Type) {
	libraryProperties, libraryPropertiesLoadError = libraryproperties.Properties(project.Path)
	if libraryPropertiesLoadError != nil {
		// TODO: can I even do this?
		libraryPropertiesSchemaValidationResult = nil
	} else {
		libraryPropertiesSchemaValidationResult = libraryproperties.Validate(libraryProperties)
	}
}

var libraryPropertiesLoadError error

// LibraryPropertiesLoadError returns the error output from loading the library.properties metadata file.
func LibraryPropertiesLoadError() error {
	return libraryPropertiesLoadError
}

var libraryProperties *properties.Map

// LibraryProperties returns the data from the library.properties metadata file.
func LibraryProperties() *properties.Map {
	return libraryProperties
}

var libraryPropertiesSchemaValidationResult *jsonschema.ValidationError

// LibraryPropertiesSchemaValidationResult returns the result of validating library.properties against the JSON schema.
func LibraryPropertiesSchemaValidationResult() *jsonschema.ValidationError {
	return libraryPropertiesSchemaValidationResult
}
