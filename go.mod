module github.com/arduino/arduino-lint

go 1.16

replace github.com/jandelgado/gcov2lcov => github.com/jandelgado/gcov2lcov v1.0.5 // v1.0.4 causes Dependabot updates to fail due to checksum mismatch (likely a moved tag). This is an unused transitive dependency, so version is irrelevant.

replace github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1 // https://github.com/oleiade/reflections/issues/14

require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/arduino/arduino-cli v0.0.0-20201210103408-bf7a3194bb63
	github.com/arduino/go-paths-helper v1.6.1
	github.com/arduino/go-properties-orderedmap v1.7.0
	github.com/client9/misspell v0.3.4
	github.com/daaku/go.zipexe v1.0.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/gliderlabs/ssh v0.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2
	github.com/h2non/filetype v1.1.0 // indirect
	github.com/juju/testing v0.0.0-20201030020617-7189b3728523 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/ory/jsonschema/v3 v3.0.6
	github.com/sirupsen/logrus v1.8.1
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	go.bug.st/relaxed-semver v0.0.0-20190922224835-391e10178d18
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
)
