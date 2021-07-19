// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package rule runs rules on a project.
package rule

import (
	"fmt"

	"github.com/arduino/arduino-lint/internal/configuration"
	"github.com/arduino/arduino-lint/internal/configuration/rulemode"
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/result"
	"github.com/arduino/arduino-lint/internal/result/feedback"
	"github.com/arduino/arduino-lint/internal/rule/ruleconfiguration"
	"github.com/sirupsen/logrus"
)

// Runner runs all rules for the given project and outputs the results.
func Runner(project project.Type) {
	feedback.Printf("Linting %s in %s\n", project.ProjectType, project.Path)

	projectdata.Initialize(project)

	for _, ruleConfiguration := range ruleconfiguration.Configurations() {
		runRule, err := shouldRun(ruleConfiguration, project)
		if err != nil {
			panic(err)
		}

		if !runRule {
			logrus.Infof("Skipping rule: %s\n", ruleConfiguration.ID)
			continue
		}

		// Output will be printed after all rules are finished when configured for "json" output format.
		feedback.VerbosePrintf("Running rule %s (%s)...\n", ruleConfiguration.ID, ruleConfiguration.Brief)

		ruleResult, ruleOutput := ruleConfiguration.RuleFunction()
		reportText := result.Results.Record(project, ruleConfiguration, ruleResult, ruleOutput)
		feedback.Print(reportText)
	}
}

// shouldRun returns whether a given rule should be run for the given project under the current tool configuration.
func shouldRun(ruleConfiguration ruleconfiguration.Type, currentProject project.Type) (bool, error) {
	configurationRuleModes := configuration.RuleModes(currentProject.SuperprojectType)

	if !(ruleConfiguration.ProjectType.Matches(currentProject.ProjectType) && ruleConfiguration.SuperprojectType.Matches(currentProject.SuperprojectType)) {
		return false, nil
	}

	return IsEnabled(ruleConfiguration, configurationRuleModes)
}

// IsEnabled returns whether a given rule is enabled under a given tool configuration.
func IsEnabled(ruleConfiguration ruleconfiguration.Type, configurationRuleModes map[rulemode.Type]bool) (bool, error) {
	for _, disableMode := range ruleConfiguration.DisableModes {
		if configurationRuleModes[disableMode] {
			return false, nil
		}
	}

	for _, enableMode := range ruleConfiguration.EnableModes {
		if configurationRuleModes[enableMode] {
			return true, nil
		}
	}

	// Use default
	for _, disableMode := range ruleConfiguration.DisableModes {
		if disableMode == rulemode.Default {
			return false, nil
		}
	}

	for _, enableMode := range ruleConfiguration.EnableModes {
		if enableMode == rulemode.Default {
			return true, nil
		}
	}

	return false, fmt.Errorf("Rule %s is incorrectly configured", ruleConfiguration.ID)
}
