// Source: https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/assets/cobra/docsgen/go.mod
module github.com/arduino/arduino-lint/docsgen

go 1.16

replace github.com/arduino/arduino-lint => ../

replace github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1 // https://github.com/oleiade/reflections/issues/14

require (
	github.com/arduino/arduino-lint v0.0.0
	github.com/spf13/cobra v1.1.1
)
