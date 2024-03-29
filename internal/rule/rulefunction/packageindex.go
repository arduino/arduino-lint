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

package rulefunction

import (
	"strings"

	"github.com/arduino/arduino-lint/internal/project/packageindex"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
)

// The rule functions for package indexes.

// PackageIndexMissing checks whether a file resembling a package index was found in the specified project folder.
func PackageIndexMissing() (result ruleresult.Type, output string) {
	if projectdata.ProjectPath() == nil {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PackageIndexFilenameInvalid checks whether the package index's filename is valid for 3rd party projects.
func PackageIndexFilenameInvalid() (result ruleresult.Type, output string) {
	if projectdata.ProjectPath() == nil {
		return ruleresult.NotRun, "Package index not found"
	}

	if packageindex.HasValidFilename(projectdata.ProjectPath(), false) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.ProjectPath().Base()
}

// PackageIndexOfficialFilenameInvalid checks whether the package index's filename is valid for official projects.
func PackageIndexOfficialFilenameInvalid() (result ruleresult.Type, output string) {
	if projectdata.ProjectPath() == nil {
		return ruleresult.NotRun, "Package index not found"
	}

	if packageindex.HasValidFilename(projectdata.ProjectPath(), true) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, projectdata.ProjectPath().Base()
}

// PackageIndexJSONFormat checks whether the package index file is a valid JSON document.
func PackageIndexJSONFormat() (result ruleresult.Type, output string) {
	if projectdata.ProjectPath() == nil {
		return ruleresult.NotRun, "Package index not found"
	}

	if isValidJSON(projectdata.ProjectPath()) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, ""
}

// PackageIndexFormat checks for invalid package index data format.
func PackageIndexFormat() (result ruleresult.Type, output string) {
	if projectdata.ProjectPath() == nil {
		return ruleresult.NotRun, "Package index not found"
	}

	if projectdata.PackageIndexCLILoadError() != nil {
		return ruleresult.Fail, projectdata.PackageIndexCLILoadError().Error()
	}

	return ruleresult.Pass, ""
}

// PackageIndexAdditionalProperties checks for additional properties in the package index root.
func PackageIndexAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	if schema.ProhibitedAdditionalProperties("", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesMissing checks for missing packages property.
func PackageIndexPackagesMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	if schema.RequiredPropertyMissing("/packages", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesIncorrectType checks for incorrect type of packages[].
func PackageIndexPackagesIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	if schema.PropertyTypeMismatch("/packages", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
		return ruleresult.Fail, ""
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesAdditionalProperties checks for additional properties in packages[].
func PackageIndexPackagesAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.ProhibitedAdditionalProperties(packageData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesNameMissing checks for missing packages[].name property.
func PackageIndexPackagesNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesNameIncorrectType checks for incorrect type of the packages[].name property.
func PackageIndexPackagesNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesNameLTMinLength checks for packages[].name property less than the minimum length.
func PackageIndexPackagesNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyLessThanMinLength(packageData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesNameIsArduino checks for packages[].name being "arduino".
func PackageIndexPackagesNameIsArduino() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.ValidationErrorMatch(
			"^#"+packageData.JSONPointer+"/name$",
			"/patternObjects/notArduino",
			"",
			"",
			projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification],
		) {
			// Since the package name is implicit in the rule itself, it makes most sense to use the JSON pointer to identify.
			nonCompliantIDs = append(nonCompliantIDs, packageData.JSONPointer)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesMaintainerMissing checks for missing packages[].maintainer property.
func PackageIndexPackagesMaintainerMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/maintainer", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesMaintainerIncorrectType checks for incorrect type of the packages[].maintainer property.
func PackageIndexPackagesMaintainerIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/maintainer", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesMaintainerLTMinLength checks for packages[].maintainer property less than the minimum length.
func PackageIndexPackagesMaintainerLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyLessThanMinLength(packageData.JSONPointer+"/maintainer", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesMaintainerStartsWithArduino checks for packages[].maintainer starting with "arduino".
func PackageIndexPackagesMaintainerStartsWithArduino() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.ValidationErrorMatch(
			"^#"+packageData.JSONPointer+"/maintainer$",
			"/patternObjects/notStartsWithArduino",
			"",
			"",
			projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict],
		) {
			// Since the package name is implicit in the rule itself, it makes most sense to use the JSON pointer to identify.
			nonCompliantIDs = append(nonCompliantIDs, packageData.JSONPointer)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesWebsiteURLMissing checks for missing packages[].websiteURL property.
func PackageIndexPackagesWebsiteURLMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/websiteURL", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesWebsiteURLIncorrectType checks for incorrect type of the packages[].websiteURL property.
func PackageIndexPackagesWebsiteURLIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/websiteURL", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesWebsiteURLInvalidFormat checks for incorrect format of the packages[].websiteURL property.
func PackageIndexPackagesWebsiteURLInvalidFormat() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyFormatMismatch(packageData.JSONPointer+"/websiteURL", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesWebsiteURLDeadLink checks for dead links in packages[].websiteURL.
func PackageIndexPackagesWebsiteURLDeadLink() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, data := range projectdata.PackageIndexPackages() {
		url, ok := data.Object["websiteURL"].(string)
		if !ok {
			continue
		}

		if url == "" {
			continue
		}

		if checkURL(url) == nil {
			continue
		}

		nonCompliantIDs = append(nonCompliantIDs, data.ID)
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesEmailMissing checks for missing packages[].email property.
func PackageIndexPackagesEmailMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/email", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesEmailIncorrectType checks for incorrect type of the packages[].email property.
func PackageIndexPackagesEmailIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/email", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpIncorrectType checks for incorrect type of the packages[].help property.
func PackageIndexPackagesHelpIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/help", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpAdditionalProperties checks for additional properties in packages[].help.
func PackageIndexPackagesHelpAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.ProhibitedAdditionalProperties(packageData.JSONPointer+"/help", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpOnlineMissing checks for missing packages[].help.online property.
func PackageIndexPackagesHelpOnlineMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpOnlineIncorrectType checks for incorrect type of the packages[].help.online property.
func PackageIndexPackagesHelpOnlineIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpOnlineInvalidFormat checks for incorrect format of the packages[].help.online property.
func PackageIndexPackagesHelpOnlineInvalidFormat() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyFormatMismatch(packageData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesHelpOnlineDeadLink checks for dead links in packages[].help.online.
func PackageIndexPackagesHelpOnlineDeadLink() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, data := range projectdata.PackageIndexPackages() {
		help, ok := data.Object["help"].(map[string]interface{})
		if !ok {
			continue
		}

		url, ok := help["online"].(string)
		if !ok {
			continue
		}

		if url == "" {
			continue
		}

		if checkURL(url) == nil {
			continue
		}

		nonCompliantIDs = append(nonCompliantIDs, data.ID)
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMissing checks for missing packages[].platforms[] property.
func PackageIndexPackagesPlatformsMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/platforms", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsIncorrectType checks for incorrect type of packages[].platforms.
func PackageIndexPackagesPlatformsIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/platforms", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsAdditionalProperties checks for additional properties in packages[].platforms[].
func PackageIndexPackagesPlatformsAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.ProhibitedAdditionalProperties(platformData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsNameMissing checks for missing packages[].platforms[].name property.
func PackageIndexPackagesPlatformsNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsNameIncorrectType checks for incorrect type of the packages[].platforms[].name property.
func PackageIndexPackagesPlatformsNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsNameLTMinLength checks for packages[].platforms[].name property less than the minimum length.
func PackageIndexPackagesPlatformsNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyLessThanMinLength(platformData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchitectureMissing checks for missing packages[].platforms[].architecture property.
func PackageIndexPackagesPlatformsArchitectureMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/architecture", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchitectureIncorrectType checks for incorrect type of the packages[].platforms[].architecture property.
func PackageIndexPackagesPlatformsArchitectureIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/architecture", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchitectureLTMinLength checks for packages[].platforms[].architecture property less than the minimum length.
func PackageIndexPackagesPlatformsArchitectureLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyLessThanMinLength(platformData.JSONPointer+"/architecture", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsVersionMissing checks for missing packages[].platforms[].version property.
func PackageIndexPackagesPlatformsVersionMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsVersionIncorrectType checks for incorrect type of the packages[].platforms[].version property.
func PackageIndexPackagesPlatformsVersionIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsVersionNonRelaxedSemver checks whether the packages[].platforms[].version property is "relaxed semver" compliant.
func PackageIndexPackagesPlatformsVersionNonRelaxedSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyPatternMismatch(platformData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsVersionNonSemver checks whether the packages[].platforms[].version property is semver compliant.
func PackageIndexPackagesPlatformsVersionNonSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyPatternMismatch(platformData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDeprecatedIncorrectType checks for incorrect type of the packages[].platforms[].deprecated property.
func PackageIndexPackagesPlatformsDeprecatedIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/deprecated", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsCategoryMissing checks for missing packages[].platforms[].category property.
func PackageIndexPackagesPlatformsCategoryMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/category", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsCategoryIncorrectType checks for incorrect type of the packages[].platforms[].category property.
func PackageIndexPackagesPlatformsCategoryIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/category", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsCategoryThirdPartyInvalid checks for invalid value of the packages[].platforms[].category property for 3rd party platforms.
func PackageIndexPackagesPlatformsCategoryThirdPartyInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyEnumMismatch(platformData.JSONPointer+"/category", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpMissing checks for missing packages[].platforms[].help property.
func PackageIndexPackagesPlatformsHelpMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/help", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpIncorrectType checks for incorrect type of the packages[].platforms[].help property.
func PackageIndexPackagesPlatformsHelpIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/help", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpAdditionalProperties checks for additional properties in packages[].help.
func PackageIndexPackagesPlatformsHelpAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.ProhibitedAdditionalProperties(platformData.JSONPointer+"/help", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpOnlineMissing checks for missing packages[].platforms[].help.online property.
func PackageIndexPackagesPlatformsHelpOnlineMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpOnlineIncorrectType checks for incorrect type of the packages[].platforms[].help.online property.
func PackageIndexPackagesPlatformsHelpOnlineIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpOnlineInvalidFormat checks for incorrect format of the packages[].platforms[].help.online property.
func PackageIndexPackagesPlatformsHelpOnlineInvalidFormat() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyFormatMismatch(platformData.JSONPointer+"/help/online", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsHelpOnlineDeadLink checks for dead links in packages[].platforms[].help.online.
func PackageIndexPackagesPlatformsHelpOnlineDeadLink() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, data := range projectdata.PackageIndexPlatforms() {
		help, ok := data.Object["help"].(map[string]interface{})
		if !ok {
			continue
		}

		url, ok := help["online"].(string)
		if !ok {
			continue
		}

		if url == "" {
			continue
		}

		if checkURL(url) == nil {
			continue
		}

		nonCompliantIDs = append(nonCompliantIDs, data.ID)
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsURLMissing checks for missing packages[].platforms[].url property.
func PackageIndexPackagesPlatformsURLMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsURLIncorrectType checks for incorrect type of the packages[].platforms[].url property.
func PackageIndexPackagesPlatformsURLIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsURLInvalidFormat checks for incorrect format of the packages[].platforms[].url property.
func PackageIndexPackagesPlatformsURLInvalidFormat() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyFormatMismatch(platformData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsURLDeadLink checks for dead links in packages[].platforms[].url.
func PackageIndexPackagesPlatformsURLDeadLink() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, data := range projectdata.PackageIndexPlatforms() {
		url, ok := data.Object["url"].(string)
		if !ok {
			continue
		}

		if url == "" {
			continue
		}

		if checkURL(url) == nil {
			continue
		}

		nonCompliantIDs = append(nonCompliantIDs, data.ID)
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchiveFileNameMissing checks for missing packages[].platforms[].archiveFileName property.
func PackageIndexPackagesPlatformsArchiveFileNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchiveFileNameIncorrectType checks for incorrect type of the packages[].platforms[].archiveFileName property.
func PackageIndexPackagesPlatformsArchiveFileNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchiveFileNameLTMinLength checks for packages[].platforms[].archiveFileName property less than the minimum length.
func PackageIndexPackagesPlatformsArchiveFileNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyLessThanMinLength(platformData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsArchiveFileNameInvalid checks for invalid format of packages[].platforms[].archiveFileName property.
func PackageIndexPackagesPlatformsArchiveFileNameInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyPatternMismatch(platformData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsChecksumMissing checks for missing packages[].platforms[].checksum property.
func PackageIndexPackagesPlatformsChecksumMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsChecksumIncorrectType checks for incorrect type of the packages[].platforms[].checksum property.
func PackageIndexPackagesPlatformsChecksumIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsChecksumInvalid checks for invalid format of packages[].platforms[].checksum property.
func PackageIndexPackagesPlatformsChecksumInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyPatternMismatch(platformData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsChecksumDiscouragedAlgorithm checks for use of discouraged hash algorithm in packages[].platforms[].checksum property.
func PackageIndexPackagesPlatformsChecksumDiscouragedAlgorithm() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.ValidationErrorMatch("^#"+platformData.JSONPointer+"/checksum$", "/patternObjects/usesSHA256", "", "", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsSizeMissing checks for missing packages[].platforms[].size property.
func PackageIndexPackagesPlatformsSizeMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsSizeIncorrectType checks for incorrect type of the packages[].platforms[].size property.
func PackageIndexPackagesPlatformsSizeIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsSizeInvalid checks for invalid format of packages[].platforms[].size property.
func PackageIndexPackagesPlatformsSizeInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyPatternMismatch(platformData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsMissing checks for missing packages[].platforms[].boards[] property.
func PackageIndexPackagesPlatformsBoardsMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/boards", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsIncorrectType checks for incorrect type of the packages[].platforms[].boards property.
func PackageIndexPackagesPlatformsBoardsIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/boards", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsAdditionalProperties checks for additional properties in packages[].platforms[].boards[].
func PackageIndexPackagesPlatformsBoardsAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, boardData := range projectdata.PackageIndexBoards() {
		if schema.ProhibitedAdditionalProperties(boardData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, boardData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsNameMissing checks for missing packages[].platforms[].boards[].name property.
func PackageIndexPackagesPlatformsBoardsNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, boardData := range projectdata.PackageIndexBoards() {
		if schema.RequiredPropertyMissing(boardData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, boardData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsNameIncorrectType checks for incorrect type of the packages[].platforms[].boards[].name property.
func PackageIndexPackagesPlatformsBoardsNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, boardData := range projectdata.PackageIndexBoards() {
		if schema.PropertyTypeMismatch(boardData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, boardData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsBoardsNameLTMinLength checks for packages[].platforms[].board[].name property less than the minimum length.
func PackageIndexPackagesPlatformsBoardsNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, boardData := range projectdata.PackageIndexBoards() {
		if schema.PropertyLessThanMinLength(boardData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, boardData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesMissing checks for missing packages[].platforms[].toolsDependencies[] property.
func PackageIndexPackagesPlatformsToolsDependenciesMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.RequiredPropertyMissing(platformData.JSONPointer+"/toolsDependencies", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesIncorrectType checks for incorrect type of the packages[].platforms[].toolsDependencies property.
func PackageIndexPackagesPlatformsToolsDependenciesIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/toolsDependencies", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesAdditionalProperties checks for additional properties in packages[].platforms[].toolsDependencies[].
func PackageIndexPackagesPlatformsToolsDependenciesAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.ProhibitedAdditionalProperties(dependencyData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesPackagerMissing checks for missing packages[].platforms[].toolsDependencies[].packager property.
func PackageIndexPackagesPlatformsToolsDependenciesPackagerMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesPackagerIncorrectType checks for incorrect type of the packages[].platforms[].toolsDependencies[].packager property.
func PackageIndexPackagesPlatformsToolsDependenciesPackagerIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesPackagerLTMinLength checks for packages[].platforms[].toolsDependencies[].packager property less than the minimum length.
func PackageIndexPackagesPlatformsToolsDependenciesPackagerLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesNameMissing checks for missing packages[].platforms[].toolsDependencies[].name property.
func PackageIndexPackagesPlatformsToolsDependenciesNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesNameIncorrectType checks for incorrect type of the packages[].platforms[].toolsDependencies[].name property.
func PackageIndexPackagesPlatformsToolsDependenciesNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesNameLTMinLength checks for packages[].platforms[].toolsDependencies[].name property less than the minimum length.
func PackageIndexPackagesPlatformsToolsDependenciesNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesVersionMissing checks for missing packages[].platforms[].toolsDependencies[].version property.
func PackageIndexPackagesPlatformsToolsDependenciesVersionMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesVersionIncorrectType checks for incorrect type of the packages[].platforms[].toolsDependencies[].packager property.
func PackageIndexPackagesPlatformsToolsDependenciesVersionIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesVersionNonRelaxedSemver checks whether the packages[].platforms[].toolsDependencies[].version property is "relaxed semver" compliant.
func PackageIndexPackagesPlatformsToolsDependenciesVersionNonRelaxedSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyPatternMismatch(dependencyData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsToolsDependenciesVersionNonSemver checks whether the packages[].platforms[].toolsDependencies[].version property is semver compliant.
func PackageIndexPackagesPlatformsToolsDependenciesVersionNonSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexToolsDependencies() {
		if schema.PropertyPatternMismatch(dependencyData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesIncorrectType checks for incorrect type of the packages[].platforms[].discoveryDependencies property.
func PackageIndexPackagesPlatformsDiscoveryDependenciesIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/discoveryDependencies", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesAdditionalProperties checks for additional properties in packages[].platforms[].discoveryDependencies[].
func PackageIndexPackagesPlatformsDiscoveryDependenciesAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.ProhibitedAdditionalProperties(dependencyData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerMissing checks for missing packages[].platforms[].discoveryDependencies[].packager property.
func PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerIncorrectType checks for incorrect type of the packages[].platforms[].discoveryDependencies[].packager property.
func PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerLTMinLength checks for packages[].platforms[].discoveryDependencies[].packager property less than the minimum length.
func PackageIndexPackagesPlatformsDiscoveryDependenciesPackagerLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesNameMissing checks for missing packages[].platforms[].discoveryDependencies[].name property.
func PackageIndexPackagesPlatformsDiscoveryDependenciesNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesNameIncorrectType checks for incorrect type of the packages[].platforms[].discoveryDependencies[].name property.
func PackageIndexPackagesPlatformsDiscoveryDependenciesNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsDiscoveryDependenciesNameLTMinLength checks for packages[].platforms[].discoveryDependencies[].name property less than the minimum length.
func PackageIndexPackagesPlatformsDiscoveryDependenciesNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexDiscoveryDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesIncorrectType checks for incorrect type of the packages[].platforms[].monitorDependencies property.
func PackageIndexPackagesPlatformsMonitorDependenciesIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, platformData := range projectdata.PackageIndexPlatforms() {
		if schema.PropertyTypeMismatch(platformData.JSONPointer+"/monitorDependencies", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, platformData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesAdditionalProperties checks for additional properties in packages[].platforms[].monitorDependencies[].
func PackageIndexPackagesPlatformsMonitorDependenciesAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.ProhibitedAdditionalProperties(dependencyData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesPackagerMissing checks for missing packages[].platforms[].monitorDependencies[].packager property.
func PackageIndexPackagesPlatformsMonitorDependenciesPackagerMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesPackagerIncorrectType checks for incorrect type of the packages[].platforms[].monitorDependencies[].packager property.
func PackageIndexPackagesPlatformsMonitorDependenciesPackagerIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesPackagerLTMinLength checks for packages[].platforms[].monitorDependencies[].packager property less than the minimum length.
func PackageIndexPackagesPlatformsMonitorDependenciesPackagerLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/packager", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesNameMissing checks for missing packages[].platforms[].monitorDependencies[].name property.
func PackageIndexPackagesPlatformsMonitorDependenciesNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.RequiredPropertyMissing(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesNameIncorrectType checks for incorrect type of the packages[].platforms[].monitorDependencies[].name property.
func PackageIndexPackagesPlatformsMonitorDependenciesNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.PropertyTypeMismatch(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesPlatformsMonitorDependenciesNameLTMinLength checks for packages[].platforms[].monitorDependencies[].name property less than the minimum length.
func PackageIndexPackagesPlatformsMonitorDependenciesNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, dependencyData := range projectdata.PackageIndexMonitorDependencies() {
		if schema.PropertyLessThanMinLength(dependencyData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, dependencyData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsMissing checks for missing packages[].tools property.
func PackageIndexPackagesToolsMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.RequiredPropertyMissing(packageData.JSONPointer+"/tools", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsIncorrectType checks for incorrect type of packages[].tools.
func PackageIndexPackagesToolsIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, packageData := range projectdata.PackageIndexPackages() {
		if schema.PropertyTypeMismatch(packageData.JSONPointer+"/tools", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, packageData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsAdditionalProperties checks for additional properties in packages[].tools[].
func PackageIndexPackagesToolsAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.ProhibitedAdditionalProperties(toolData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsNameMissing checks for missing packages[].tools[].name property.
func PackageIndexPackagesToolsNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.RequiredPropertyMissing(toolData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsNameIncorrectType checks for incorrect type of the packages[].tools[].name property.
func PackageIndexPackagesToolsNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyTypeMismatch(toolData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsNameLTMinLength checks for packages[].tools[].name property less than the minimum length.
func PackageIndexPackagesToolsNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyLessThanMinLength(toolData.JSONPointer+"/name", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsVersionMissing checks for missing packages[].tools[].version property.
func PackageIndexPackagesToolsVersionMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.RequiredPropertyMissing(toolData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsVersionIncorrectType checks for incorrect type of the packages[].tools[].version property.
func PackageIndexPackagesToolsVersionIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyTypeMismatch(toolData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsVersionNonRelaxedSemver checks whether the packages[].tools[].version property is "relaxed semver" compliant.
func PackageIndexPackagesToolsVersionNonRelaxedSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyPatternMismatch(toolData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsVersionNonSemver checks whether the packages[].tools[].version property is semver compliant.
func PackageIndexPackagesToolsVersionNonSemver() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyPatternMismatch(toolData.JSONPointer+"/version", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsMissing checks for missing packages[].tools[].systems[] property.
func PackageIndexPackagesToolsSystemsMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.RequiredPropertyMissing(toolData.JSONPointer+"/systems", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsIncorrectType checks for incorrect type of the packages[].tools[].systems property.
func PackageIndexPackagesToolsSystemsIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, toolData := range projectdata.PackageIndexTools() {
		if schema.PropertyTypeMismatch(toolData.JSONPointer+"/systems", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, toolData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsAdditionalProperties checks for additional properties in packages[].tools[].systems[].
func PackageIndexPackagesToolsSystemsAdditionalProperties() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.ProhibitedAdditionalProperties(systemData.JSONPointer, projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsHostMissing checks for missing packages[].tools[].systems[].host property.
func PackageIndexPackagesToolsSystemsHostMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.RequiredPropertyMissing(systemData.JSONPointer+"/host", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsHostIncorrectType checks for incorrect type of the packages[].tools[].systems[].host property.
func PackageIndexPackagesToolsSystemsHostIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyTypeMismatch(systemData.JSONPointer+"/host", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsHostInvalid checks for invalid format of whether the packages[].tools[].systems[].host property.
func PackageIndexPackagesToolsSystemsHostInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyPatternMismatch(systemData.JSONPointer+"/host", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsURLMissing checks for missing packages[].tools[].systems[].url property.
func PackageIndexPackagesToolsSystemsURLMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.RequiredPropertyMissing(systemData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsURLIncorrectType checks for incorrect type of the packages[].tools[].systems[].url property.
func PackageIndexPackagesToolsSystemsURLIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyTypeMismatch(systemData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsURLInvalidFormat checks for incorrect format of the packages[].tools[].systems[].url property.
func PackageIndexPackagesToolsSystemsURLInvalidFormat() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyFormatMismatch(systemData.JSONPointer+"/url", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsURLDeadLink checks for dead links in packages[].tools[].systems[].url.
func PackageIndexPackagesToolsSystemsURLDeadLink() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, data := range projectdata.PackageIndexSystems() {
		url, ok := data.Object["url"].(string)
		if !ok {
			continue
		}

		if url == "" {
			continue
		}

		if checkURL(url) == nil {
			continue
		}

		nonCompliantIDs = append(nonCompliantIDs, data.ID)
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsArchiveFileNameMissing checks for missing packages[].tools[].systems[].archiveFileName property.
func PackageIndexPackagesToolsSystemsArchiveFileNameMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.RequiredPropertyMissing(systemData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsArchiveFileNameIncorrectType checks for incorrect type of the packages[].tools[].systems[].archiveFileName property.
func PackageIndexPackagesToolsSystemsArchiveFileNameIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyTypeMismatch(systemData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsArchiveFileNameLTMinLength checks for packages[].tools[].systems[].archiveFileName property less than the minimum length.
func PackageIndexPackagesToolsSystemsArchiveFileNameLTMinLength() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyLessThanMinLength(systemData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsArchiveFileNameInvalid checks for invalid format of packages[].tools[].systems[].archiveFileName property.
func PackageIndexPackagesToolsSystemsArchiveFileNameInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyPatternMismatch(systemData.JSONPointer+"/archiveFileName", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsChecksumMissing checks for missing packages[].tools[].systems[].checksum property.
func PackageIndexPackagesToolsSystemsChecksumMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.RequiredPropertyMissing(systemData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsChecksumIncorrectType checks for incorrect type of the packages[].tools[].systems[].checksum property.
func PackageIndexPackagesToolsSystemsChecksumIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyTypeMismatch(systemData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsChecksumInvalid checks for invalid format of packages[].tools[].systems[].checksum property.
func PackageIndexPackagesToolsSystemsChecksumInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyPatternMismatch(systemData.JSONPointer+"/checksum", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsChecksumDiscouragedAlgorithm checks for use of discouraged hash algorithm in packages[].tools[].systems[].checksum property.
func PackageIndexPackagesToolsSystemsChecksumDiscouragedAlgorithm() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.ValidationErrorMatch("^#"+systemData.JSONPointer+"/checksum$", "/patternObjects/usesSHA256", "", "", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Strict]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsSizeMissing checks for missing packages[].tools[].systems[].size property.
func PackageIndexPackagesToolsSystemsSizeMissing() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.RequiredPropertyMissing(systemData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsSizeIncorrectType checks for incorrect type of the packages[].tools[].systems[].size property.
func PackageIndexPackagesToolsSystemsSizeIncorrectType() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyTypeMismatch(systemData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}

// PackageIndexPackagesToolsSystemsSizeInvalid checks for invalid format of packages[].tools[].systems[].size property.
func PackageIndexPackagesToolsSystemsSizeInvalid() (result ruleresult.Type, output string) {
	if projectdata.PackageIndexLoadError() != nil {
		return ruleresult.NotRun, "Error loading package index"
	}

	nonCompliantIDs := []string{}
	for _, systemData := range projectdata.PackageIndexSystems() {
		if schema.PropertyPatternMismatch(systemData.JSONPointer+"/size", projectdata.PackageIndexSchemaValidationResult()[compliancelevel.Specification]) {
			nonCompliantIDs = append(nonCompliantIDs, systemData.ID)
		}
	}

	if len(nonCompliantIDs) > 0 {
		return ruleresult.Fail, brokenOutputList(nonCompliantIDs)
	}

	return ruleresult.Pass, ""
}
