// Package feedback provides feedback to the user.
package feedback

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Errorf behaves like fmt.Printf but also logs the error.
func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

// Error behaves like fmt.Print but also logs the error.
func Error(errorMessage string) {
	fmt.Printf(errorMessage)
	logrus.Error(fmt.Sprint(errorMessage))
}
