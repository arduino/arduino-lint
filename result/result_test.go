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

package result

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/util/test"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteReport(t *testing.T) {
	flags := test.ConfigurationFlags()

	projectPath, err := os.Getwd() // Path to an arbitrary folder that is guaranteed to exist.
	require.Nil(t, err)
	projectPaths := []string{projectPath}

	reportFolderPathString, err := ioutil.TempDir("", "arduino-check-result-TestWriteReport")
	require.Nil(t, err)
	defer os.RemoveAll(reportFolderPathString) // clean up
	reportFolderPath := paths.New(reportFolderPathString)

	reportFilePath := reportFolderPath.Join("report-file.json")
	_, err = reportFilePath.Create() // Create file using the report folder path.
	require.Nil(t, err)

	flags.Set("report-file", reportFilePath.Join("report-file.json").String())
	assert.Nil(t, configuration.Initialize(flags, projectPaths))
	assert.Error(t, Results.WriteReport(), "Parent folder creation should fail due to a collision with an existing file at that path")

	reportFilePath = reportFolderPath.Join("report-file-subfolder", "report-file-subsubfolder", "report-file.json")
	flags.Set("report-file", reportFilePath.String())
	assert.Nil(t, configuration.Initialize(flags, projectPaths))
	assert.NoError(t, Results.WriteReport(), "Creation of multiple levels of parent folders")

	reportFile, err := reportFilePath.Open()
	require.Nil(t, err)
	reportFileInfo, err := reportFile.Stat()
	require.Nil(t, err)
	reportFileBytes := make([]byte, reportFileInfo.Size())
	_, err = reportFile.Read(reportFileBytes)
	require.Nil(t, err)
	assert.True(t, assert.ObjectsAreEqualValues(reportFileBytes, Results.jsonReportRaw()), "Report file contents are correct")
}
