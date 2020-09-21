package checkfunctions

import (
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
)

func PdeSketchExtension() (result checkresult.Type, output string) {
	directoryListing, _ := checkdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()
	pdeSketches := ""
	for _, filePath := range directoryListing {
		if filePath.Ext() == ".pde" {
			if pdeSketches == "" {
				pdeSketches = filePath.Base()
			} else {
				pdeSketches += ", " + filePath.Base()
			}
		}
	}

	if pdeSketches != "" {
		return checkresult.Fail, pdeSketches
	}

	return checkresult.Pass, ""
}
