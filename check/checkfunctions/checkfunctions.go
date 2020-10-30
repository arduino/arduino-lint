// Package checkfunctions contains the functions that implement each check.
package checkfunctions

import (
	"github.com/arduino/arduino-check/check/checkresult"
)

// Type is the function signature for the check functions.
// The `output` result is the contextual information that will be inserted into the check's message template.
type Type func() (result checkresult.Type, output string)
