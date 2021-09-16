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

// Package main generates Markdown documentation for Arduino Lint's rules.
package main

import (
	"bytes"
	"os"
	"text/template"

	"github.com/JohannesKaufmann/html-to-markdown/escape"
	"github.com/arduino/arduino-lint/internal/cli"
	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/rulelevel"
	"github.com/arduino/go-paths-helper"
	"github.com/olekukonko/tablewriter"
)

func main() {
	if len(os.Args) < 2 {
		print("error: Please provide the output folder argument")
		os.Exit(1)
	}
	outputPath := paths.New(os.Args[1])

	generateRulesDocumentation(ruleconfiguration.Configurations(), outputPath)
}

// generateRulesDocumentation generates documentation in Markdown language for the rules defined by the provided
// configurations and writes them to a file for each project type at the specified path.
func generateRulesDocumentation(ruleConfigurations []ruleconfiguration.Type, outputPath *paths.Path) {
	projectTypeReferences := map[projecttype.Type]string{
		projecttype.Sketch:       "https://arduino.github.io/arduino-cli/latest/sketch-specification/",
		projecttype.Library:      "https://arduino.github.io/arduino-cli/latest/library-specification/",
		projecttype.Platform:     "https://arduino.github.io/arduino-cli/latest/platform-specification/",
		projecttype.PackageIndex: "https://arduino.github.io/arduino-cli/latest/package_index_json-specification/",
	}

	templateFunctions := template.FuncMap{
		// Some the rule config text is intended for use in both tool output and in the reference, so can't be formatted at
		// the source as Markdown. Incidental markup characters in that text must be escaped.
		"escape": escape.MarkdownCharacters,
	}

	projectRulesIntroTemplate := template.Must(template.New("messageTemplate").Parse(
		"Arduino Lint provides {{.RuleCount}} rules for the [`{{.ProjectType}}`]({{.ProjectTypeReference}}) project type:\n",
	))
	ruleDocumentationTemplate := template.Must(template.New("messageTemplate").Funcs(templateFunctions).Parse(`
---

<a id="{{.ID}}"></a>

## {{escape .Brief}} (` + "`" + `{{.ID}}` + "`" + `)

{{.Description}}

{{if .Reference}}More information: [**here**]({{.Reference}})<br />{{end}}
Enabled for superproject type: {{.SuperprojectType}}<br />
Category: {{.Category}}<br />
Subcategory: {{.Subcategory}}

##### Rule levels

`))

	rulesDocumentation := make(
		map[projecttype.Type]struct {
			content bytes.Buffer
			count   int
		},
	)

	// Generate the rule documentation, indexed by project type.
	for _, ruleConfiguration := range ruleConfigurations {
		projectRulesData := rulesDocumentation[ruleConfiguration.ProjectType]
		projectRulesData.count++

		// Fill the template with the rule's configuration data.
		ruleDocumentationTemplate.Execute(&projectRulesData.content, ruleConfiguration)

		// Generate a table of the rule violation levels.
		// This is too complex to handle via the template so it is appended to the templated text.
		levelsData := ruleLevels(ruleConfiguration)
		var table bytes.Buffer
		tableWriter := tablewriter.NewWriter(&table)
		tableWriter.SetAutoFormatHeaders(false)
		tableWriter.SetHeader(levelsData[0])
		tableWriter.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		tableWriter.SetCenterSeparator("|")
		tableWriter.AppendBulk(levelsData[1:])
		tableWriter.Render()
		if _, err := table.WriteTo(&projectRulesData.content); err != nil {
			panic(err)
		}

		rulesDocumentation[ruleConfiguration.ProjectType] = projectRulesData
	}

	// Write the rule documentation to file, one for each project type.
	outputPath.MkdirAll()
	for projectType, projectRules := range rulesDocumentation {
		projectPage := new(bytes.Buffer)
		projectRulesIntroTemplate.Execute(projectPage, struct {
			RuleCount            int
			ProjectType          string
			ProjectTypeReference string
		}{
			RuleCount:            projectRules.count,
			ProjectType:          projectType.String(),
			ProjectTypeReference: projectTypeReferences[projectType],
		})
		projectRules.content.WriteTo(projectPage)

		outputFile := outputPath.Join(projectType.String() + ".md")
		if err := outputFile.WriteFile(projectPage.Bytes()); err != nil {
			panic(err)
		}
	}
}

