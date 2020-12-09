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

package result

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checklevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/util/test"
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
	fmt.Printf("paths: %s", configuration.TargetPaths())
	assert.Equal(t, paths.NewPathList(workingDirectoryPath), results.Configuration.Paths)
	assert.Equal(t, projecttype.Sketch.String(), results.Configuration.ProjectType)
	assert.False(t, results.Configuration.Recursive)
}

func TestRecord(t *testing.T) {
	flags := test.ConfigurationFlags()
	require.Nil(t, configuration.Initialize(flags, projectPaths))

	checkedProject := project.Type{
		Path:             paths.New("/foo/bar"),
		ProjectType:      projecttype.Sketch,
		SuperprojectType: projecttype.Library,
	}

	var results Type
	results.Initialize()
	checkConfiguration := checkconfigurations.Configurations()[0]
	checkOutput := "foo"
	summaryText := results.Record(checkedProject, checkConfiguration, checkresult.Fail, checkOutput)
	assert.Equal(t, fmt.Sprintf("Check %s result: %s\n%s: %s\n", checkConfiguration.ID, checkresult.Fail, checklevel.Error, message(checkConfiguration.MessageTemplate, checkOutput)), summaryText)
	summaryText = results.Record(checkedProject, checkConfiguration, checkresult.NotRun, checkOutput)
	assert.Equal(t, fmt.Sprintf("Check %s result: %s\n%s: %s\n", checkConfiguration.ID, checkresult.NotRun, checklevel.Notice, checkOutput), summaryText, "Non-fail result should not use message")
	summaryText = results.Record(checkedProject, checkConfiguration, checkresult.Pass, "")
	assert.Equal(t, "", "", summaryText, "Non-failure result with no check function output should result in an empty summary")

	flags.Set("verbose", "true")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	checkResult := checkresult.Pass
	results.Initialize()
	results.Record(checkedProject, checkConfiguration, checkResult, checkOutput)
	projectReport := results.Projects[0]
	assert.Equal(t, checkedProject.Path, projectReport.Path)
	assert.Equal(t, checkedProject.ProjectType.String(), projectReport.ProjectType)
	projectConfigurationReport := projectReport.Configuration
	assert.Equal(t, checkmode.Compliance(configuration.CheckModes(checkedProject.ProjectType)), projectConfigurationReport.Compliance)
	assert.Equal(t, configuration.CheckModes(checkedProject.ProjectType)[checkmode.LibraryManagerSubmission], projectConfigurationReport.LibraryManagerSubmit)
	assert.Equal(t, configuration.CheckModes(checkedProject.ProjectType)[checkmode.LibraryManagerIndexed], projectConfigurationReport.LibraryManagerUpdate)
	assert.Equal(t, configuration.CheckModes(checkedProject.ProjectType)[checkmode.Official], projectConfigurationReport.Official)
	assert.Equal(t, 1, len(results.Projects[0].Checks), "Passing check reports should be written to report in verbose mode")
	checkReport := projectReport.Checks[0]
	assert.Equal(t, checkConfiguration.Category, checkReport.Category)
	assert.Equal(t, checkConfiguration.Subcategory, checkReport.Subcategory)
	assert.Equal(t, checkConfiguration.ID, checkReport.ID)
	assert.Equal(t, checkConfiguration.Brief, checkReport.Brief)
	assert.Equal(t, checkConfiguration.Description, checkReport.Description)
	assert.Equal(t, checkResult.String(), checkReport.Result)
	checkLevel, _ := checklevel.CheckLevel(checkConfiguration, checkResult)
	assert.Equal(t, checkLevel.String(), checkReport.Level)
	assert.Equal(t, checkOutput, checkReport.Message)

	flags.Set("verbose", "false")
	require.Nil(t, configuration.Initialize(flags, projectPaths))
	results.Initialize()
	results.Record(checkedProject, checkConfiguration, checkresult.Pass, checkOutput)
	assert.Equal(t, 0, len(results.Projects[0].Checks), "Passing check reports should not be written to report in non-verbose mode")

	results.Initialize()
	results.Record(checkedProject, checkConfiguration, checkresult.Fail, checkOutput)
	require.Equal(t, 1, len(projectReport.Checks), "Failing check reports should be written to report in non-verbose mode")

	assert.Len(t, results.Projects, 1)
	previousProjectPath := checkedProject.Path
	checkedProject.Path = paths.New("/foo/baz")
	results.Record(checkedProject, checkConfiguration, checkresult.Fail, checkOutput)
	assert.Len(t, results.Projects, 2)

	assert.Len(t, results.Projects[0].Checks, 1)
	checkedProject.Path = previousProjectPath
	results.Record(checkedProject, checkconfigurations.Configurations()[1], checkresult.Fail, checkOutput)
	assert.Len(t, results.Projects[0].Checks, 2)
}

