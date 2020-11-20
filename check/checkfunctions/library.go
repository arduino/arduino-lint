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

package checkfunctions

// The check functions for libraries.

import (
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkdata/schema"
	"github.com/arduino/arduino-check/check/checkdata/schema/compliancelevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/sirupsen/logrus"
)

// LibraryPropertiesFormat checks for invalid library.properties format.
func LibraryPropertiesFormat() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.Fail, checkdata.LibraryPropertiesLoadError().Error()
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldMissing checks for missing library.properties "name" field.
func LibraryPropertiesNameFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDisallowedCharacters checks for disallowed characters in the library.properties "name" field.
func LibraryPropertiesNameFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.PropertyPatternMismatch("name", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDuplicate checks whether there is an existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldDuplicate() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, ""
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Fail, name
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldNotInIndex checks whether there is no existing entry in the Library Manager index using the the library.properties `name` value.
func LibraryPropertiesNameFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	name, hasName := checkdata.LibraryProperties().GetOk("name")
	if !hasName {
		return checkresult.NotRun, ""
	}

	if nameInLibraryManagerIndex(name) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, name
}

// LibraryPropertiesVersionFieldMissing checks for missing library.properties "version" field.
func LibraryPropertiesVersionFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if schema.RequiredPropertyMissing("version", checkdata.LibraryPropertiesSchemaValidationResult()[compliancelevel.Specification], configuration.SchemasPath()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesDependsFieldNotInIndex checks whether the libraries listed in the library.properties `depends` field are in the Library Manager index.
func LibraryPropertiesDependsFieldNotInIndex() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	depends, hasDepends := checkdata.LibraryProperties().GetOk("depends")
	if !hasDepends {
		return checkresult.NotRun, ""
	}

	dependencies, err := properties.SplitQuotedString(depends, "", false)
	if err != nil {
		panic(err)
	}
	dependenciesNotInIndex := []string{}
	for _, dependency := range dependencies {
		logrus.Tracef("Checking if dependency %s is in index.", dependency)
		if !nameInLibraryManagerIndex(dependency) {
			dependenciesNotInIndex = append(dependenciesNotInIndex, dependency)
		}
	}

	if len(dependenciesNotInIndex) > 0 {
		return checkresult.Fail, strings.Join(dependenciesNotInIndex, ", ")
	}

	return checkresult.Pass, ""
}

// nameInLibraryManagerIndex returns whether there is a library in Library Manager index using the given name.
func nameInLibraryManagerIndex(name string) bool {
	libraries := checkdata.LibraryManagerIndex()["libraries"].([]interface{})
	for _, libraryInterface := range libraries {
		library := libraryInterface.(map[string]interface{})
		if library["name"].(string) == name {
			return true
		}
	}

	return false
}
