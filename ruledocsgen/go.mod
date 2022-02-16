module github.com/arduino/arduino-lint/ruledocsgen

go 1.16

replace github.com/arduino/arduino-lint => ../

replace github.com/jandelgado/gcov2lcov => github.com/jandelgado/gcov2lcov v1.0.5 // v1.0.4 causes Dependabot updates to fail due to checksum mismatch (likely a moved tag). This is an unused transitive dependency, so version is irrelevant.

require (
	github.com/JohannesKaufmann/html-to-markdown v1.3.3
	github.com/arduino/arduino-lint v0.0.0
	github.com/arduino/go-paths-helper v1.7.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stretchr/testify v1.7.0
)
