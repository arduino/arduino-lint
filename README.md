# arduino-check

`arduino-check` is a command line tool that automatically checks for common problems in your
[Arduino](https://www.arduino.cc/) projects:

- Sketches
- Libraries

## Usage

After installing `arduino-check`, run the command `arduino-check --help` for usage documentation.

A few additional configuration options only of use for internal/development use of the tool can be set via environment
variables:

- `ARDUINO_CHECK_OFFICIAL` - Set to `"true"` to run the checks that only apply to official Arduino projects.
- `ARDUINO_CHECK_LOG_LEVEL` - Messages with this level and above will be logged.
  - Supported values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`
- `ARDUINO_CHECK_LOG_FORMAT` - The output format for the logs.
  - Supported values: `text`, `json`
