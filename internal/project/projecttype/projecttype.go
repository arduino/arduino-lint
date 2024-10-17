// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License, either
// version 3 of the License, or (at your option) any later version.
// This license covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package projecttype defines the Arduino project types.
package projecttype

import (
	"fmt"
	"strings"
)

// Type is the type for Arduino project types.
//
//go:generate stringer -type=Type -linecomment
type Type int

const (
	// Sketch is used for Arduino sketch projects.
	Sketch Type = iota // sketch
	// Library is used for Arduino library projects.
	Library // library
	// Platform is used for Arduino boards platform projects.
	Platform // platform
	// PackageIndex is used for Arduino package index projects.
	PackageIndex // package-index
	// All is the catch-all for all supported Arduino project types.
	All // all
	// Not is the project type used when an Arduino project was not detected.
	Not // N/A
)

// FromString parses the --project-type flag value and returns the corresponding project type.
func FromString(projectTypeString string) (Type, error) {
	projectType, found := map[string]Type{
		Sketch.String():       Sketch,
		Library.String():      Library,
		Platform.String():     Platform,
		PackageIndex.String(): PackageIndex,
		All.String():          All,
	}[strings.ToLower(projectTypeString)]

	if found {
		return projectType, nil
	}
	return Not, fmt.Errorf("No matching project type for string %s", projectTypeString)
}

// Matches returns whether the receiver project type matches the argument project type.
func (projectTypeA Type) Matches(projectTypeB Type) bool {
	if projectTypeA == Not && projectTypeB == Not {
		return true
	} else if projectTypeA == Not || projectTypeB == Not {
		return false
	}
	return (projectTypeA == All || projectTypeB == All || projectTypeA == projectTypeB)
}
