package checkconfigurations

import (
	"github.com/arduino/arduino-check/check/checkfunctions"
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
)

type Type struct {
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
	DisableModes []checkmode.Type
	// Check is only enabled when tool is in one of these modes
	EnableModes []checkmode.Type
	// In these modes, failed check is treated as an error, otherwise it's handled as normal.
	PromoteModes  []checkmode.Type
	InfoModes     []checkmode.Type
	WarningModes  []checkmode.Type
	ErrorModes    []checkmode.Type
	CheckFunction checkfunctions.Type
}

// Checks is an array of structs that define the configuration of each check.
var Configurations = []Type{
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP001",
		Name:            "invalid format",
		Description:     "",
		MessageTemplate: "library.properties has an invalid format: {{.}}",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.CheckLibraryPropertiesFormat,
	},
}
