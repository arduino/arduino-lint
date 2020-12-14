# arduino-lint

`arduino-lint` is a command line tool that automatically checks for common problems in your
[Arduino](https://www.arduino.cc/) projects:

- Sketches
- Libraries

## Usage

After installing `arduino-lint`, run the command `arduino-lint --help` for usage documentation.

A few additional configuration options only of use for internal/development use of the tool can be set via environment
variables:

- `ARDUINO_LINT_OFFICIAL` - Set to `"true"` to run the checks that only apply to official Arduino projects.
- `ARDUINO_LINT_LOG_LEVEL` - Messages with this level and above will be logged.
  - Supported values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`
- `ARDUINO_LINT_LOG_FORMAT` - The output format for the logs.
  - Supported values: `text`, `json`
