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

// Package compliance level defines the levels of specification compliance.
package compliancelevel

// Type is the type for the compliance levels.
//go:generate stringer -type=Type -linecomment
type Type int

// The line comments set the string for each level.
const (
	Permissive    Type = iota // permissive
	Specification             // standard
	Strict                    // strict
)
