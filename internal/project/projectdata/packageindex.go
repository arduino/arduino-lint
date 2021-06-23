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

package projectdata

import (
	"bytes"
	"fmt"
	"text/template"

	clipackageindex "github.com/arduino/arduino-cli/arduino/cores/packageindex"
	"github.com/arduino/arduino-lint/internal/project/packageindex"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
)

// PackageIndexData is the type for package index data.
type PackageIndexData struct {
	ID          string                 // Identifier for display to humans
	JSONPointer string                 // Path to the data in the JSON document
	Object      map[string]interface{} // The data of the object
}

// InitializeForPackageIndex gathers the package index rule data for the specified project.
func InitializeForPackageIndex() {
	packageIndex, packageIndexLoadError = packageindex.Properties(ProjectPath())
	if ProjectPath() != nil {
		_, packageIndexCLILoadError = clipackageindex.LoadIndex(ProjectPath())
	}

	packageIndexPackages = nil
	packageIndexPlatforms = nil
	packageIndexBoards = nil
	packageIndexTools = nil
	packageIndexSystems = nil
	packageIndexSchemaValidationResult = nil
	if packageIndexLoadError == nil {
		packageIndexPackages = getPackageIndexData(PackageIndex(), "", "packages", "", "{{index . 0}}", []string{"name"})

		for _, packageData := range PackageIndexPackages() {
			packageIndexPlatforms = append(packageIndexPlatforms, getPackageIndexData(packageData.Object, packageData.JSONPointer, "platforms", packageData.ID, ":{{index . 0}}@{{index . 1}}", []string{"architecture", "version"})...)
		}

		for _, platformData := range PackageIndexPlatforms() {
			packageIndexBoards = append(packageIndexBoards, getPackageIndexData(platformData.Object, platformData.JSONPointer, "boards", platformData.ID, " - {{index . 0}}", []string{"name"})...)
		}

		for _, packageData := range PackageIndexPackages() {
			packageIndexTools = append(packageIndexTools, getPackageIndexData(packageData.Object, packageData.JSONPointer, "tools", packageData.ID, ":{{index . 0}}@{{index . 1}}", []string{"name", "version"})...)
		}

		for _, toolData := range PackageIndexTools() {
			packageIndexSystems = append(packageIndexSystems, getPackageIndexData(toolData.Object, toolData.JSONPointer, "systems", toolData.ID, " - {{index . 0}}", []string{"host"})...)
		}

		packageIndexSchemaValidationResult = packageindex.Validate(PackageIndex())
	}
}

var packageIndex map[string]interface{}

// PackageIndex returns the package index data.
func PackageIndex() map[string]interface{} {
	return packageIndex
}

var packageIndexLoadError error

// PackageIndexLoadError returns the error from loading the package index.
func PackageIndexLoadError() error {
	return packageIndexLoadError
}

var packageIndexCLILoadError error

// PackageIndexCLILoadError returns the error return of Arduino CLI's packageindex.LoadIndex().
func PackageIndexCLILoadError() error {
	return packageIndexCLILoadError
}

var packageIndexPackages []PackageIndexData

// PackageIndexPackages returns the slice of package data for the package index.
func PackageIndexPackages() []PackageIndexData {
	return packageIndexPackages
}

var packageIndexPlatforms []PackageIndexData

// PackageIndexPlatforms returns the slice of platform data for the package index.
func PackageIndexPlatforms() []PackageIndexData {
	return packageIndexPlatforms
}

var packageIndexBoards []PackageIndexData

// PackageIndexBoards returns the slice of board data for the package index.
func PackageIndexBoards() []PackageIndexData {
	return packageIndexBoards
}

var packageIndexTools []PackageIndexData

// PackageIndexTools returns the slice of tool data for the package index.
func PackageIndexTools() []PackageIndexData {
	return packageIndexTools
}

var packageIndexSystems []PackageIndexData

// PackageIndexSystems returns the slice of system data for the package index.
func PackageIndexSystems() []PackageIndexData {
	return packageIndexSystems
}

var packageIndexSchemaValidationResult map[compliancelevel.Type]schema.ValidationResult

// PackageIndexSchemaValidationResult returns the result of validating the package index against the JSON schema.
func PackageIndexSchemaValidationResult() map[compliancelevel.Type]schema.ValidationResult {
	return packageIndexSchemaValidationResult
}

func getPackageIndexData(interfaceObject map[string]interface{}, pointerPrefix string, dataKey string, iDPrefix string, iDSuffixTemplateString string, iDSuffixKeys []string) []PackageIndexData {
	var data []PackageIndexData

	interfaceSlice, ok := interfaceObject[dataKey].([]interface{})
	if !ok {
		return data
	}

	for index, interfaceElement := range interfaceSlice {
		interfaceElementData := PackageIndexData{
			JSONPointer: fmt.Sprintf("%s/%s/%v", pointerPrefix, dataKey, index),
			Object:      nil,
		}

		object, ok := interfaceElement.(map[string]interface{})
		if ok {
			interfaceElementData.Object = object
		}

		objectID := func() string {
			// In the event missing data prevents creating a standard reference ID for the data, use the JSON pointer.
			fallbackID := interfaceElementData.JSONPointer

			if iDPrefix != "" && iDPrefix == pointerPrefix {
				// Parent object uses fallback ID, so this one must even if it was possible to generate a true suffix.
				return fallbackID
			}

			// Gather the ID suffix components.
			iDSuffixComponents := []string{}
			for _, key := range iDSuffixKeys {
				component, ok := object[key].(string)
				if !ok {
					return fallbackID
				}
				if component == "" {
					return fallbackID
				}
				iDSuffixComponents = append(iDSuffixComponents, component)
			}

			// Fill the ID suffix components into the template.
			iDSuffixTemplate := template.Must(template.New("iDSuffixTemplate").Parse(iDSuffixTemplateString))
			iDSuffix := new(bytes.Buffer)
			iDSuffixTemplate.Execute(iDSuffix, iDSuffixComponents)

			return iDPrefix + iDSuffix.String()
		}
		interfaceElementData.ID = objectID()

		data = append(data, interfaceElementData)
	}

	return data
}
