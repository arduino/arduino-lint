package schema

import (
	"runtime"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/require"
)

func TestPathURI(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		require.Equal(t, "file:///c:/foo%20bar", pathURI(paths.New("c:/foo bar")))
	default:
		require.Equal(t, "file:///foo%20bar", pathURI(paths.New("/foo bar")))
	}
}
