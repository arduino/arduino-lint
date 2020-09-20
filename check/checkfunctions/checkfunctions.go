package checkfunctions

import "github.com/arduino/arduino-check/check/checkresult"

// output is the contextual information that will be added to the stock message
type Type func() (result checkresult.Type, output string)

func CheckLibraryPropertiesFormat() (result checkresult.Type, output string) {
	return
}
