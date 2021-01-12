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

// Package command implements the arduino-lint commands.
package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/result"
	"github.com/arduino/arduino-lint/internal/result/feedback"
	"github.com/arduino/arduino-lint/internal/result/outputformat"
	"github.com/arduino/arduino-lint/internal/rule"
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
			if configuration.Version() == "" {
				fmt.Print("0.0.0+" + configuration.Commit())
			} else {
				fmt.Print(configuration.Version())
			}
			fmt.Println(" " + configuration.BuildTimestamp())
		} else {
			versionObject := struct {
				Version        string `json:"version"`
				Commit         string `json:"commit"`
				BuildTimestamp string `json:"buildTimestamp"`
			}{
				Version:        configuration.Version(),
				Commit:         configuration.Commit(),
				BuildTimestamp: configuration.BuildTimestamp(),
			}
			versionJSON, err := json.MarshalIndent(versionObject, "", "  ")
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
		rule.Runner(project)

		// Rules are finished for this project, so summarize its rule results in the report.
		result.Results.AddProjectSummary(project)

		// Print the project rule results summary.
		feedback.Printf("\n%s\n", result.Results.ProjectSummaryText(project))
		feedback.Print("\n-------------------\n\n")
	}

	// All projects have been linted, so summarize their rule results in the report.
	result.Results.AddSummary()

	if configuration.OutputFormat() == outputformat.Text {
		if len(projects) > 1 {
			// There are multiple projects, print the summary of rule results for all projects.
			fmt.Println(result.Results.SummaryText())
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
