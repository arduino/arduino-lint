package main

import (
	"os"

	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result/feedback"
)

func main() {
	configuration.Initialize()
	projects, err := project.FindProjects()
	if err != nil {
		feedback.Errorf("Error while finding projects: %v", err)
		os.Exit(1)
	}
	for _, project := range projects {
		check.RunChecks(project)
	}
	// TODO: set exit status according to check results
}
