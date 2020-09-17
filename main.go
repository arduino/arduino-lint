package main

import (
	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/projects"
)

func main() {
	configuration.Initialize()
	projects := projects.FindProjects()
	for _, project := range projects {
		check.RunChecks(project)
	}
}