func TestAddProjectSummary(t *testing.T) {
	checkedProject := project.Type{
		Path:             paths.New("/foo/bar"),
		ProjectType:      projecttype.Sketch,
		SuperprojectType: projecttype.Library,
	}

	testTables := []struct {
		results              []checkresult.Type
		levels               []checklevel.Type
		verbose              string
		expectedPass         bool
		expectedWarningCount int
		expectedErrorCount   int
	}{
		{
			[]checkresult.Type{checkresult.Pass, checkresult.Pass},
			[]checklevel.Type{checklevel.Info, checklevel.Info},
			"true",
			true,
			0,
			0,
		},
		{
			[]checkresult.Type{checkresult.Pass, checkresult.Pass},
			[]checklevel.Type{checklevel.Info, checklevel.Info},
			"false",
			true,
			0,
			0,
		},
		{
			[]checkresult.Type{checkresult.Pass, checkresult.Fail},
			[]checklevel.Type{checklevel.Info, checklevel.Warning},
			"false",
			true,
			1,
			0,
		},
		{
			[]checkresult.Type{checkresult.Fail, checkresult.Fail},
			[]checklevel.Type{checklevel.Error, checklevel.Warning},
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

		checkIndex := 0
		for testDataIndex, result := range testTable.results {
			results.Record(checkedProject, checkconfigurations.Configurations()[0], result, "")
			if (result == checkresult.Fail) || configuration.Verbose() {
				level := testTable.levels[testDataIndex].String()
				results.Projects[0].Checks[checkIndex].Level = level
				checkIndex += 1
			}
		}
		results.AddProjectSummary(checkedProject)
		assert.Equal(t, testTable.expectedPass, results.Projects[0].Summary.Pass)
		assert.Equal(t, testTable.expectedWarningCount, results.Projects[0].Summary.WarningCount)
		assert.Equal(t, testTable.expectedErrorCount, results.Projects[0].Summary.ErrorCount)
		assert.Equal(t, fmt.Sprintf("\nFinished checking project. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v\n\n", testTable.expectedWarningCount, testTable.expectedErrorCount, testTable.expectedPass), results.ProjectSummaryText(checkedProject))
	}
}

func TestAddSummary(t *testing.T) {
	checkedProject := project.Type{
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
			checkedProject.Path = paths.New(fmt.Sprintf("/foo/bar%v", projectIndex)) // Use a unique path to generate a new project report.
			results.Record(checkedProject, checkconfigurations.Configurations()[0], checkresult.Pass, "")
			results.AddProjectSummary(checkedProject)
			results.Projects[projectIndex].Summary = projectSummary
		}
		results.AddSummary()
		assert.Equal(t, testTable.expectedPass, results.Summary.Pass)
		assert.Equal(t, testTable.expectedPass, results.Passed())
		assert.Equal(t, testTable.expectedWarningCount, results.Summary.WarningCount)
		assert.Equal(t, testTable.expectedErrorCount, results.Summary.ErrorCount)
		assert.Equal(t, fmt.Sprintf("Finished checking projects. Results:\nWarning count: %v\nError count: %v\nChecks passed: %v\n", testTable.expectedWarningCount, testTable.expectedErrorCount, testTable.expectedPass), results.SummaryText())
	}
}

func TestWriteReport(t *testing.T) {
	flags := test.ConfigurationFlags()

	reportFolderPathString, err := ioutil.TempDir("", "arduino-check-result-TestWriteReport")
	require.Nil(t, err)
	defer os.RemoveAll(reportFolderPathString) // clean up
	reportFolderPath := paths.New(reportFolderPathString)

	reportFilePath := reportFolderPath.Join("report-file.json")
	_, err = reportFilePath.Create() // Create file using the report folder path.
	require.Nil(t, err)

	flags.Set("report-file", reportFilePath.Join("report-file.json").String())
	assert.Nil(t, configuration.Initialize(flags, projectPaths))
	assert.Error(t, Results.WriteReport(), "Parent folder creation should fail due to a collision with an existing file at that path")

	reportFilePath = reportFolderPath.Join("report-file-subfolder", "report-file-subsubfolder", "report-file.json")
	flags.Set("report-file", reportFilePath.String())
	assert.Nil(t, configuration.Initialize(flags, projectPaths))
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