// ruleLevels returns the level of a rule violation for each of the relevant Arduino Lint configurations.
func ruleLevels(ruleConfiguration ruleconfiguration.Type) [][]string {
	complianceModes := []rulemode.Type{
		rulemode.Permissive,
		rulemode.Specification,
		rulemode.Strict,
	}

	libraryManagerModes := []rulemode.Type{
		rulemode.LibraryManagerSubmission,
		rulemode.LibraryManagerIndexed,
	}

	// `--library-manager` flag values are defined separately from modes due to the need to also document levels for the "false" value.
	libraryManagerFlagValues := []string{
		rulemode.LibraryManagerSubmission.String(),
		rulemode.LibraryManagerIndexed.String(),
		"false",
	}

	ruleConfigurationModeFields := [][]rulemode.Type{
		ruleConfiguration.DisableModes,
		ruleConfiguration.EnableModes,
		ruleConfiguration.InfoModes,
		ruleConfiguration.WarningModes,
		ruleConfiguration.ErrorModes,
	}

	lmFlagDependentLevel := func() bool {
		// Determine whether the `--library-manager` flag setting affects this rule's level
		for _, ruleConfigurationModeField := range ruleConfigurationModeFields {
			for _, modeConfiguration := range ruleConfigurationModeField {
				for _, libraryManagerMode := range libraryManagerModes {
					if modeConfiguration == libraryManagerMode {
						return true
					}
				}
			}
		}

		return false
	}

	var levelsData [][]string
	if lmFlagDependentLevel() {
		// The `--library-manager` flag setting is used by the rule's configuration, so provide compliance vs. library-manager vs. level data.
		levelsData = append(levelsData, []string{"`compliance`", "`library-manager`", "Level"})
		for _, complianceMode := range complianceModes {
			for _, libraryManagerFlagValue := range libraryManagerFlagValues {
				flags := cli.Root().PersistentFlags()
				if err := flags.Set("compliance", complianceMode.String()); err != nil {
					panic(err)
				}
				if err := flags.Set("library-manager", libraryManagerFlagValue); err != nil {
					panic(err)
				}
				if err := configuration.Initialize(flags, []string{}); err != nil {
					panic(err)
				}
				ruleModes := configuration.RuleModes(ruleConfiguration.ProjectType)
				levelsData = append(levelsData, []string{complianceMode.String(), libraryManagerFlagValue, ruleLevel(ruleConfiguration, ruleModes)})
			}
		}
	} else {
		// The `--library-manager` flag setting is not used by the rule's configuration, so only provide compliance vs. level data.
		levelsData = append(levelsData, []string{"`compliance`", "Level"})
		for _, complianceMode := range complianceModes {
			flags := cli.Root().PersistentFlags()
			if err := flags.Set("compliance", complianceMode.String()); err != nil {
				panic(err)
			}
			if err := configuration.Initialize(flags, []string{}); err != nil {
				panic(err)
			}
			ruleModes := configuration.RuleModes(ruleConfiguration.ProjectType)
			levelsData = append(levelsData, []string{complianceMode.String(), ruleLevel(ruleConfiguration, ruleModes)})
		}
	}

	return levelsData
}

// ruleLevel returns the string representation of the violation level of the given rule in the given mode.
func ruleLevel(ruleConfiguration ruleconfiguration.Type, ruleModes map[rulemode.Type]bool) string {
	enabled, err := rule.IsEnabled(ruleConfiguration, ruleModes)
	if err != nil {
		panic(err)
	}
	if enabled {
		ruleLevel, err := rulelevel.FailRuleLevel(ruleConfiguration, ruleModes)
		if err != nil {
			panic(err)
		}
		return ruleLevel.String()
	}

	return "disabled"
}
