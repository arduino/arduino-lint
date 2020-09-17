package checkfunctions

import "github.com/arduino/arduino-check/check/checkresult"

// output is the contextual information that will be added to the stock message
type CheckFunction func() (result checkresult.Result, output string)

func CheckLibraryPropertiesFormat() (result checkresult.Result, output string) {
	return
}
