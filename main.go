package main

import (
	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project"
)

func main() {
	configuration.Initialize()
	projects := project.FindProjects()
	for _, project := range projects {
		check.RunChecks(project)
	}
	// TODO: set exit status according to check results
}
