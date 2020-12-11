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

// Package check runs checks on a project.
package check

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result"
	"github.com/arduino/arduino-check/result/feedback"
	"github.com/sirupsen/logrus"
)

// RunChecks runs all checks for the given project and outputs the results.
func RunChecks(project project.Type) {
	feedback.Printf("\nChecking %s in %s\n", project.ProjectType, project.Path)

	checkdata.Initialize(project, configuration.SchemasPath())

	for _, checkConfiguration := range checkconfigurations.Configurations() {
		runCheck, err := shouldRun(checkConfiguration, project)
		if err != nil {
			feedback.Errorf("Error while determining whether to run check: %v", err)
			os.Exit(1)
		}

		if !runCheck {
			logrus.Infof("Skipping check: %s\n", checkConfiguration.ID)
			continue
		}

		// Output will be printed after all checks are finished when configured for "json" output format.
		feedback.VerbosePrintf("Running check %s...\n", checkConfiguration.ID)

		checkResult, checkOutput := checkConfiguration.CheckFunction()
		reportText := result.Results.Record(project, checkConfiguration, checkResult, checkOutput)
		if (checkResult == checkresult.Fail) || configuration.Verbose() {
			feedback.Println(reportText)
		}
	}

	// Checks are finished for this project, so summarize its check results in the report.
	result.Results.AddProjectSummary(project)

	// Print the project check results summary.
	feedback.Printf("\n%s\n", result.Results.ProjectSummaryText(project))
}

// shouldRun returns whether a given check should be run for the given project under the current tool configuration.
func shouldRun(checkConfiguration checkconfigurations.Type, currentProject project.Type) (bool, error) {
	configurationCheckModes := configuration.CheckModes(currentProject.SuperprojectType)

	if checkConfiguration.ProjectType != currentProject.ProjectType {
		return false, nil
	}

	return IsEnabled(checkConfiguration, configurationCheckModes)
}

func IsEnabled(checkConfiguration checkconfigurations.Type, configurationCheckModes map[checkmode.Type]bool) (bool, error) {
	for _, disableMode := range checkConfiguration.DisableModes {
		if configurationCheckModes[disableMode] {
			return false, nil
		}
	}

	for _, enableMode := range checkConfiguration.EnableModes {
		if configurationCheckModes[enableMode] {
			return true, nil
		}
	}

	// Use default
	for _, disableMode := range checkConfiguration.DisableModes {
		if disableMode == checkmode.Default {
			return false, nil
		}
	}

	for _, enableMode := range checkConfiguration.EnableModes {
		if enableMode == checkmode.Default {
			return true, nil
		}
	}

	return false, fmt.Errorf("Check %s is incorrectly configured", checkConfiguration.ID)
}
