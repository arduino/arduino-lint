// This file is part of arduino-lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package outputformat defines the output formats
package outputformat

import (
	"fmt"
	"strings"
)

// Type is the type for output formats
//go:generate stringer -type=Type -linecomment
type Type int

const (
	Text Type = iota // text
	JSON             // json
)

// FromString parses the --format flag value and returns the corresponding output format type.
func FromString(outputFormatString string) (Type, error) {
	formatType, found := map[string]Type{
		Text.String(): Text,
		JSON.String(): JSON,
	}[strings.ToLower(outputFormatString)]

	if found {
		return formatType, nil
	}
	return Text, fmt.Errorf("No matching output format for string %s", outputFormatString)
}
