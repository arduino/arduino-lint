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

// Package command implements the arduino-check commands.
package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result"
	"github.com/arduino/arduino-check/result/feedback"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/spf13/cobra"
)

// ArduinoCheck is the root command function.
func ArduinoCheck(rootCommand *cobra.Command, cliArguments []string) {
	if err := configuration.Initialize(rootCommand.Flags(), cliArguments); err != nil {
		feedback.Errorf("Configuration error: %v", err)
		os.Exit(1)
	}

	if configuration.VersionMode() {
		if configuration.OutputFormat() == outputformat.Text {
			fmt.Println(configuration.Version() + " " + configuration.BuildTimestamp())
		} else {
			versionObject := struct {
				Version        string `json:"version"`
				BuildTimestamp string `json:"buildTimestamp"`
			}{
				Version:        configuration.Version(),
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
		check.RunChecks(project)
	}

	// All projects have been checked, so summarize their check results in the report.
	result.Results.AddSummary()

	if configuration.OutputFormat() == outputformat.Text {
		if len(projects) > 1 {
			// There are multiple projects, print the summary of check results for all projects.
			fmt.Print(result.Results.SummaryText())
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
