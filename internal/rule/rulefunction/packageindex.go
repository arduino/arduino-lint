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
