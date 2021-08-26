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

package result

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/rulelevel"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/arduino-lint/internal/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var projectPaths []string

func init() {
	projectPath, err := os.Getwd() // Path to an arbitrary folder that is guaranteed to exist.
	if err != nil {
		panic(err)
	}
	projectPaths = []string{projectPath}
}

func TestInitialize(t *testing.T) {
	flags := test.ConfigurationFlags()
	flags.Set("project-type", "sketch")
	flags.Set("recursive", "false")
	workingDirectoryPath, err := os.Getwd() // A convenient path that is guaranteed to exist.
	require.Nil(t, err)

	err = configuration.Initialize(flags, []string{workingDirectoryPath})
	require.Nil(t, err)
	var results Type
	results.Initialize()
	assert.Equal(t, paths.NewPathList(workingDirectoryPath), results.Configuration.Paths)
	assert.Equal(t, projecttype.Sketch.String(), results.Configuration.ProjectType)
	assert.False(t, results.Configuration.Recursive)
}

func TestRecord(t *testing.T) {
	flags := test.ConfigurationFlags()
	require.Nil(t, configuration.Initialize(flags, projectPaths))

	lintedProject := project.Type{
		Path:             paths.New("/foo/bar"),
		ProjectType:      projecttype.Sketch,
		SuperprojectType: projecttype.Library,
	}

	var results Type
	results.Initialize()
	ruleConfiguration := ruleconfiguration.Configurations()[0]
	ruleOutput := "foo"
	flags.Set("verbose", "true")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	ruleConfiguration.Reference = ""
	summaryText := results.Record(lintedProject, ruleConfiguration, ruleresult.Fail, ruleOutput)
	outputAssertion := "Rule LS001 result: fail\nERROR: Path does not contain a valid Arduino library.\n"
	assert.Equal(t, outputAssertion, summaryText, "No reference URL")
	ruleConfiguration.Reference = "https://arduino.github.io/arduino-cli/latest/library-specification"
	summaryText = results.Record(lintedProject, ruleConfiguration, ruleresult.Fail, ruleOutput)
	outputAssertion = "Rule LS001 result: fail\nERROR: Path does not contain a valid Arduino library.                         \n       See: https://arduino.github.io/arduino-cli/latest/library-specification\n"
	assert.Equal(t, outputAssertion, summaryText, "Reference URL is appended if one is defined")
	summaryText = results.Record(lintedProject, ruleConfiguration, ruleresult.NotRun, ruleOutput)
	assert.Equal(t, fmt.Sprintf("Rule %s result: %s\n%s: %s\n", ruleConfiguration.ID, ruleresult.NotRun, rulelevel.Notice, ruleOutput), summaryText, "Non-fail result should not use message")
	summaryText = results.Record(lintedProject, ruleConfiguration, ruleresult.Pass, "")
	assert.Equal(t, fmt.Sprintf("Rule %s result: %s\n", ruleConfiguration.ID, ruleresult.Pass), summaryText, "Non-failure result with no rule function output should only use preface")
	flags.Set("verbose", "false")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	ruleConfigurationCopy := ruleConfiguration
	ruleConfigurationCopy.MessageTemplate = "bar"
	ruleConfigurationCopy.Reference = ""
	summaryText = results.Record(lintedProject, ruleConfigurationCopy, ruleresult.Fail, ruleOutput)
	outputAssertion = "ERROR: bar (Rule LS001)\n"
	assert.Equal(t, outputAssertion, summaryText, "Rule ID is appended to non-verbose fail message on same line when rule message is single line")
	ruleConfigurationCopy.MessageTemplate = "bar\nbaz"
	summaryText = results.Record(lintedProject, ruleConfigurationCopy, ruleresult.Fail, ruleOutput)
	outputAssertion = "ERROR: bar         \n       baz         \n       (Rule LS001)\n"
	assert.Equal(t, outputAssertion, summaryText, "Rule ID is appended to non-verbose fail message on same line when rule message is multiple lines")
	summaryText = results.Record(lintedProject, ruleConfiguration, ruleresult.NotRun, ruleOutput)
	assert.Equal(t, "", summaryText, "Non-fail result should not result in output in non-verbose mode")
	summaryText = results.Record(lintedProject, ruleConfiguration, ruleresult.Pass, "")
	assert.Equal(t, "", summaryText, "Non-fail result should not result in output in non-verbose mode")

	flags.Set("verbose", "true")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	ruleResult := ruleresult.Pass
	results.Initialize()
	results.Record(lintedProject, ruleConfiguration, ruleResult, ruleOutput)
	projectReport := results.Projects[0]
	assert.Equal(t, lintedProject.Path, projectReport.Path)
	assert.Equal(t, lintedProject.ProjectType.String(), projectReport.ProjectType)
	projectConfigurationReport := projectReport.Configuration
	assert.Equal(t, rulemode.Compliance(configuration.RuleModes(lintedProject.ProjectType)), projectConfigurationReport.Compliance)
	assert.Equal(t, rulemode.LibraryManager(configuration.RuleModes(lintedProject.ProjectType)), projectConfigurationReport.LibraryManager)
	assert.Equal(t, configuration.RuleModes(lintedProject.ProjectType)[rulemode.Official], projectConfigurationReport.Official)
	assert.Equal(t, 1, len(results.Projects[0].Rules), "Passing rule reports should be written to report in verbose mode")
	ruleReport := projectReport.Rules[0]
	assert.Equal(t, ruleConfiguration.Category, ruleReport.Category)
	assert.Equal(t, ruleConfiguration.Subcategory, ruleReport.Subcategory)
	assert.Equal(t, ruleConfiguration.ID, ruleReport.ID)
	assert.Equal(t, ruleConfiguration.Brief, ruleReport.Brief)
	assert.Equal(t, ruleConfiguration.Description, ruleReport.Description)
	assert.Equal(t, ruleResult.String(), ruleReport.Result)
	ruleLevel, _ := rulelevel.RuleLevel(ruleConfiguration, ruleResult, lintedProject)
	assert.Equal(t, ruleLevel.String(), ruleReport.Level)
	assert.Equal(t, ruleOutput, ruleReport.Message)

	flags.Set("verbose", "false")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	results.Initialize()
	results.Record(lintedProject, ruleConfiguration, ruleresult.Pass, ruleOutput)
	assert.Equal(t, 0, len(results.Projects[0].Rules), "Passing rule reports should not be written to report in non-verbose mode")

	results.Initialize()
	results.Record(lintedProject, ruleConfiguration, ruleresult.Fail, ruleOutput)
	require.Equal(t, 1, len(projectReport.Rules), "Failing rule reports should be written to report in non-verbose mode")

	assert.Len(t, results.Projects, 1)
	previousProjectPath := lintedProject.Path
	lintedProject.Path = paths.New("/foo/baz")
	results.Record(lintedProject, ruleConfiguration, ruleresult.Fail, ruleOutput)
	assert.Len(t, results.Projects, 2)

	assert.Len(t, results.Projects[0].Rules, 1)
	lintedProject.Path = previousProjectPath
	results.Record(lintedProject, ruleconfiguration.Configurations()[1], ruleresult.Fail, ruleOutput)
	assert.Len(t, results.Projects[0].Rules, 2)
}

