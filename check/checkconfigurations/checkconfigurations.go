package checkconfigurations

import (
	"github.com/arduino/arduino-check/check/checkfunctions"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/projects/projecttype"
)

type Configuration struct {
	ProjectType projecttype.Type
	// Arbitrary text for the log
	Category    string
	Subcategory string
	// Unique check identifier
	ID          string
	Name        string
	Description string
	// The warning/error message template displayed when the check fails
	// The check function output will be filled into the template
	MessageTemplate string
	// Check is disabled when tool is in any of these modes
	DisableModes []configuration.CheckMode
	// Check is only enabled when tool is in one of these modes
	EnableModes []configuration.CheckMode
	// In these modes, failed check is treated as an error, otherwise it's handled as normal.
	PromoteModes  []configuration.CheckMode
	InfoModes     []configuration.CheckMode
	WarningModes  []configuration.CheckMode
	ErrorModes    []configuration.CheckMode
	CheckFunction checkfunctions.CheckFunction
}

// Checks is an array of structs that define the configuration of each check.
var Configurations = []Configuration{
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP001",
		Name:            "invalid format",
		Description:     "",
		MessageTemplate: "library.properties has an invalid format: {{.}}",
		DisableModes:    nil,
		EnableModes:     []configuration.CheckMode{configuration.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []configuration.CheckMode{configuration.Default},
		CheckFunction:   checkfunctions.CheckLibraryPropertiesFormat,
	},
}
