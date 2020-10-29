// Package check runs checks on a project.
package check

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result"
	"github.com/arduino/arduino-check/result/feedback"
	"github.com/sirupsen/logrus"
)

// RunChecks runs all checks for the given project and outputs the results.
func RunChecks(project project.Type) {
	fmt.Printf("Checking %s in %s\n", project.ProjectType, project.Path)

	checkdata.Initialize(project)

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

		// Output will be printed after all checks are finished when configured for "json" output format
		if configuration.OutputFormat() == "text" {
			fmt.Printf("Running check %s: ", checkConfiguration.ID)
		}
		checkResult, checkOutput := checkConfiguration.CheckFunction()
		reportText := result.Results.Record(project, checkConfiguration, checkResult, checkOutput)
		if configuration.OutputFormat() == "text" {
			fmt.Print(reportText)
		}
	}

	// Checks are finished for this project, so summarize its check results in the report.
	result.Results.AddProjectSummary(project)

	if configuration.OutputFormat() == "text" {
		// Print the project check results summary.
		fmt.Print(result.Results.ProjectSummaryText(project))
	}
}

// shouldRun returns whether a given check should be run for the given project under the current tool configuration.
func shouldRun(checkConfiguration checkconfigurations.Type, currentProject project.Type) (bool, error) {
	configurationCheckModes := configuration.CheckModes(currentProject.SuperprojectType)

	if checkConfiguration.ProjectType != currentProject.ProjectType {
		return false, nil
	}

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
