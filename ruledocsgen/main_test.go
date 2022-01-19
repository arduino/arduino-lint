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

package main

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/rulefunction"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath *paths.Path

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = workingDirectory.Join("testdata")
}

func TestAll(t *testing.T) {
	ruleConfigurations := []ruleconfiguration.Type{
		{
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.All,
			Category:         "structure",
			Subcategory:      "general",
			ID:               "LS001",
			Brief:            "invalid library",
			Description:      "The path does not contain a valid Arduino library.",
			MessageTemplate:  "Path does not contain a valid Arduino library.",
			Reference:        "https://arduino.github.io/arduino-cli/latest/library-specification",
			DisableModes:     nil,
			EnableModes:      []rulemode.Type{rulemode.Default},
			InfoModes:        nil,
			WarningModes:     nil,
			ErrorModes:       []rulemode.Type{rulemode.Default},
			RuleFunction:     rulefunction.LibraryInvalid,
		},
		{
			ProjectType:      projecttype.Library,
			SuperprojectType: projecttype.Library,
			Category:         "structure",
			Subcategory:      "miscellaneous",
			ID:               "LS007",
			Brief:            ".exe file",
			Description:      "A file with `.exe` file extension was found under the library folder. Presence of this file blocks addition to the Library Manager index.",
			MessageTemplate:  ".exe file(s) found. Presence of these files blocks addition to the Library Manager index:\n{{.}}",
			Reference:        "",
			DisableModes:     []rulemode.Type{rulemode.Default},
			EnableModes:      []rulemode.Type{rulemode.LibraryManagerSubmission, rulemode.LibraryManagerIndexed, rulemode.LibraryManagerIndexing},
			InfoModes:        nil,
			WarningModes:     nil,
			ErrorModes:       []rulemode.Type{rulemode.Default},
			RuleFunction:     rulefunction.LibraryHasExe,
		},
		{
			ProjectType:      projecttype.Sketch,
			SuperprojectType: projecttype.All,
			Category:         "structure",
			Subcategory:      "root folder",
			ID:               "SS001",
			Brief:            "name mismatch",
			Description:      "There is no `.ino` sketch file with name matching the sketch folder. The primary sketch file name must match the folder for the sketch to be valid.",
			MessageTemplate:  "Sketch file/folder name mismatch. The primary sketch file name must match the folder: {{.}}",
			Reference:        "https://arduino.github.io/arduino-cli/latest/sketch-specification/#primary-sketch-file",
			DisableModes:     nil,
			EnableModes:      []rulemode.Type{rulemode.Default},
			InfoModes:        nil,
			WarningModes:     []rulemode.Type{rulemode.Permissive},
			ErrorModes:       []rulemode.Type{rulemode.Default},
			RuleFunction:     rulefunction.SketchNameMismatch,
		},
		{
			ProjectType:      projecttype.Platform,
			SuperprojectType: projecttype.All,
			Category:         "configuration files",
			Subcategory:      "boards.txt",
			ID:               "PF001",
			Brief:            "boards.txt missing",
			Description:      "The `boards.txt` configuration file was not found in the platform folder",
			MessageTemplate:  "Required boards.txt is missing. Expected at: {{.}}",
			Reference:        "https://arduino.github.io/arduino-cli/latest/platform-specification/#boardstxt",
			DisableModes:     nil,
			EnableModes:      []rulemode.Type{rulemode.Default},
			InfoModes:        nil,
			WarningModes:     nil,
			ErrorModes:       []rulemode.Type{rulemode.Default},
			RuleFunction:     rulefunction.BoardsTxtMissing,
		},
		{
			ProjectType:      projecttype.Platform,
			SuperprojectType: projecttype.All,
			Category:         "configuration files",
			Subcategory:      "boards.txt",
			ID:               "PF009",
			Brief:            "use of compiler.<pattern type>.extra_flags & foo 'bar' \"baz\"",
			MessageTemplate:  "Board ID(s) {{.}} use compiler.<pattern type>.extra_flags properties. These are intended to be left for use by the user.",
			Description:      "A board definition in the platform's `boards.txt` configuration file is using one of the `compiler.<pattern type>.extra_flags` properties (e.g., `compiler.cpp.extra_flags`). These are intended to be left for use by the user as a standardized interface for customizing the compilation commands. The platform author can define as many arbitrary properties as they like, so there is no need for them to take the user's properties.",
			Reference:        "",
			DisableModes:     nil,
			EnableModes:      []rulemode.Type{rulemode.Default},
			InfoModes:        nil,
			WarningModes:     []rulemode.Type{rulemode.Default},
			ErrorModes:       []rulemode.Type{rulemode.Strict},
			RuleFunction:     rulefunction.BoardsTxtUserExtraFlagsUsage,
		},
		{
			ProjectType:      projecttype.PackageIndex,
			SuperprojectType: projecttype.All,
			Category:         "data",
			Subcategory:      "package",
			ID:               "IA004",
			Brief:            "packages[*].name missing",
			Description:      "The package index has a package without a `name` property.",
			MessageTemplate:  "Missing packages[*].name property in package(s): {{.}}",
			Reference:        "https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#json-index-file-contents",
			DisableModes:     nil,
			EnableModes:      []rulemode.Type{rulemode.Default},
			InfoModes:        nil,
			WarningModes:     nil,
			ErrorModes:       []rulemode.Type{rulemode.Default},
			RuleFunction:     rulefunction.PackageIndexPackagesNameMissing,
		},
	}

	outputPath, err := paths.MkTempDir("", "backup-test-testall")
	require.NoError(t, err)
	defer outputPath.RemoveAll()

	generateRulesDocumentation(ruleConfigurations, outputPath)

	assert.True(t, outputPath.Exist())

	for _, outputFileName := range []string{"sketch.md", "library.md", "platform.md", "package-index.md"} {
		assert.True(t, outputPath.Join(outputFileName).Exist(), outputFileName)
		rules, err := outputPath.Join(outputFileName).ReadFileAsLines()
		require.NoError(t, err)
		goldenRules, err := testDataPath.Join("golden", outputFileName).ReadFileAsLines()
		assert.Equal(t, goldenRules, rules, outputFileName)
	}
}
