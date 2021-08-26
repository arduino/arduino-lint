Arduino Lint provides 2 rules for the [`library`](https://arduino.github.io/arduino-cli/latest/library-specification/) project type:

---

<a id="LS001"></a>

## invalid library (`LS001`)

The path does not contain a valid Arduino library.

More information: [**here**](https://arduino.github.io/arduino-cli/latest/library-specification)<br />
Enabled for superproject type: all<br />
Category: structure<br />
Subcategory: general

##### Rule levels

| `compliance`  | Level |
|---------------|-------|
| permissive    | ERROR |
| specification | ERROR |
| strict        | ERROR |

---

<a id="LS007"></a>

## .exe file (`LS007`)

A file with `.exe` file extension was found under the library folder. Presence of this file blocks addition to the Library Manager index.


Enabled for superproject type: library<br />
Category: structure<br />
Subcategory: miscellaneous

##### Rule levels

| `compliance`  | `library-manager` |  Level   |
|---------------|-------------------|----------|
| permissive    | submit            | ERROR    |
| permissive    | update            | ERROR    |
| permissive    | false             | disabled |
| specification | submit            | ERROR    |
| specification | update            | ERROR    |
| specification | false             | disabled |
| strict        | submit            | ERROR    |
| strict        | update            | ERROR    |
| strict        | false             | disabled |