func TestAddProjectSummary(t *testing.T) {
	lintedProject := project.Type{
		Path:             paths.New("/foo/bar"),
		ProjectType:      projecttype.Sketch,
		SuperprojectType: projecttype.Library,
	}

	testTables := []struct {
		results              []ruleresult.Type
		levels               []rulelevel.Type
		verbose              string
		expectedPass         bool
		expectedWarningCount int
		expectedErrorCount   int
	}{
		{
			[]ruleresult.Type{ruleresult.Pass, ruleresult.Pass},
			[]rulelevel.Type{rulelevel.Info, rulelevel.Info},
			"true",
			true,
			0,
			0,
		},
		{
			[]ruleresult.Type{ruleresult.Pass, ruleresult.Pass},
			[]rulelevel.Type{rulelevel.Info, rulelevel.Info},
			"false",
			true,
			0,
			0,
		},
		{
			[]ruleresult.Type{ruleresult.Pass, ruleresult.Fail},
			[]rulelevel.Type{rulelevel.Info, rulelevel.Warning},
			"false",
			true,
			1,
			0,
		},
		{
			[]ruleresult.Type{ruleresult.Fail, ruleresult.Fail},
			[]rulelevel.Type{rulelevel.Error, rulelevel.Warning},
			"false",
			false,
			1,
			1,
		},
	}

	for _, testTable := range testTables {
		flags := test.ConfigurationFlags()
		flags.Set("verbose", testTable.verbose)
		require.Nil(t, configuration.Initialize(flags, projectPaths))

		var results Type
		results.Initialize()

		ruleIndex := 0
		for testDataIndex, result := range testTable.results {
			results.Record(lintedProject, ruleconfiguration.Configurations()[0], result, "")
			if (result == ruleresult.Fail) || configuration.Verbose() {
				level := testTable.levels[testDataIndex].String()
				results.Projects[0].Rules[ruleIndex].Level = level
				ruleIndex++
			}
		}
		results.AddProjectSummary(lintedProject)
		assert.Equal(t, testTable.expectedPass, results.Projects[0].Summary.Pass)
		assert.Equal(t, testTable.expectedWarningCount, results.Projects[0].Summary.WarningCount)
		assert.Equal(t, testTable.expectedErrorCount, results.Projects[0].Summary.ErrorCount)
		if testTable.expectedErrorCount == 0 && testTable.expectedWarningCount == 0 {
			assert.Equal(t, "Linter results for project: no errors or warnings", results.ProjectSummaryText(lintedProject))
		} else {
			assert.Equal(t, fmt.Sprintf("Linter results for project: %v ERRORS, %v WARNINGS", testTable.expectedErrorCount, testTable.expectedWarningCount), results.ProjectSummaryText(lintedProject))
		}
	}
}

