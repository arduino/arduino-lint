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

// PackageIndexPackagesNameLTMinLength checks for incorrect type of the packages[].name property.
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

// PackageIndexPackagesMaintainerLTMinLength checks for incorrect type of the packages[].maintainer property.
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
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
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
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
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
		return ruleresult.Fail, strings.Join(nonCompliantIDs, ", ")
	}

	return ruleresult.Pass, ""
}
