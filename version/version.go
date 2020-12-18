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

package version

import (
	"fmt"
)

var (
	defaultVersion = "0.0.0-git"
	version        = ""
	commit         = ""
	status         = "alpha"
	date           = ""
)

// Info FIXMEDOC
type Info struct {
	Application string `json:"Application"`
	Version     string `json:"Version"`
	Commit      string `json:"Commit"`
	Status      string `json:"Status"`
	Date        string `json:"Date"`
}

// NewInfo FIXMEDOC
func NewInfo(application string) *Info {
	return &Info{
		Application: application,
		Version:     version,
		Commit:      commit,
		Status:      status,
		Date:        date,
	}
}

func (i *Info) String() string {
	return fmt.Sprintf("%s %s Version: %s Commit: %s Date: %s", i.Application, i.Status, i.Version, i.Commit, i.Date)
}

//nolint:gochecknoinits
func init() {
	if version == "" {
		version = defaultVersion
	}
}