func TestAddSummary(t *testing.T) {
	lintedProject := project.Type{
		Path:             paths.New("/foo/bar"),
		ProjectType:      projecttype.Sketch,
		SuperprojectType: projecttype.Library,
	}

	testTables := []struct {
		projectSummaries     []summaryReportType
		expectedPass         bool
		expectedWarningCount int
		expectedErrorCount   int
	}{
		{
			[]summaryReportType{
				{
					Pass:         true,
					WarningCount: 0,
					ErrorCount:   0,
				},
				{
					Pass:         true,
					WarningCount: 0,
					ErrorCount:   0,
				},
			},
			true,
			0,
			0,
		},
		{
			[]summaryReportType{
				{
					Pass:         true,
					WarningCount: 1,
					ErrorCount:   0,
				},
				{
					Pass:         true,
					WarningCount: 2,
					ErrorCount:   0,
				},
			},
			true,
			3,
			0,
		},
		{
			[]summaryReportType{
				{
					Pass:         false,
					WarningCount: 1,
					ErrorCount:   0,
				},
				{
					Pass:         true,
					WarningCount: 2,
					ErrorCount:   2,
				},
			},
			false,
			3,
			2,
		},
	}

	for _, testTable := range testTables {
		var results Type
		for projectIndex, projectSummary := range testTable.projectSummaries {
			lintedProject.Path = paths.New(fmt.Sprintf("/foo/bar%v", projectIndex)) // Use a unique path to generate a new project report.
			results.Record(lintedProject, ruleconfiguration.Configurations()[0], ruleresult.Pass, "")
			results.AddProjectSummary(lintedProject)
			results.Projects[projectIndex].Summary = projectSummary
		}
		results.AddSummary()
		assert.Equal(t, testTable.expectedPass, results.Summary.Pass)
		assert.Equal(t, testTable.expectedPass, results.Passed())
		assert.Equal(t, testTable.expectedWarningCount, results.Summary.WarningCount)
		assert.Equal(t, testTable.expectedErrorCount, results.Summary.ErrorCount)
		if testTable.expectedErrorCount == 0 && testTable.expectedWarningCount == 0 {
			assert.Equal(t, "Linter results for projects: no errors or warnings", results.SummaryText())
		} else {
			assert.Equal(t, fmt.Sprintf("Linter results for projects: %v ERRORS, %v WARNINGS", testTable.expectedErrorCount, testTable.expectedWarningCount), results.SummaryText())
		}
	}
}

func TestWriteReport(t *testing.T) {
	flags := test.ConfigurationFlags()

	reportFolderPathString, err := ioutil.TempDir("", "arduino-lint-result-TestWriteReport")
	require.Nil(t, err)
	defer os.RemoveAll(reportFolderPathString) // clean up
	reportFolderPath := paths.New(reportFolderPathString)

	reportFilePath := reportFolderPath.Join("report-file.json")
	_, err = reportFilePath.Create() // Create file using the report folder path.
	require.Nil(t, err)

	flags.Set("report-file", reportFilePath.Join("report-file.json").String())
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	assert.Error(t, Results.WriteReport(), "Parent folder creation should fail due to a collision with an existing file at that path")

	reportFilePath = reportFolderPath.Join("report-file-subfolder", "report-file-subsubfolder", "report-file.json")
	flags.Set("report-file", reportFilePath.String())
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	assert.NoError(t, Results.WriteReport(), "Creation of multiple levels of parent folders")

	reportFile, err := reportFilePath.Open()
	require.Nil(t, err)
	reportFileInfo, err := reportFile.Stat()
	require.Nil(t, err)
	reportFileBytes := make([]byte, reportFileInfo.Size())
	_, err = reportFile.Read(reportFileBytes)
	require.Nil(t, err)
	assert.True(t, assert.ObjectsAreEqualValues(reportFileBytes, Results.jsonReportRaw()), "Report file contents are correct")
}
