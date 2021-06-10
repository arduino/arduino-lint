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
Package packageindex provides functions specific to linting the package index files of the Arduino Boards Manager.
See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification
*/
package packageindex

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/arduino-lint/internal/rule/schema/schemadata"
	"github.com/arduino/go-paths-helper"
)

// Properties parses the package index from the given path and returns the data.
func Properties(packageIndexPath *paths.Path) (map[string]interface{}, error) {
	if packageIndexPath == nil {
		return nil, fmt.Errorf("Package index path is nil")
	}
	rawIndex, err := packageIndexPath.ReadFile()
	if err != nil {
		return nil, err
	}
	var indexData map[string]interface{}
	err = json.Unmarshal(rawIndex, &indexData)
	if err != nil {
		return nil, err
	}

	return indexData, nil
}

var schemaObject = make(map[compliancelevel.Type]schema.Schema)

// Validate validates boards.txt data against the JSON schema and returns a map of the result for each compliance level.
func Validate(packageIndex map[string]interface{}) map[compliancelevel.Type]schema.ValidationResult {
	referencedSchemaFilenames := []string{
		"general-definitions-schema.json",
		"arduino-package-index-definitions-schema.json",
	}

	var validationResults = make(map[compliancelevel.Type]schema.ValidationResult)

	if schemaObject[compliancelevel.Permissive].Compiled == nil { // Only compile the schemas once.
		schemaObject[compliancelevel.Permissive] = schema.Compile("arduino-package-index-permissive-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Specification] = schema.Compile("arduino-package-index-schema.json", referencedSchemaFilenames, schemadata.Asset)
		schemaObject[compliancelevel.Strict] = schema.Compile("arduino-package-index-strict-schema.json", referencedSchemaFilenames, schemadata.Asset)
	}

	validationResults[compliancelevel.Permissive] = schema.Validate(packageIndex, schemaObject[compliancelevel.Permissive])
	validationResults[compliancelevel.Specification] = schema.Validate(packageIndex, schemaObject[compliancelevel.Specification])
	validationResults[compliancelevel.Strict] = schema.Validate(packageIndex, schemaObject[compliancelevel.Strict])

	return validationResults
}

var empty struct{}

// Reference: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
var validExtensions = map[string]struct{}{
	".json": empty,
}

// HasValidExtension returns whether the file at the given path has a valid package index extension.
func HasValidExtension(filePath *paths.Path) bool {
	_, hasValidExtension := validExtensions[filePath.Ext()]
	return hasValidExtension
}

// Regular expressions for official and non-official package index filenames
// See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
var validFilenameRegex = map[bool]*regexp.Regexp{
	true:  regexp.MustCompile(`^package_(.+_)*index.json$`),
	false: regexp.MustCompile(`^package_(.+_)+index.json$`),
}

// HasValidFilename returns whether the file at the given path has a valid package index filename.
func HasValidFilename(filePath *paths.Path, officialRuleMode bool) bool {
	regex := validFilenameRegex[officialRuleMode]
	filename := filePath.Base()
	return regex.MatchString(filename)
}

// Find searches the provided path for a file that has a name resembling a package index and returns the path to that file.
func Find(folderPath *paths.Path) (*paths.Path, error) {
	exist, err := folderPath.ExistCheck()
	if !exist {
		return nil, fmt.Errorf("Error opening path %s: %s", folderPath, err)
	}

	if folderPath.IsNotDir() {
		return folderPath, nil
	}

	directoryListing, err := folderPath.ReadDir()
	if err != nil {
		return nil, err
	}

	directoryListing.FilterOutDirs()
	for _, potentialPackageIndexFile := range directoryListing {
		if HasValidFilename(potentialPackageIndexFile, true) {
			return potentialPackageIndexFile, nil
		}
	}
	for _, potentialPackageIndexFile := range directoryListing {
		if HasValidExtension(potentialPackageIndexFile) {
			return potentialPackageIndexFile, nil
		}
	}

	return nil, nil
}
