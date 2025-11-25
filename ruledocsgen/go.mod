module github.com/arduino/arduino-lint/ruledocsgen

go 1.25.2

replace github.com/arduino/arduino-lint => ../

replace github.com/jandelgado/gcov2lcov => github.com/jandelgado/gcov2lcov v1.0.5 // v1.0.4 causes Dependabot updates to fail due to checksum mismatch (likely a moved tag). This is an unused transitive dependency, so version is irrelevant.

require (
	github.com/JohannesKaufmann/html-to-markdown v1.6.0
	github.com/arduino/arduino-lint v0.0.0
	github.com/arduino/go-paths-helper v1.14.0
	github.com/olekukonko/tablewriter v1.1.1
	github.com/stretchr/testify v1.11.1
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.1.6 // indirect
	github.com/arduino/arduino-cli v0.35.4-0.20241001142927-1f8d0f6c0dd3 // indirect
	github.com/arduino/go-properties-orderedmap v1.8.1 // indirect
	github.com/arduino/go-win32-utils v1.0.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/clipperhouse/displaywidth v0.3.1 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.2.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/cmaglie/pb v1.0.27 // indirect
	github.com/codeclysm/extract/v4 v4.0.0 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-git/go-git/v5 v5.16.4 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gobuffalo/pop/v5 v5.3.3 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jandelgado/gcov2lcov v1.0.4 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/juju/errors v1.0.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/leonelquinteros/gotext v1.4.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.1.2 // indirect
	github.com/ory/go-acc v0.2.6 // indirect
	github.com/ory/jsonschema/v3 v3.0.4 // indirect
	github.com/ory/viper v1.7.5 // indirect
	github.com/ory/x v0.0.272 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/seatgeek/logrus-gelf-formatter v0.0.0-20210414080842-5b05eb8ff761 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.17.0 // indirect
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.bug.st/cleanup v1.0.0 // indirect
	go.bug.st/downloader/v2 v2.1.1 // indirect
	go.bug.st/relaxed-semver v0.11.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.44.0 // indirect
	go.opentelemetry.io/otel v1.18.0 // indirect
	go.opentelemetry.io/otel/metric v1.18.0 // indirect
	go.opentelemetry.io/otel/trace v1.18.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
