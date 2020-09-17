package check

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checklevel"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/projects"
)

func shouldRun(checkConfiguration checkconfigurations.Configuration, currentProject projects.Project) bool {
	checkModes := configuration.CheckModes(projects.SuperprojectType(currentProject))

	if checkConfiguration.ProjectType != currentProject.Type {
		return false
	}

	for _, disableMode := range checkConfiguration.DisableModes {
		if checkModes[disableMode] == true {
			return false
		}
	}

	for _, enableMode := range checkConfiguration.EnableModes {
		if checkModes[enableMode] == true {
			return true
		}
	}
	return false
}

// TODO: make checkOutput a struct to allow for more advanced message templating
func message(templateText string, checkOutput string) string {
	messageTemplate := template.Must(template.New("messageTemplate").Parse(templateText))

	messageBuffer := new(bytes.Buffer)
	messageTemplate.Execute(messageBuffer, checkOutput)

	return messageBuffer.String()
}

func RunChecks(project projects.Project) {
	fmt.Printf("Checking %s in %s\n", project.Type.String(), project.Path.String())

	checkdata.Initialize(project)

	for _, checkConfiguration := range checkconfigurations.Configurations {
		if !shouldRun(checkConfiguration, project) {
			// TODO: this should only be printed to log and in verbose mode
			fmt.Printf("Skipping check: %s\n", checkConfiguration.ID)
			continue
		}

		fmt.Printf("Running check %s: ", checkConfiguration.ID)
		result, output := checkConfiguration.CheckFunction()
		fmt.Printf("%s\n", result.String())
		if result != checkresult.Pass {
			fmt.Printf("%s: %s\n", checklevel.CheckLevel(checkConfiguration).String(), message(checkConfiguration.MessageTemplate, output))
		}
	}
}
