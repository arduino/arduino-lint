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

// Package rulelevel defines the level assigned to a rule violation.
package rulelevel

import (
	"fmt"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
)

// Type is the type for the rule levels.
//
//go:generate stringer -type=Type -linecomment
type Type int

// The line comments set the string for each level.
const (
	Info    Type = iota // INFO
	Warning             // WARNING
	Error               // ERROR
	Notice              // NOTICE
)

// RuleLevel determines the rule level assigned to the given result of the given rule under the current tool configuration.
func RuleLevel(ruleConfiguration ruleconfiguration.Type, ruleResult ruleresult.Type, lintedProject project.Type) (Type, error) {
	if ruleResult != ruleresult.Fail {
		return Notice, nil // Level provided by FailRuleLevel() is only relevant for failure result.
	}
	configurationRuleModes := configuration.RuleModes(lintedProject.SuperprojectType)
	return FailRuleLevel(ruleConfiguration, configurationRuleModes)
}

// FailRuleLevel determines the level of a failed rule for the given rule modes.
func FailRuleLevel(ruleConfiguration ruleconfiguration.Type, configurationRuleModes map[rulemode.Type]bool) (Type, error) {
	for _, errorMode := range ruleConfiguration.ErrorModes {
		if configurationRuleModes[errorMode] {
			return Error, nil
		}
	}

	for _, warningMode := range ruleConfiguration.WarningModes {
		if configurationRuleModes[warningMode] {
			return Warning, nil
		}
	}

	for _, infoMode := range ruleConfiguration.InfoModes {
		if configurationRuleModes[infoMode] {
			return Info, nil
		}
	}

	// Use default level
	for _, errorMode := range ruleConfiguration.ErrorModes {
		if errorMode == rulemode.Default {
			return Error, nil
		}
	}

	for _, warningMode := range ruleConfiguration.WarningModes {
		if warningMode == rulemode.Default {
			return Warning, nil
		}
	}

	for _, infoMode := range ruleConfiguration.InfoModes {
		if infoMode == rulemode.Default {
			return Info, nil
		}
	}

	return Notice, fmt.Errorf("Rule %s is incorrectly configured", ruleConfiguration.ID)
}
