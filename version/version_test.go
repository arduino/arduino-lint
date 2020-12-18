package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBuildInjectedInfo tests the Info strings passed to the binary at build time
// in order to have this test green launch your testing using the provided task (see /Taskfile.yml) or use:
//     go test -run TestBuildInjectedInfo -v ./... -ldflags '
//       -X github.com/arduino/arduino-lint/version.version=0.0.0-test.preview
//       -X github.com/arduino/arduino-lint/version.commit=deadbeef'
func TestBuildInjectedInfo(t *testing.T) {
	goldenAppName := "arduino-lint"
	goldenInfo := Info{
		Application: goldenAppName,
		Version:     "0.0.0-test.preview",
		Commit:      "deadbeef",
	}
	info := NewInfo(goldenAppName)
	require.Equal(t, goldenInfo.Application, info.Application)
	require.Equal(t, goldenInfo.Version, info.Version)
	require.Equal(t, goldenInfo.Commit, info.Commit)
}
