// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package checkfunctions contains the functions that implement each check.
package checkfunctions

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/go-paths-helper"
)

// Type is the function signature for the check functions.
// The `output` result is the contextual information that will be inserted into the check's message template.
type Type func() (result checkresult.Type, output string)

// validProjectPathBaseName checks whether the provided library folder or sketch filename contains prohibited characters.
func validProjectPathBaseName(name string) bool {
	baseNameRegexp := regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")
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

// MissingReadme checks if the project has a readme that will be recognized by GitHub.
func MissingReadme() (result checkresult.Type, output string) {
	// https://github.com/github/markup/blob/master/README.md
	readmeRegexp := regexp.MustCompile(`(?i)^readme\.(markdown)|(mdown)|(mkdn)|(md)|(textile)|(rdoc)|(org)|(creole)|(mediawiki)|(wiki)|(rst)|(asciidoc)|(adoc)|(asc)|(pod)|(txt)$`)

	// https://docs.github.com/en/free-pro-team@latest/github/creating-cloning-and-archiving-repositories/about-readmes#about-readmes
	if pathContainsReadme(checkdata.ProjectPath(), readmeRegexp) ||
		(checkdata.ProjectPath().Join("docs").Exist() && pathContainsReadme(checkdata.ProjectPath().Join("docs"), readmeRegexp)) ||
		(checkdata.ProjectPath().Join(".github").Exist() && pathContainsReadme(checkdata.ProjectPath().Join(".github"), readmeRegexp)) {
		return checkresult.Pass, ""
	}

	return checkresult.Fail, ""
}

// pathContainsReadme checks if the provided path contains a readme file recognized by GitHub.
func pathContainsReadme(path *paths.Path, readmeRegexp *regexp.Regexp) bool {
	listing, err := path.ReadDir()
	if err != nil {
		panic(err)
	}
	listing.FilterOutDirs()

	for _, file := range listing {
		if readmeRegexp.MatchString(file.Base()) {
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
