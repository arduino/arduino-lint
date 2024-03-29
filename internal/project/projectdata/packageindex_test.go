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

package projectdata

import (
	"testing"

	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
)

var packageIndexTestDataPath *paths.Path

func init() {
	workingDirectory, err := paths.Getwd()
	if err != nil {
		panic(err)
	}
	packageIndexTestDataPath = workingDirectory.Join("testdata", "packageindexes")
}

func TestInitializeForPackageIndex(t *testing.T) {
	testTables := []struct {
		testName                                       string
		path                                           *paths.Path
		packageIndexAssertion                          assert.ValueAssertionFunc
		packageIndexLoadErrorAssertion                 assert.ValueAssertionFunc
		packageIndexCLILoadErrorAssertion              assert.ValueAssertionFunc
		packageIndexPackagesAssertion                  assert.ValueAssertionFunc
		packageIndexPackagesDataAssertion              []PackageIndexData
		packageIndexPlatformsAssertion                 assert.ValueAssertionFunc
		packageIndexPlatformsDataAssertion             []PackageIndexData
		packageIndexBoardsAssertion                    assert.ValueAssertionFunc
		packageIndexBoardsDataAssertion                []PackageIndexData
		packageIndexToolsDependenciesAssertion         assert.ValueAssertionFunc
		packageIndexToolsDependenciesDataAssertion     []PackageIndexData
		packageIndexDiscoveryDependenciesAssertion     assert.ValueAssertionFunc
		packageIndexDiscoveryDependenciesDataAssertion []PackageIndexData
		packageIndexMonitorDependenciesAssertion       assert.ValueAssertionFunc
		packageIndexMonitorDependenciesDataAssertion   []PackageIndexData
		packageIndexToolsAssertion                     assert.ValueAssertionFunc
		packageIndexToolsDataAssertion                 []PackageIndexData
		packageIndexSystemsAssertion                   assert.ValueAssertionFunc
		packageIndexSystemsDataAssertion               []PackageIndexData
		packageIndexSchemaValidationResultAssertion    assert.ValueAssertionFunc
	}{
		{
			testName:                          "Valid",
			path:                              packageIndexTestDataPath.Join("valid-package-index", "package_foo_index.json"),
			packageIndexAssertion:             assert.NotNil,
			packageIndexLoadErrorAssertion:    assert.Nil,
			packageIndexCLILoadErrorAssertion: assert.Nil,
			packageIndexPackagesAssertion:     assert.NotNil,
			packageIndexPackagesDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1",
					JSONPointer: "/packages/0",
				},
				{
					ID:          "foopackager2",
					JSONPointer: "/packages/1",
				},
			},
			packageIndexPlatformsAssertion: assert.NotNil,
			packageIndexPlatformsDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1:avr@1.0.0",
					JSONPointer: "/packages/0/platforms/0",
				},
				{
					ID:          "foopackager1:avr@1.0.1",
					JSONPointer: "/packages/0/platforms/1",
				},
				{
					ID:          "foopackager2:samd@2.0.0",
					JSONPointer: "/packages/1/platforms/0",
				},
				{
					ID:          "foopackager2:mbed@1.1.1",
					JSONPointer: "/packages/1/platforms/1",
				},
			},
			packageIndexBoardsAssertion: assert.NotNil,
			packageIndexBoardsDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1:avr@1.0.0 >> My Board",
					JSONPointer: "/packages/0/platforms/0/boards/0",
				},
				{
					ID:          "foopackager1:avr@1.0.0 >> My Board Pro",
					JSONPointer: "/packages/0/platforms/0/boards/1",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> My Board",
					JSONPointer: "/packages/0/platforms/1/boards/0",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> My Board Pro",
					JSONPointer: "/packages/0/platforms/1/boards/1",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> My Board",
					JSONPointer: "/packages/1/platforms/0/boards/0",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> My Board Pro",
					JSONPointer: "/packages/1/platforms/0/boards/1",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> My Board",
					JSONPointer: "/packages/1/platforms/1/boards/0",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> My Board Pro",
					JSONPointer: "/packages/1/platforms/1/boards/1",
				},
			},
			packageIndexToolsDependenciesAssertion: assert.NotNil,
			packageIndexToolsDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1:avr@1.0.0 >> arduino:avr-gcc@4.8.1-arduino5",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "foopackager1:avr@1.0.0 >> arduino:avrdude@6.0.1-arduino5",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> arduino:avr-gcc@7.3.0-atmel3.6.1-arduino7",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> arduino:avrdude@6.3.0-arduino17",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/1",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> arduino:arm-none-eabi-gcc@7-2017q4",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> arduino:bossac@1.7.0-arduino3",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> arduino:openocd@0.11.0-arduino2",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> arduino:arm-none-eabi-gcc@7-2017q4",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/1",
				},
			},
			packageIndexDiscoveryDependenciesAssertion: assert.NotNil,
			packageIndexDiscoveryDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1:avr@1.0.1 >> arduino:ble-discovery",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/0",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> barpackager:carrier-pigeon-discovery",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/1",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> arduino:ble-discovery",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/0",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> bazpackager:signal-flag-discovery",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/1",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> quxpackager:sneakernet-discovery",
					JSONPointer: "/packages/1/platforms/1/discoveryDependencies/0",
				},
			},
			packageIndexMonitorDependenciesAssertion: assert.NotNil,
			packageIndexMonitorDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager1:avr@1.0.1 >> arduino:network-monitor",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/0",
				},
				{
					ID:          "foopackager1:avr@1.0.1 >> barpackager:coop-monitor",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/1",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> arduino:network-monitor",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/0",
				},
				{
					ID:          "foopackager2:samd@2.0.0 >> bazpackager:flag-monitor",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/1",
				},
				{
					ID:          "foopackager2:mbed@1.1.1 >> quxpackager:shoe-monitor",
					JSONPointer: "/packages/1/platforms/1/monitorDependencies/0",
				},
			},
			packageIndexToolsAssertion: assert.NotNil,
			packageIndexToolsDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager2:openocd@0.10.0-arduino1-static",
					JSONPointer: "/packages/1/tools/0",
				},
				{
					ID:          "foopackager2:CMSIS@4.0.0-atmel",
					JSONPointer: "/packages/1/tools/1",
				},
			},
			packageIndexSystemsAssertion: assert.NotNil,
			packageIndexSystemsDataAssertion: []PackageIndexData{
				{
					ID:          "foopackager2:openocd@0.10.0-arduino1-static >> i386-apple-darwin11",
					JSONPointer: "/packages/1/tools/0/systems/0",
				},
				{
					ID:          "foopackager2:openocd@0.10.0-arduino1-static >> x86_64-linux-gnu",
					JSONPointer: "/packages/1/tools/0/systems/1",
				},
				{
					ID:          "foopackager2:CMSIS@4.0.0-atmel >> arm-linux-gnueabihf",
					JSONPointer: "/packages/1/tools/1/systems/0",
				},
				{
					ID:          "foopackager2:CMSIS@4.0.0-atmel >> i686-mingw32",
					JSONPointer: "/packages/1/tools/1/systems/1",
				},
			},
			packageIndexSchemaValidationResultAssertion: assert.NotNil,
		},
		{
			testName:                          "Missing IDs",
			path:                              packageIndexTestDataPath.Join("missing-ids", "package_foo_index.json"),
			packageIndexAssertion:             assert.NotNil,
			packageIndexLoadErrorAssertion:    assert.Nil,
			packageIndexCLILoadErrorAssertion: assert.Nil,
			packageIndexPackagesAssertion:     assert.NotNil,
			packageIndexPackagesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0",
					JSONPointer: "/packages/0",
				},
				{
					ID:          "foopackager2",
					JSONPointer: "/packages/1",
				},
			},
			packageIndexPlatformsAssertion: assert.NotNil,
			packageIndexPlatformsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0",
					JSONPointer: "/packages/0/platforms/0",
				},
				{
					ID:          "/packages/0/platforms/1",
					JSONPointer: "/packages/0/platforms/1",
				},
				{
					ID:          "/packages/1/platforms/0",
					JSONPointer: "/packages/1/platforms/0",
				},
				{
					ID:          "/packages/1/platforms/1",
					JSONPointer: "/packages/1/platforms/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0",
					JSONPointer: "/packages/1/platforms/2",
				},
			},
			packageIndexBoardsAssertion: assert.NotNil,
			packageIndexBoardsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/boards/0",
					JSONPointer: "/packages/0/platforms/0/boards/0",
				},
				{
					ID:          "/packages/0/platforms/0/boards/1",
					JSONPointer: "/packages/0/platforms/0/boards/1",
				},
				{
					ID:          "/packages/0/platforms/1/boards/0",
					JSONPointer: "/packages/0/platforms/1/boards/0",
				},
				{
					ID:          "/packages/0/platforms/1/boards/1",
					JSONPointer: "/packages/0/platforms/1/boards/1",
				},
				{
					ID:          "/packages/1/platforms/0/boards/0",
					JSONPointer: "/packages/1/platforms/0/boards/0",
				},
				{
					ID:          "/packages/1/platforms/0/boards/1",
					JSONPointer: "/packages/1/platforms/0/boards/1",
				},
				{
					ID:          "/packages/1/platforms/1/boards/0",
					JSONPointer: "/packages/1/platforms/1/boards/0",
				},
				{
					ID:          "/packages/1/platforms/1/boards/1",
					JSONPointer: "/packages/1/platforms/1/boards/1",
				},
				{
					ID:          "/packages/1/platforms/2/boards/0",
					JSONPointer: "/packages/1/platforms/2/boards/0",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> My Board Pro",
					JSONPointer: "/packages/1/platforms/2/boards/1",
				},
			},
			packageIndexToolsDependenciesAssertion: assert.NotNil,
			packageIndexToolsDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/toolsDependencies/0",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/toolsDependencies/1",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/toolsDependencies/0",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/toolsDependencies/1",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/2",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/2",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> arduino:CMSIS@4.5.0",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/3",
				},
			},
			packageIndexDiscoveryDependenciesAssertion: assert.NotNil,
			packageIndexDiscoveryDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/discoveryDependencies/0",
					JSONPointer: "/packages/0/platforms/0/discoveryDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/discoveryDependencies/1",
					JSONPointer: "/packages/0/platforms/0/discoveryDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/discoveryDependencies/0",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/discoveryDependencies/1",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/1/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/1/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> quxpackager:sneakernet-discovery",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/2",
				},
			},
			packageIndexMonitorDependenciesAssertion: assert.NotNil,
			packageIndexMonitorDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/monitorDependencies/0",
					JSONPointer: "/packages/0/platforms/0/monitorDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/monitorDependencies/1",
					JSONPointer: "/packages/0/platforms/0/monitorDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/monitorDependencies/0",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/monitorDependencies/1",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/1/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/1/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> quxpackager:shoe-monitor",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/2",
				},
			},
			packageIndexToolsAssertion: assert.NotNil,
			packageIndexToolsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/1/tools/0",
					JSONPointer: "/packages/1/tools/0",
				},
				{
					ID:          "/packages/1/tools/1",
					JSONPointer: "/packages/1/tools/1",
				},
				{
					ID:          "foopackager2:bossac@1.9.1-arduino2",
					JSONPointer: "/packages/1/tools/2",
				},
			},
			packageIndexSystemsAssertion: assert.NotNil,
			packageIndexSystemsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/1/tools/0/systems/0",
					JSONPointer: "/packages/1/tools/0/systems/0",
				},
				{
					ID:          "/packages/1/tools/0/systems/1",
					JSONPointer: "/packages/1/tools/0/systems/1",
				},
				{
					ID:          "/packages/1/tools/1/systems/0",
					JSONPointer: "/packages/1/tools/1/systems/0",
				},
				{
					ID:          "/packages/1/tools/1/systems/1",
					JSONPointer: "/packages/1/tools/1/systems/1",
				},
				{
					ID:          "/packages/1/tools/2/systems/0",
					JSONPointer: "/packages/1/tools/2/systems/0",
				},
			},
			packageIndexSchemaValidationResultAssertion: assert.NotNil,
		},
		{
			testName:                          "Empty IDs",
			path:                              packageIndexTestDataPath.Join("empty-ids", "package_foo_index.json"),
			packageIndexAssertion:             assert.NotNil,
			packageIndexLoadErrorAssertion:    assert.Nil,
			packageIndexCLILoadErrorAssertion: assert.Nil,
			packageIndexPackagesAssertion:     assert.NotNil,
			packageIndexPackagesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0",
					JSONPointer: "/packages/0",
				},
				{
					ID:          "foopackager2",
					JSONPointer: "/packages/1",
				},
			},
			packageIndexPlatformsAssertion: assert.NotNil,
			packageIndexPlatformsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0",
					JSONPointer: "/packages/0/platforms/0",
				},
				{
					ID:          "/packages/0/platforms/1",
					JSONPointer: "/packages/0/platforms/1",
				},
				{
					ID:          "/packages/1/platforms/0",
					JSONPointer: "/packages/1/platforms/0",
				},
				{
					ID:          "/packages/1/platforms/1",
					JSONPointer: "/packages/1/platforms/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0",
					JSONPointer: "/packages/1/platforms/2",
				},
			},
			packageIndexBoardsAssertion: assert.NotNil,
			packageIndexBoardsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/boards/0",
					JSONPointer: "/packages/0/platforms/0/boards/0",
				},
				{
					ID:          "/packages/0/platforms/0/boards/1",
					JSONPointer: "/packages/0/platforms/0/boards/1",
				},
				{
					ID:          "/packages/0/platforms/1/boards/0",
					JSONPointer: "/packages/0/platforms/1/boards/0",
				},
				{
					ID:          "/packages/0/platforms/1/boards/1",
					JSONPointer: "/packages/0/platforms/1/boards/1",
				},
				{
					ID:          "/packages/1/platforms/0/boards/0",
					JSONPointer: "/packages/1/platforms/0/boards/0",
				},
				{
					ID:          "/packages/1/platforms/0/boards/1",
					JSONPointer: "/packages/1/platforms/0/boards/1",
				},
				{
					ID:          "/packages/1/platforms/1/boards/0",
					JSONPointer: "/packages/1/platforms/1/boards/0",
				},
				{
					ID:          "/packages/1/platforms/1/boards/1",
					JSONPointer: "/packages/1/platforms/1/boards/1",
				},
				{
					ID:          "/packages/1/platforms/2/boards/0",
					JSONPointer: "/packages/1/platforms/2/boards/0",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> My Board Pro",
					JSONPointer: "/packages/1/platforms/2/boards/1",
				},
			},
			packageIndexToolsDependenciesAssertion: assert.NotNil,
			packageIndexToolsDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/toolsDependencies/0",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/toolsDependencies/1",
					JSONPointer: "/packages/0/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/toolsDependencies/0",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/toolsDependencies/1",
					JSONPointer: "/packages/0/platforms/1/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/0/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/1/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/0",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/1",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/toolsDependencies/2",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/2",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> arduino:CMSIS@4.5.0",
					JSONPointer: "/packages/1/platforms/2/toolsDependencies/3",
				},
			},
			packageIndexDiscoveryDependenciesAssertion: assert.NotNil,
			packageIndexDiscoveryDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/discoveryDependencies/0",
					JSONPointer: "/packages/0/platforms/0/discoveryDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/discoveryDependencies/1",
					JSONPointer: "/packages/0/platforms/0/discoveryDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/discoveryDependencies/0",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/discoveryDependencies/1",
					JSONPointer: "/packages/0/platforms/1/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/0/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/1/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/1/discoveryDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/discoveryDependencies/0",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/discoveryDependencies/1",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> quxpackager:sneakernet-discovery",
					JSONPointer: "/packages/1/platforms/2/discoveryDependencies/2",
				},
			},
			packageIndexMonitorDependenciesAssertion: assert.NotNil,
			packageIndexMonitorDependenciesDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/0/platforms/0/monitorDependencies/0",
					JSONPointer: "/packages/0/platforms/0/monitorDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/0/monitorDependencies/1",
					JSONPointer: "/packages/0/platforms/0/monitorDependencies/1",
				},
				{
					ID:          "/packages/0/platforms/1/monitorDependencies/0",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/0",
				},
				{
					ID:          "/packages/0/platforms/1/monitorDependencies/1",
					JSONPointer: "/packages/0/platforms/1/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/0/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/0/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/0/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/1/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/1/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/1/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/1/monitorDependencies/1",
				},
				{
					ID:          "/packages/1/platforms/2/monitorDependencies/0",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/0",
				},
				{
					ID:          "/packages/1/platforms/2/monitorDependencies/1",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/1",
				},
				{
					ID:          "foopackager2:megaavr@1.0.0 >> quxpackager:shoe-monitor",
					JSONPointer: "/packages/1/platforms/2/monitorDependencies/2",
				},
			},
			packageIndexToolsAssertion: assert.NotNil,
			packageIndexToolsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/1/tools/0",
					JSONPointer: "/packages/1/tools/0",
				},
				{
					ID:          "/packages/1/tools/1",
					JSONPointer: "/packages/1/tools/1",
				},
				{
					ID:          "foopackager2:bossac@1.9.1-arduino2",
					JSONPointer: "/packages/1/tools/2",
				},
			},
			packageIndexSystemsAssertion: assert.NotNil,
			packageIndexSystemsDataAssertion: []PackageIndexData{
				{
					ID:          "/packages/1/tools/0/systems/0",
					JSONPointer: "/packages/1/tools/0/systems/0",
				},
				{
					ID:          "/packages/1/tools/0/systems/1",
					JSONPointer: "/packages/1/tools/0/systems/1",
				},
				{
					ID:          "/packages/1/tools/1/systems/0",
					JSONPointer: "/packages/1/tools/1/systems/0",
				},
				{
					ID:          "/packages/1/tools/1/systems/1",
					JSONPointer: "/packages/1/tools/1/systems/1",
				},
				{
					ID:          "/packages/1/tools/2/systems/0",
					JSONPointer: "/packages/1/tools/2/systems/0",
				},
			},
			packageIndexSchemaValidationResultAssertion: assert.NotNil,
		},
		{
			testName:                                    "Invalid package index",
			path:                                        packageIndexTestDataPath.Join("invalid-package-index", "package_foo_index.json"),
			packageIndexAssertion:                       assert.Nil,
			packageIndexLoadErrorAssertion:              assert.NotNil,
			packageIndexCLILoadErrorAssertion:           assert.NotNil,
			packageIndexPackagesAssertion:               assert.Nil,
			packageIndexPlatformsAssertion:              assert.Nil,
			packageIndexBoardsAssertion:                 assert.Nil,
			packageIndexToolsDependenciesAssertion:      assert.Nil,
			packageIndexDiscoveryDependenciesAssertion:  assert.Nil,
			packageIndexMonitorDependenciesAssertion:    assert.Nil,
			packageIndexToolsAssertion:                  assert.Nil,
			packageIndexSystemsAssertion:                assert.Nil,
			packageIndexSchemaValidationResultAssertion: assert.Nil,
		},
		{
			testName:                                    "Invalid JSON",
			path:                                        packageIndexTestDataPath.Join("invalid-JSON", "package_foo_index.json"),
			packageIndexAssertion:                       assert.Nil,
			packageIndexLoadErrorAssertion:              assert.NotNil,
			packageIndexCLILoadErrorAssertion:           assert.NotNil,
			packageIndexPackagesAssertion:               assert.Nil,
			packageIndexPlatformsAssertion:              assert.Nil,
			packageIndexBoardsAssertion:                 assert.Nil,
			packageIndexToolsDependenciesAssertion:      assert.Nil,
			packageIndexDiscoveryDependenciesAssertion:  assert.Nil,
			packageIndexMonitorDependenciesAssertion:    assert.Nil,
			packageIndexToolsAssertion:                  assert.Nil,
			packageIndexSystemsAssertion:                assert.Nil,
			packageIndexSchemaValidationResultAssertion: assert.Nil,
		},
	}

	for _, testTable := range testTables {

		testProject := project.Type{
			Path:             testTable.path,
			ProjectType:      projecttype.PackageIndex,
			SuperprojectType: projecttype.PackageIndex,
		}
		Initialize(testProject)

		testTable.packageIndexLoadErrorAssertion(t, PackageIndexLoadError(), testTable.testName)
		testTable.packageIndexCLILoadErrorAssertion(t, PackageIndexCLILoadError(), testTable.testName)
		if PackageIndexLoadError() == nil {
			testTable.packageIndexAssertion(t, PackageIndex(), testTable.testName)
		}

		testTable.packageIndexPackagesAssertion(t, PackageIndexPackages(), testTable.testName)
		if PackageIndexPackages() != nil {
			for index, packageIndexPackage := range PackageIndexPackages() {
				assert.Equal(t, testTable.packageIndexPackagesDataAssertion[index].ID, packageIndexPackage.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexPackagesDataAssertion[index].JSONPointer, packageIndexPackage.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexPlatformsAssertion(t, PackageIndexPlatforms(), testTable.testName)
		if PackageIndexPlatforms() != nil {
			for index, packageIndexPlatform := range PackageIndexPlatforms() {
				assert.Equal(t, testTable.packageIndexPlatformsDataAssertion[index].ID, packageIndexPlatform.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexPlatformsDataAssertion[index].JSONPointer, packageIndexPlatform.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexBoardsAssertion(t, PackageIndexBoards(), testTable.testName)
		if PackageIndexBoards() != nil {
			for index, packageIndexBoard := range PackageIndexBoards() {
				assert.Equal(t, testTable.packageIndexBoardsDataAssertion[index].ID, packageIndexBoard.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexBoardsDataAssertion[index].JSONPointer, packageIndexBoard.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexToolsDependenciesAssertion(t, PackageIndexToolsDependencies(), testTable.testName)
		if PackageIndexToolsDependencies() != nil {
			for index, packageIndexToolsDependency := range PackageIndexToolsDependencies() {
				assert.Equal(t, testTable.packageIndexToolsDependenciesDataAssertion[index].ID, packageIndexToolsDependency.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexToolsDependenciesDataAssertion[index].JSONPointer, packageIndexToolsDependency.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexDiscoveryDependenciesAssertion(t, PackageIndexDiscoveryDependencies(), testTable.testName)
		if PackageIndexDiscoveryDependencies() != nil {
			for index, packageIndexDiscoveryDependency := range PackageIndexDiscoveryDependencies() {
				assert.Equal(t, testTable.packageIndexDiscoveryDependenciesDataAssertion[index].ID, packageIndexDiscoveryDependency.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexDiscoveryDependenciesDataAssertion[index].JSONPointer, packageIndexDiscoveryDependency.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexMonitorDependenciesAssertion(t, PackageIndexMonitorDependencies(), testTable.testName)
		if PackageIndexMonitorDependencies() != nil {
			for index, packageIndexMonitorDependency := range PackageIndexMonitorDependencies() {
				assert.Equal(t, testTable.packageIndexMonitorDependenciesDataAssertion[index].ID, packageIndexMonitorDependency.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexMonitorDependenciesDataAssertion[index].JSONPointer, packageIndexMonitorDependency.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexToolsAssertion(t, PackageIndexTools(), testTable.testName)
		if PackageIndexTools() != nil {
			for index, packageIndexTool := range PackageIndexTools() {
				assert.Equal(t, testTable.packageIndexToolsDataAssertion[index].ID, packageIndexTool.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexToolsDataAssertion[index].JSONPointer, packageIndexTool.JSONPointer, testTable.testName)
			}
		}

		testTable.packageIndexSystemsAssertion(t, PackageIndexSystems(), testTable.testName)
		if PackageIndexSystems() != nil {
			for index, packageIndexSystem := range PackageIndexSystems() {
				assert.Equal(t, testTable.packageIndexSystemsDataAssertion[index].ID, packageIndexSystem.ID, testTable.testName)
				assert.Equal(t, testTable.packageIndexSystemsDataAssertion[index].JSONPointer, packageIndexSystem.JSONPointer, testTable.testName)
			}
		}
	}
}
