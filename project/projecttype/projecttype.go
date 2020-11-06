// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
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

// Type is the type for Arduino project types.
//go:generate stringer -type=Type -linecomment
type Type int

const (
	Sketch       Type = iota // sketch
	Library                  // library
	Platform                 // boards platform
	PackageIndex             // Boards Manager package index
	All                      // any project type
	Not                      // N/A
)

// Matches returns whether the receiver project type matches the argument project type
func (projectTypeA Type) Matches(projectTypeB Type) bool {
	if projectTypeA == Not && projectTypeB == Not {
		return true
	} else if projectTypeA == Not || projectTypeB == Not {
		return false
	}
	return (projectTypeA == All || projectTypeB == All || projectTypeA == projectTypeB)
}
