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

// Package result records rule results and provides reports and summary text on those results.
package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/rulelevel"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/go-paths-helper"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// Results is the global instance of the rule results result.Type struct
var Results Type

// Type is the type for the rule results data
type Type struct {
	Configuration toolConfigurationReportType `json:"configuration"`
	Projects      []projectReportType         `json:"projects"`
	Summary       summaryReportType           `json:"summary"`
}

// toolConfigurationReportType is the type for the Arduino Lint tool configuration.
type toolConfigurationReportType struct {
	Paths       paths.PathList `json:"paths"`
	ProjectType string         `json:"projectType"`
	Recursive   bool           `json:"recursive"`
}

// projectReportType is the type for the individual project reports.
type projectReportType struct {
	Path          *paths.Path                    `json:"path"`
	ProjectType   string                         `json:"projectType"`
	Configuration projectConfigurationReportType `json:"configuration"`
	Rules         []ruleReportType               `json:"rules"`
	Summary       summaryReportType              `json:"summary"`
}

// projectConfigurationReportType is the type for the individual project tool configurations.
type projectConfigurationReportType struct {
	Compliance     string `json:"compliance"`
	LibraryManager string `json:"libraryManager"`
	Official       bool   `json:"official"`
}

// ruleReportType is the type of the rule reports.
type ruleReportType struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	ID          string `json:"ID"`
	Brief       string `json:"brief"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

// summaryReportType is the type of the rule result summary reports.
type summaryReportType struct {
	Pass         bool `json:"pass"`
	WarningCount int  `json:"warningCount"`
	ErrorCount   int  `json:"errorCount"`
}

// Initialize adds the tool configuration data to the results data.
func (results *Type) Initialize() {
	*results = *new(Type)
	results.Configuration = toolConfigurationReportType{
		Paths:       configuration.TargetPaths(),
		ProjectType: configuration.SuperprojectTypeFilter().String(),
		Recursive:   configuration.Recursive(),
	}
}

var blankLineRegexp = regexp.MustCompile("\n[[:space:]]*\n")

// Record records the result of a rule and returns a text summary for it.
func (results *Type) Record(lintedProject project.Type, ruleConfiguration ruleconfiguration.Type, ruleResult ruleresult.Type, ruleOutput string) string {
	ruleLevel, err := rulelevel.RuleLevel(ruleConfiguration, ruleResult, lintedProject)
	if err != nil {
		panic(fmt.Errorf("Error while determining rule level: %v", err))
	}

	ruleMessage := ""
	if ruleResult == ruleresult.Fail {
		ruleMessage = message(ruleConfiguration.MessageTemplate, ruleOutput)
		if ruleConfiguration.Reference != "" {
			ruleMessage = fmt.Sprintf("%s\nSee: %s", ruleMessage, ruleConfiguration.Reference)
		}
	} else {
		// Rules may provide an explanation for their non-fail result.
		// The message template should not be used in this case, since it is written for a failure result.
		ruleMessage = ruleOutput
	}

	summaryText := ""

	formatRuleText := func(level rulelevel.Type, message string) string {
		width := 120 // Wrap text to this width
		prefix := fmt.Sprintf("%s: ", level)

		formattedOutput := &strings.Builder{}

		tableConfigBuilder := tablewriter.NewConfigBuilder()
		// Configure column widths so that the text will be wrapped to the appropriate width.
		tableConfigBuilder.ForColumn(0).WithMaxWidth(len(prefix))
		tableConfigBuilder.ForColumn(1).WithMaxWidth(width - len(prefix))
		// A trailing space is intentionally added to the "prefix" string. Trimming must be disabled to preserve that space.
		tableConfigBuilder.WithTrimSpace(tw.Off)
		tableConfig := tableConfigBuilder.Build()

		tableRendition := tw.Rendition{
			// Do not add border characters to the table.
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				/*
					Do not add a separator character between columns. The trailing space on the "prefix" string serves as the
					separator.
				*/
				Separators: tw.SeparatorsNone,
			},
		}

		tableRenderer := renderer.NewBlueprint(tableRendition)

		table := tablewriter.NewTable(
			formattedOutput,
			tablewriter.WithConfig(tableConfig),
			tablewriter.WithRenderer(tableRenderer),
			// Do not add margins to the cell content.
			tablewriter.WithPadding(tw.PaddingNone),
		)

		if err := table.Append([]string{prefix, message}); err != nil {
			panic(err)
		}

		if err := table.Render(); err != nil {
			panic(err)
		}

		return formattedOutput.String()
	}

	if configuration.Verbose() {
		summaryText = fmt.Sprintf("Rule %s result: %s\n", ruleConfiguration.ID, ruleResult)
		// Add explanation of rule result if present.
		if ruleMessage != "" {
			summaryText += formatRuleText(ruleLevel, ruleMessage)
		}
	} else {
		if ruleResult == ruleresult.Fail {
			if strings.Contains(ruleMessage, "\n") {
				summaryText = formatRuleText(ruleLevel, fmt.Sprintf("%s\n(Rule %s)", ruleMessage, ruleConfiguration.ID))
			} else {
				summaryText = formatRuleText(ruleLevel, fmt.Sprintf("%s (Rule %s)", ruleMessage, ruleConfiguration.ID))
			}
		}
	}

	reportExists, projectReportIndex := results.getProjectReportIndex(lintedProject.Path)
	if !reportExists {
		// There is no existing report for this project.
		results.Projects = append(
			results.Projects,
			projectReportType{
				Path:        lintedProject.Path,
				ProjectType: lintedProject.ProjectType.String(),
				Configuration: projectConfigurationReportType{
					Compliance:     rulemode.Compliance(configuration.RuleModes(lintedProject.ProjectType)),
					LibraryManager: rulemode.LibraryManager(configuration.RuleModes(lintedProject.ProjectType)),
					Official:       configuration.RuleModes(lintedProject.ProjectType)[rulemode.Official],
				},
				Rules: []ruleReportType{},
			},
		)
	}

	if (ruleResult == ruleresult.Fail) || configuration.Verbose() {
		ruleReport := ruleReportType{
			Category:    ruleConfiguration.Category,
			Subcategory: ruleConfiguration.Subcategory,
			ID:          ruleConfiguration.ID,
			Brief:       ruleConfiguration.Brief,
			Description: ruleConfiguration.Description,
			Result:      ruleResult.String(),
			Level:       ruleLevel.String(),
			Message:     ruleMessage,
		}
		results.Projects[projectReportIndex].Rules = append(results.Projects[projectReportIndex].Rules, ruleReport)
	}

	return summaryText
}

// AddProjectSummary summarizes the results of all rules on the given project and adds it to the report.
func (results *Type) AddProjectSummary(lintedProject project.Type) {
	reportExists, projectReportIndex := results.getProjectReportIndex(lintedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary", lintedProject.Path))
	}

	pass := true
	warningCount := 0
	errorCount := 0
	for _, ruleReport := range results.Projects[projectReportIndex].Rules {
		if ruleReport.Result == ruleresult.Fail.String() {
			if ruleReport.Level == rulelevel.Warning.String() {
				warningCount++
			} else if ruleReport.Level == rulelevel.Error.String() {
				errorCount++
				pass = false
			}
		}
	}

	results.Projects[projectReportIndex].Summary = summaryReportType{
		Pass:         pass,
		WarningCount: warningCount,
		ErrorCount:   errorCount,
	}
}

// ProjectSummaryText returns a text summary of the rule results for the given project.
func (results Type) ProjectSummaryText(lintedProject project.Type) string {
	reportExists, projectReportIndex := results.getProjectReportIndex(lintedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary text", lintedProject.Path))
	}

	projectSummaryReport := "Linter results for project: "
	projectSummaryData := results.Projects[projectReportIndex].Summary
	if projectSummaryData.ErrorCount == 0 && projectSummaryData.WarningCount == 0 {
		projectSummaryReport += "no errors or warnings"
	} else {
		projectSummaryReport += fmt.Sprintf("%v ERRORS, %v WARNINGS", projectSummaryData.ErrorCount, projectSummaryData.WarningCount)
	}

	return projectSummaryReport
}

// AddSummary summarizes the rule results for all projects and adds it to the report.
func (results *Type) AddSummary() {
	pass := true
	warningCount := 0
	errorCount := 0
	for _, projectReport := range results.Projects {
		if !projectReport.Summary.Pass {
			pass = false
		}
		warningCount += projectReport.Summary.WarningCount
		errorCount += projectReport.Summary.ErrorCount
	}

	results.Summary = summaryReportType{
		Pass:         pass,
		WarningCount: warningCount,
		ErrorCount:   errorCount,
	}
}

// SummaryText returns a text summary of the cumulative rule results.
func (results Type) SummaryText() string {
	summaryReport := "Linter results for projects: "
	if results.Summary.ErrorCount == 0 && results.Summary.WarningCount == 0 {
		summaryReport += "no errors or warnings"
	} else {
		summaryReport += fmt.Sprintf("%v ERRORS, %v WARNINGS", results.Summary.ErrorCount, results.Summary.WarningCount)
	}

	return summaryReport
}

// JSONReport returns a JSON formatted report of rules on all projects in string encoding.
func (results Type) JSONReport() string {
	return string(results.jsonReportRaw())
}

// jsonReportRaw returns the report marshaled into JSON format in byte encoding.
func (results Type) jsonReportRaw() []byte {
	var marshaledReportBuffer bytes.Buffer
	jsonEncoder := json.NewEncoder(io.Writer(&marshaledReportBuffer))
	// By default, the json package HTML-sanitizes strings during marshaling (https://golang.org/pkg/encoding/json/#Marshal)
	// This means that the simple json.MarshalIndent() approach would result in the report containing gibberish.
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	err := jsonEncoder.Encode(results)
	if err != nil {
		panic(fmt.Sprintf("Error while formatting rules report: %v", err))
	}

	return marshaledReportBuffer.Bytes()
}

// WriteReport writes a report for all projects to the specified file.
func (results Type) WriteReport() error {
	reportFilePath := configuration.ReportFilePath()
	reportFilePathParentExists, err := reportFilePath.Parent().ExistCheck()
	if err != nil {
		return fmt.Errorf("Problem processing --report-file flag value %v: %v", reportFilePath, err)
	}
	if !reportFilePathParentExists {
		err = reportFilePath.Parent().MkdirAll()
		if err != nil {
			return fmt.Errorf("Unable to create report file path (%v): %v", reportFilePath.Parent(), err)
		}
	}

	err = reportFilePath.WriteFile(results.jsonReportRaw())
	if err != nil {
		return fmt.Errorf("While writing report: %v", err)
	}

	return nil
}

// Passed returns whether the rules passed cumulatively.
func (results Type) Passed() bool {
	return results.Summary.Pass
}

// getProjectReportIndex returns the index of the existing entry in the results.Projects array for the given project, or the next available index if there is no existing entry.
func (results Type) getProjectReportIndex(projectPath *paths.Path) (bool, int) {
	var index int
	var projectReport projectReportType
	for index, projectReport = range results.Projects {
		if projectReport.Path == projectPath {
			return true, index
		}
	}

	// There is no element in the report for this project.
	return false, len(results.Projects)
}

// message fills the message template provided by the rule configuration with the rule output.
// TODO: make ruleOutput a struct to allow for more advanced message templating
func message(templateText string, ruleOutput string) string {
	messageTemplate := template.Must(template.New("messageTemplate").Parse(templateText))

	messageBuffer := new(bytes.Buffer)
	messageTemplate.Execute(messageBuffer, ruleOutput)

	return messageBuffer.String()
}
