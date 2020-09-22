package main

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result"
	"github.com/arduino/arduino-check/result/feedback"
)

func main() {
	configuration.Initialize()
	// Must be called after configuration.Initialize()
	result.Initialize()

	projects, err := project.FindProjects()
	if err != nil {
		feedback.Errorf("Error while finding projects: %v", err)
		os.Exit(1)
	}

	for _, project := range projects {
		check.RunChecks(project)
	}

	// All projects have been checked, so summarize their check results in the report.
	result.AddSummaryReport()

	if configuration.OutputFormat() == "text" {
		if len(projects) > 1 {
			// There are multiple projects, print the summary of check results for all projects.
			fmt.Print(result.SummaryText())
		}
	} else {
		// Print the complete JSON formatted report.
		fmt.Println(result.JSONReport())
	}

	if !result.Passed() {
		os.Exit(1)
	}
}
