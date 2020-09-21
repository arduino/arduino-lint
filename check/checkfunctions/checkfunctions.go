package checkfunctions

import (
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
)

// output is the contextual information that will be added to the stock message
type Type func() (result checkresult.Type, output string)

func LibraryPropertiesFormat() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.Fail, checkdata.LibraryPropertiesLoadError().Error()
	}
	return checkresult.Pass, ""
}

func LibraryPropertiesNameFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldMissing("name", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

func LibraryPropertiesNameFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldPatternMismatch("name", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

func LibraryPropertiesVersionFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldMissing("version", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}
