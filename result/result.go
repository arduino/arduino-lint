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

// Package result records check results and provides reports and summary text on those results.
package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checklevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/go-paths-helper"
)

// Results is the global instance of the check results result.Type struct
var Results Type

// Type is the type for the check results data
type Type struct {
	Configuration toolConfigurationReportType `json:"configuration"`
	Projects      []projectReportType         `json:"projects"`
	Summary       summaryReportType           `json:"summary"`
}

type toolConfigurationReportType struct {
	Paths       paths.PathList `json:"paths"`
	ProjectType string         `json:"projectType"`
	Recursive   bool           `json:"recursive"`
}

type projectReportType struct {
	Path          *paths.Path                    `json:"path"`
	ProjectType   string                         `json:"projectType"`
	Configuration projectConfigurationReportType `json:"configuration"`
	Checks        []checkReportType              `json:"checks"`
	Summary       summaryReportType              `json:"summary"`
}

type projectConfigurationReportType struct {
	Compliance     string `json:"compliance"`
	LibraryManager string `json:"libraryManager"`
	Official       bool   `json:"official"`
}

type checkReportType struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	ID          string `json:"ID"`
	Brief       string `json:"brief"`
	Description string `json:"description"`
	Result      string `json:"result"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

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

// Record records the result of a check and returns a text summary for it.
func (results *Type) Record(checkedProject project.Type, checkConfiguration checkconfigurations.Type, checkResult checkresult.Type, checkOutput string) string {
	checkLevel, err := checklevel.CheckLevel(checkConfiguration, checkResult)
	if err != nil {
		panic(fmt.Errorf("Error while determining check level: %v", err))
	}

	summaryText := fmt.Sprintf("Check %s result: %s", checkConfiguration.ID, checkResult)

	checkMessage := ""
	if checkResult == checkresult.Fail {
		checkMessage = message(checkConfiguration.MessageTemplate, checkOutput)
	} else {
		// Checks may provide an explanation for their non-fail result.
		// The message template should not be used in this case, since it is written for a failure result.
		checkMessage = checkOutput
	}

	// Add explanation of check result if present.
	if checkMessage != "" {
		summaryText += fmt.Sprintf("\n%s: %s", checkLevel, checkMessage)
	}

	reportExists, projectReportIndex := results.getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		// There is no existing report for this project.
		results.Projects = append(
			results.Projects,
			projectReportType{
				Path:        checkedProject.Path,
				ProjectType: checkedProject.ProjectType.String(),
				Configuration: projectConfigurationReportType{
					Compliance:     checkmode.Compliance(configuration.CheckModes(checkedProject.ProjectType)),
					LibraryManager: checkmode.LibraryManager(configuration.CheckModes(checkedProject.ProjectType)),
					Official:       configuration.CheckModes(checkedProject.ProjectType)[checkmode.Official],
				},
				Checks: []checkReportType{},
			},
		)
	}

	if (checkResult == checkresult.Fail) || configuration.Verbose() {
		checkReport := checkReportType{
			Category:    checkConfiguration.Category,
			Subcategory: checkConfiguration.Subcategory,
			ID:          checkConfiguration.ID,
			Brief:       checkConfiguration.Brief,
			Description: checkConfiguration.Description,
			Result:      checkResult.String(),
			Level:       checkLevel.String(),
			Message:     checkMessage,
		}
		results.Projects[projectReportIndex].Checks = append(results.Projects[projectReportIndex].Checks, checkReport)
	}

	return summaryText
}

// AddProjectSummary summarizes the results of all checks on the given project and adds it to the report.
func (results *Type) AddProjectSummary(checkedProject project.Type) {
	reportExists, projectReportIndex := results.getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary", checkedProject.Path))
	}

	pass := true
	warningCount := 0
	errorCount := 0
	for _, checkReport := range results.Projects[projectReportIndex].Checks {
		if checkReport.Result == checkresult.Fail.String() {
			if checkReport.Level == checklevel.Warning.String() {
				warningCount += 1
			} else if checkReport.Level == checklevel.Error.String() {
				errorCount += 1
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

// ProjectSummaryText returns a text summary of the check results for the given project.
func (results Type) ProjectSummaryText(checkedProject project.Type) string {
	reportExists, projectReportIndex := results.getProjectReportIndex(checkedProject.Path)
	if !reportExists {
		panic(fmt.Sprintf("Unable to find report for %v when generating report summary text", checkedProject.Path))
	}

	projectSummaryReport := results.Projects[projectReportIndex].Summary
	return fmt.Sprintf("Finished checking project. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v", projectSummaryReport.WarningCount, projectSummaryReport.ErrorCount, projectSummaryReport.Pass)
}

// AddSummary summarizes the check results for all projects and adds it to the report.
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

// SummaryText returns a text summary of the cumulative check results.
func (results Type) SummaryText() string {
	return fmt.Sprintf("Finished checking projects. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v", results.Summary.WarningCount, results.Summary.ErrorCount, results.Summary.Pass)
}

// JSONReport returns a JSON formatted report of checks on all projects.
func (results Type) JSONReport() string {
	return string(results.jsonReportRaw())
}

func (results Type) jsonReportRaw() []byte {
	reportJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Error while formatting checks report: %v", err))
	}

	return reportJSON
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

// Passed returns whether the checks passed cumulatively.
func (results Type) Passed() bool {
	return results.Summary.Pass
}

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

// message fills the message template provided by the check configuration with the check output.
// TODO: make checkOutput a struct to allow for more advanced message templating
func message(templateText string, checkOutput string) string {
	messageTemplate := template.Must(template.New("messageTemplate").Parse(templateText))

	messageBuffer := new(bytes.Buffer)
	messageTemplate.Execute(messageBuffer, checkOutput)

	return messageBuffer.String()
}
