**Arduino Lint** is a command line tool that checks for common problems with [Arduino](https://www.arduino.cc/)
projects.

Its focus is on the structure, metadata, and configuration of Arduino projects, rather than the code. [Rules](rules.md)
cover specification compliance, Library Manager submission requirements, and best practices.

## Installation

See the [installation instructions](installation.md).

## Getting started

Once installed, you only need to open a terminal at your project folder and run the command:

```
arduino-lint
```

This will automatically detect the project type and check it against the relevant rules.

The default configuration of **Arduino Lint** provides for the most common use case, but you have the option of changing
settings via [command line flags](commands/arduino-lint.md):

### Compliance setting

The `--compliance` flag allows you to configure the strictness of the applied rules. The three compliance level values
accepted by this flag are:

- `permissive` - failure will occur only when severe rule violations are found. Although a project that passes at the
  permissive setting will work with the current Arduino development software versions, it may not be fully
  specification-compliant, risking incompatibility or a poor experience for the users.
- `specification` - the default setting, enforces compliance with the official Arduino project specifications
  ([sketch](https://arduino.github.io/arduino-cli/latest/sketch-specification/),
  [library](https://arduino.github.io/arduino-cli/latest/library-specification/),
  [platform](https://arduino.github.io/arduino-cli/latest/platform-specification/)).
- `strict` - enforces best practices, above and beyond the minimum requirements for specification compliance. Use this
  setting to ensure the best experience for the users of the project.

### Library Manager setting

[Arduino Library Manager](https://docs.arduino.cc/software/ide-v1/tutorials/installing-libraries#using-the-library-manager)
is the best way to provide installation and updates of Arduino libraries. In order to be accepted for inclusion in
Library Manager, a library is required to meet
[some requirements](https://github.com/arduino/library-registry/blob/main/FAQ.md#readme).

**Arduino Lint** provides checks for these requirements as well, controlled by the `--library-manager` flag.

The Library Manager submission-specific rules are enabled via `--library-manager submit`. Even if your library isn't yet
ready to be added to Library Manager, it's a good idea to use this setting to ensure no incompatibilities are
introduced.

Once your library is in the Library Manager index, each release is automatically picked up and made available to the
Arduino community. Releases are also subject to special rules. The command `arduino-lint --library-manager update` will
tell you whether your library is compliant with these rules.

### Integration

The `--format` flag configures the format of `arduino-lint`'s output. The default `--format text` setting provides human
readable output. For automation or integration with other tools, the machine readable output provided by `--format json`
may be more convenient. This setting exposes every detail of the rules that were applied.

The `--report-file` flag causes `arduino-lint` to write the JSON output to the specified file.

### Environment variables

Additional configuration options intended for internal use or development can be set via environment variables:

- `ARDUINO_LINT_OFFICIAL` - Set to `"true"` to run the checks that only apply to official Arduino projects.
- `ARDUINO_LINT_LIBRARY_MANAGER_INDEXING` - Set to `"true"` to run the checks that apply when adding releases to the
  Library Manager index.
- `ARDUINO_LINT_LOG_LEVEL` - Messages with this level and above will be logged.
  - Supported values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`
- `ARDUINO_LINT_LOG_FORMAT` - The output format for the logs.
  - Supported values: `text`, `json`

## Continuous integration

**Arduino Lint** would be a great addition to your
[continuous integration](https://wikipedia.org/wiki/Continuous_integration) system. Running the tool after each change
to the project can allow you to identify any problems that were introduced.

This is easily done by using the `arduino/arduino-lint-action` [GitHub Actions](https://docs.github.com/actions) action:
https://github.com/arduino/arduino-lint-action

Add [a simple workflow file](https://github.com/arduino/arduino-lint-action#usage) to the repository of your Arduino
project and GitHub will automatically run Arduino Lint on every pull request and push.

## Support and feedback

You can discuss or get assistance with using **Arduino Lint** on the
[Arduino Forum](https://forum.arduino.cc/c/using-arduino/project-guidance/19).

Feedback is welcome! Please submit feature requests or bug reports to the
[issue tracker](CONTRIBUTING.md#issue-reports).
