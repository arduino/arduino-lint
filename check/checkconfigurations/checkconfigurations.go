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
		Subcategory:     "general",
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
		CheckFunction:   checkfunctions.LibraryPropertiesFormat,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP002",
		Name:            "missing name field",
		Description:     "",
		MessageTemplate: "missing name field in library.properties",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldMissing,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "name field",
		ID:              "LP003",
		Name:            "disallowed characters",
		Description:     "",
		MessageTemplate: "disallowed characters in library.properties name field. See: https://arduino.github.io/arduino-cli/latest/library-specification/#libraryproperties-file-format",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.LibraryPropertiesNameFieldDisallowedCharacters,
	},
	{
		ProjectType:     projecttype.Library,
		Category:        "library.properties",
		Subcategory:     "version field",
		ID:              "LP004",
		Name:            "missing version field",
		Description:     "",
		MessageTemplate: "missing version field in library.properties",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    nil,
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.LibraryPropertiesVersionFieldMissing,
	},
	{
		ProjectType:     projecttype.Sketch,
		Category:        "structure",
		Subcategory:     "",
		ID:              "SS001",
		Name:            ".pde extension",
		Description:     "The .pde extension is used by both Processing sketches and Arduino sketches. Processing sketches should either be in the \"data\" subfolder of the sketch or in the \"extras\" folder of the library. Arduino sketches should use the modern .ino extension",
		MessageTemplate: "{{.}} uses deprecated .pde file extension. Use .ino for Arduino sketches",
		DisableModes:    nil,
		EnableModes:     []checkmode.Type{checkmode.Default},
		PromoteModes:    nil,
		InfoModes:       nil,
		WarningModes:    []checkmode.Type{checkmode.Permissive},
		ErrorModes:      []checkmode.Type{checkmode.Default},
		CheckFunction:   checkfunctions.PdeSketchExtension,
	},
}
