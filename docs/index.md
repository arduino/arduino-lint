**arduino-lint** is a command line tool that checks for common problems with [Arduino](https://www.arduino.cc/)
projects.

Its focus is on the structure, metadata, and configuration of Arduino projects, rather than the code. Rules cover
[specification](https://arduino.github.io/arduino-cli/latest/library-specification) compliance, Library Manager
submission [requirements](https://github.com/arduino/Arduino/wiki/Library-Manager-FAQ), and best practices.

## Installation

See the [installation instructions](installation.md).

## Getting started

Once installed, you only need to open a terminal at your project folder and run the command:

```
arduino-lint
```

This will automatically search for projects, detect their type, and run the appropriate checks on them.

The default configuration of **arduino-lint** provides for the most common use case, but you have the option of changing
settings via [command line flags](commands/arduino-lint.md):

### Compliance setting

The `--compliance` flag allows you to configure the strictness of the checks.

`--compliance permissive` will cause the checks to fail only when severe problems are found. Although a project that
passes at the permissive setting will work with the current Arduino development software versions, it may not be fully
specification-compliant, risking incompatibility or a poor experience for the users.

`--compliance specification`, the default setting, enforces compliance with the official Arduino project specifications.

`--compliance strict` enforces best practices, above and beyond the minimum requirements for specification compliance.
Use this setting to ensure the best experience for the users of the project.

### Library Manager setting

Arduino Library Manager is the best way to provide installation and updates of Arduino libraries. In order to be
accepted for inclusion in Library Manager, a library is required to meet some requirements.

**arduino-lint** provides checks for these requirements as well, controlled by the `--library-manager` flag.

The checks for the Library Manager submission requirements are enabled via `--library-manager submit`. Even if your
library isn't yet ready to be added to Library Manager, it's a good idea to use this setting to ensure no
incompatibilities are introduced.

Once your library is in the Library Manager index, each release is automatically picked up and made available to the
Arduino community. Releases are also subjected to special checks, which must pass for it to be added to the index. The
command `arduino-lint --library-manager update` will tell you whether your library will pass these checks.

### Integration

The `--format` flag configures the format of **arduino-lint**'s output. The default `--format text` setting provides
human readable output. For automation or integration with other tools, the machine readable output provided by
`--format json` may be more convenient. This setting exposes all the details of the checks that are run.

The `--report-file` flag causes **arduino-lint** to write the JSON output to the specified file.

### Environment variables

Additional configuration options intended for internal use or development can be set via environment variables:

- `ARDUINO_LINT_OFFICIAL` - Set to `"true"` to run the checks that only apply to official Arduino projects.
- `ARDUINO_LINT_LOG_LEVEL` - Messages with this level and above will be logged.
  - Supported values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`
- `ARDUINO_LINT_LOG_FORMAT` - The output format for the logs.
  - Supported values: `text`, `json`

## Continuous integration

**arduino-lint** would be a great addition to your
[continuous integration](https://en.wikipedia.org/wiki/Continuous_integration) system. Running the checks after each
change to the project will allow you to identify any problems that were introduced.

This is easily done by using the `arduino/arduino-lint-action` GitHub Actions action:
https://github.com/arduino/arduino-lint-action

Add a simple workflow file to the repository of your Arduino project and the checks will automatically be run on every
pull request and push.

## Support and feedback

You can discuss or get assistance with using **arduino-lint** on the
[Arduino Forum](https://forum.arduino.cc/index.php?board=3.0).

Feedback is welcome! Please submit feature requests or bug reports to the
[issue tracker](CONTRIBUTING.md#issue-reports).
