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
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project"
)

func shouldRun(checkConfiguration checkconfigurations.Type, currentProject project.Type) bool {
	configurationCheckModes := configuration.CheckModes(currentProject.SuperprojectType)

	if checkConfiguration.ProjectType != currentProject.ProjectType {
		return false
	}

	for _, disableMode := range checkConfiguration.DisableModes {
		if configurationCheckModes[disableMode] == true {
			return false
		}
	}

	for _, enableMode := range checkConfiguration.EnableModes {
		if configurationCheckModes[enableMode] == true {
			return true
		}
	}

	// Use default
	for _, disableMode := range checkConfiguration.DisableModes {
		if disableMode == checkmode.Default {
			return false
		}
	}

	for _, enableMode := range checkConfiguration.EnableModes {
		if enableMode == checkmode.Default {
			return true
		}
	}

	// TODO: this should return an error
	return false
}

// TODO: make checkOutput a struct to allow for more advanced message templating
func message(templateText string, checkOutput string) string {
	messageTemplate := template.Must(template.New("messageTemplate").Parse(templateText))

	messageBuffer := new(bytes.Buffer)
	messageTemplate.Execute(messageBuffer, checkOutput)

	return messageBuffer.String()
}

func RunChecks(project project.Type) {
	fmt.Printf("Checking %s in %s\n", project.ProjectType.String(), project.Path.String())

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
		if (result != checkresult.Pass) && (result != checkresult.NotRun) {
			fmt.Printf("%s: %s\n", checklevel.CheckLevel(checkConfiguration).String(), message(checkConfiguration.MessageTemplate, output))
		}
	}
}
