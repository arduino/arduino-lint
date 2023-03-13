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

// Package rulefunction contains the functions that implement each rule.
package rulefunction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/arduino/arduino-lint/internal/project/projectdata"
	"github.com/arduino/arduino-lint/internal/project/sketch"
	"github.com/arduino/arduino-lint/internal/rule/ruleresult"
	"github.com/arduino/go-paths-helper"
	"github.com/sirupsen/logrus"
)

// Type is the function signature for the rule functions.
// The `output` result is the contextual information that will be inserted into the rule's message template.
type Type func() (result ruleresult.Type, output string)

// MissingReadme checks if the project has a readme that will be recognized by GitHub.
func MissingReadme() (result ruleresult.Type, output string) {
	// https://github.com/github/markup/blob/master/README.md
	readmeRegexp := regexp.MustCompile(`(?i)^readme\.((markdown)|(mdown)|(mkdn)|(md)|(textile)|(rdoc)|(org)|(creole)|(mediawiki)|(wiki)|(rst)|(asciidoc)|(adoc)|(asc)|(pod)|(txt))$`)

	// https://docs.github.com/en/free-pro-team@latest/github/creating-cloning-and-archiving-repositories/about-readmes#about-readmes
	if pathContainsRegexpMatch(projectdata.ProjectPath(), readmeRegexp) ||
		(projectdata.ProjectPath().Join("docs").Exist() && pathContainsRegexpMatch(projectdata.ProjectPath().Join("docs"), readmeRegexp)) ||
		(projectdata.ProjectPath().Join(".github").Exist() && pathContainsRegexpMatch(projectdata.ProjectPath().Join(".github"), readmeRegexp)) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, ""
}

// MissingLicenseFile checks if the project has a license file that will be recognized by GitHub.
func MissingLicenseFile() (result ruleresult.Type, output string) {
	// https://docs.github.com/en/free-pro-team@latest/github/creating-cloning-and-archiving-repositories/licensing-a-repository#detecting-a-license
	// https://github.com/licensee/licensee/blob/master/docs/what-we-look-at.md#detecting-the-license-file
	// Should be `(?i)^(((un)?licen[sc]e)|(copy(ing|right))|(ofl)|(patents))(\.(?!spdx|header|gemspec).+)?$` but regexp package doesn't support negative lookahead, so only using "preferred extensions".
	// github.com/dlclark/regexp2 does support negative lookahead, but I'd prefer to stick with the standard package.
	licenseRegexp := regexp.MustCompile(`(?i)^(((un)?licen[sc]e)|(copy(ing|right))|(ofl)|(patents))(\.((md)|(markdown)|(txt)|(html)))?$`)

	// License file must be in root of repo
	if pathContainsRegexpMatch(projectdata.ProjectPath(), licenseRegexp) {
		return ruleresult.Pass, ""
	}

	return ruleresult.Fail, ""
}

// IncorrectArduinoDotHFileNameCase checks for incorrect file name case of Arduino.h in #include directives.
func IncorrectArduinoDotHFileNameCase() (result ruleresult.Type, output string) {
	incorrectCaseRegexp := regexp.MustCompile(`^\s*#\s*include\s*["<](a((?i)rduino)|(ARDUINO))\.[hH][">]`)

	directoryListing, err := projectdata.ProjectPath().ReadDirRecursive()
	if err != nil {
		panic(err)
	}
	directoryListing.FilterOutDirs()

	for _, file := range directoryListing {
		if !sketch.HasSupportedExtension(file) { // Won't catch all possible files, but good enough.
			continue
		}

		lines, err := file.ReadFileAsLines()
		if err != nil {
			panic(err)
		}

		for lineNumber, line := range lines {
			if incorrectCaseRegexp.MatchString(line) {
				return ruleresult.Fail, fmt.Sprintf("%s:%v: %s", file, lineNumber+1, line)
			}
		}
	}

	return ruleresult.Pass, ""
}

const brokenOutputListIndent = "  " // Use this as indent for rule output that takes the form of newline-separated list.

// brokenOutputList formats the rule output as a newline-separated list.
func brokenOutputList(list []string) string {
	return brokenOutputListIndent + strings.Join(list, "\n"+brokenOutputListIndent)
}

// validProjectPathBaseName checks whether the provided library folder or sketch filename contains prohibited characters.
func validProjectPathBaseName(name string) bool {
	baseNameRegexp := regexp.MustCompile("^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$")
	return baseNameRegexp.MatchString(name)
}

func containsMisspelledPathBaseName(pathList paths.PathList, correctBaseName string, misspellingQuery string) (*paths.Path, bool) {
	misspellingRegexp := regexp.MustCompile(misspellingQuery)
	for _, path := range pathList {
		if path.Base() == correctBaseName {
			return nil, false
		}

		if misspellingRegexp.MatchString(path.Base()) {
			return path, true
		}
	}

	return nil, false
}

// containsIncorrectPathBaseCase checks whether the list of paths contains an element with base name matching the provided query in all bug case.
func containsIncorrectPathBaseCase(pathList paths.PathList, correctBaseName string) (*paths.Path, bool) {
	for _, path := range pathList {
		if path.Base() == correctBaseName {
			// There was a case-sensitive match (paths package's Exist() is not always case-sensitive, so can't be used here).
			return nil, false
		}

		if strings.EqualFold(path.Base(), correctBaseName) {
			// There was a case-insensitive match.
			return path, true
		}
	}

	return nil, false
}

// pathContainsRegexpMatch checks if the provided path contains a file name matching the given regular expression.
func pathContainsRegexpMatch(path *paths.Path, pathRegexp *regexp.Regexp) bool {
	listing, err := path.ReadDir()
	if err != nil {
		panic(err)
	}
	listing.FilterOutDirs()

	for _, file := range listing {
		if pathRegexp.MatchString(file.Base()) {
			return true
		}
	}

	return false
}

// isValidJSON checks whether the specified file is a valid JSON document.
func isValidJSON(path *paths.Path) bool {
	data, err := path.ReadFile()
	if err != nil {
		panic(err)
	}
	return json.Valid(data)
}

func checkURL(url string) error {
	logrus.Tracef("Checking URL: %s", url)
	response, err := http.Head(url)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", response.Status)
	}

	return nil
}
