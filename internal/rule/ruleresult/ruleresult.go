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

// Package ruleresult defines the possible result values returned by a rule.
package ruleresult

// Type is the type for rule results.
//
//go:generate stringer -type=Type -linecomment
type Type int

const (
	// Pass indicates rule compliance.
	Pass Type = iota // pass
	// Fail indicates a rule violation.
	Fail // fail
	// Skip indicates the rule is configured to be skipped in the current tool configuration mode.
	Skip // skipped
	// NotRun indicates an unrelated error prevented the rule from running.
	NotRun // unable to run
)
