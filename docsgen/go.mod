// Source: https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/assets/cobra/docsgen/go.mod
module github.com/arduino/arduino-lint/docsgen

go 1.16

replace github.com/arduino/arduino-lint => ../

replace github.com/jandelgado/gcov2lcov => github.com/jandelgado/gcov2lcov v1.0.5 // v1.0.4 causes Dependabot updates to fail due to checksum mismatch (likely a moved tag). This is an unused transitive dependency, so version is irrelevant.

replace github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1 // https://github.com/oleiade/reflections/issues/14

require (
	github.com/arduino/arduino-lint v0.0.0
	github.com/spf13/cobra v1.2.1
)
