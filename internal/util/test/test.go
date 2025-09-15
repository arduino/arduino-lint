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

// Package test provides resources for testing arduino-lint.
package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/spf13/pflag"
)

// ConfigurationFlags returns a set of the flags used for command line configuration of arduino-lint.
func ConfigurationFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	flags.String("compliance", "specification", "")
	flags.String("format", "text", "")
	flags.String("library-manager", "", "")
	flags.String("log-format", "text", "")
	flags.String("log-level", "panic", "")
	flags.String("project-type", "all", "")
	flags.Bool("recursive", true, "")
	flags.String("report-file", "", "")
	flags.Bool("verbose", false, "")
	flags.Bool("version", false, "")

	return flags
}

// StatusServer returns an HTTP test server that will respond with the given status code.
func StatusServer(status int) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))

	return server
}
