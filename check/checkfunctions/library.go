package checkfunctions

// The check functions for libraries.

import (
	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/project/library/libraryproperties"
)

// LibraryPropertiesFormat checks for invalid library.properties format.
func LibraryPropertiesFormat() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.Fail, checkdata.LibraryPropertiesLoadError().Error()
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldMissing checks for missing library.properties "name" field.
func LibraryPropertiesNameFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldMissing("name", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}

// LibraryPropertiesNameFieldDisallowedCharacters checks for disallowed characters in the library.properties "name" field.
func LibraryPropertiesNameFieldDisallowedCharacters() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldPatternMismatch("name", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}

	return checkresult.Pass, ""
}

// LibraryPropertiesVersionFieldMissing checks for missing library.properties "version" field.
func LibraryPropertiesVersionFieldMissing() (result checkresult.Type, output string) {
	if checkdata.LibraryPropertiesLoadError() != nil {
		return checkresult.NotRun, ""
	}

	if libraryproperties.FieldMissing("version", checkdata.LibraryPropertiesSchemaValidationResult()) {
		return checkresult.Fail, ""
	}
	return checkresult.Pass, ""
}
