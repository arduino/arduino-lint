// This file is part of arduino-lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package command implements the arduino-lint commands.
package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arduino/arduino-lint/check"
	"github.com/arduino/arduino-lint/configuration"
	"github.com/arduino/arduino-lint/project"
	"github.com/arduino/arduino-lint/result"
	"github.com/arduino/arduino-lint/result/feedback"
	"github.com/arduino/arduino-lint/result/outputformat"
	"github.com/spf13/cobra"
)

// ArduinoLint is the root command function.
func ArduinoLint(rootCommand *cobra.Command, cliArguments []string) {
	if err := configuration.Initialize(rootCommand.Flags(), cliArguments); err != nil {
		feedback.Errorf("Invalid configuration: %v", err)
		os.Exit(1)
	}

	if configuration.VersionMode() {
		if configuration.OutputFormat() == outputformat.Text {
			fmt.Println(configuration.VersionInfo())
		} else {
			versionJSON, err := json.MarshalIndent(configuration.VersionInfo(), "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(versionJSON))
		}
		return
	}

	result.Results.Initialize()

	projects, err := project.FindProjects()
	if err != nil {
		feedback.Errorf("Error while finding projects: %v", err)
		os.Exit(1)
	}

	for _, project := range projects {
		check.RunChecks(project)

		// Checks are finished for this project, so summarize its check results in the report.
		result.Results.AddProjectSummary(project)

		// Print the project check results summary.
		feedback.Printf("\n%s\n", result.Results.ProjectSummaryText(project))
	}

	// All projects have been checked, so summarize their check results in the report.
	result.Results.AddSummary()

	if configuration.OutputFormat() == outputformat.Text {
		if len(projects) > 1 {
			// There are multiple projects, print the summary of check results for all projects.
			fmt.Printf("\n%s\n", result.Results.SummaryText())
		}
	} else {
		// Print the complete JSON formatted report.
		fmt.Println(result.Results.JSONReport())
	}

	if configuration.ReportFilePath() != nil {
		// Write report file.
		if err := result.Results.WriteReport(); err != nil {
			feedback.Error(err.Error())
			os.Exit(1)
		}
	}

	if !result.Results.Passed() {
		os.Exit(1)
	}
}
