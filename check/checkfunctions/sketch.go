package checkfunctions

// The check functions for sketches.

import (
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
)

func PdeSketchExtension() (result checkresult.Type, output string) {
	directoryListing, _ := checkdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()
	pdeSketches := []string{}
	for _, filePath := range directoryListing {
		if filePath.Ext() == ".pde" {
			pdeSketches = append(pdeSketches, filePath.Base())
		}
	}

	if len(pdeSketches) > 0 {
		return checkresult.Fail, strings.Join(pdeSketches, ", ")
	}

	return checkresult.Pass, ""
}
