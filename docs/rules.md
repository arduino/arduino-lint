**Arduino Lint** inspects Arduino projects for common problems. This is done by checking the project against a series of
"rules", each of which is targeted to detecting a specific potential issue. Only the rules of relevance to the project
being linted are applied.

## Rule documentation

Additional information is available for each of the **Arduino Lint** rules, organized by project type:

- [Sketch](rules/sketch.md)
- [Library](rules/library.md)
- [Boards platform](rules/platform.md)
- [Package index](rules/package-index.md)

## Rule ID

In order to allow particular rules to be referenced unequivocally, each has been assigned a permanent unique
identification code (e.g., `SS001`).

## Rule level

In addition to checking for critical flaws, **Arduino Lint** also advocates for best practices in Arduino projects. For
this reason, not all rule violations are treated as fatal linting errors.

A violation of a rule is assigned a level according to its severity. In cases where a rule violation indicates a serious
problem with the project, the violation is treated as an error, and will result in a non-zero exit status from the
`arduino-lint` command. In cases where the violation indicates a possible problem, or where the rule is a recommendation
for an optional improvement to enhance the project user's experience, the violation is treated as a warning. It is hoped
that these warning-level violations will be given consideration by the user, but they do not affect the `arduino-lint`
exit status.

Of the hundreds of rules provided by **Arduino Lint**, only the ones relevant to the current target project are applied,
with the rest disabled.

The rule levels and enabled subset of rules is dependent on the target project type and how the user has configured
Arduino Lint via the [command line flags](commands/arduino-lint.md) and
[environment variables](index.md#environment-variables).

## Projects and "superprojects"

Arduino projects may contain other Arduino projects as subprojects. For example, the libraries
[bundled](https://arduino.github.io/arduino-cli/latest/platform-specification/#platform-bundled-libraries) with an
Arduino boards platform. These subprojects may, in turn, contain their own subprojects, such as the example sketches
included with a platform bundled library.

**Arduino Lint** also lints any subprojects the target project might contain. The type of the top level "superproject"
is a factor in the configuration of some rules. For example, there is no need to take Library Manager requirements into
consideration in the case of a platform bundled library, since it will never be distributed via that system.
