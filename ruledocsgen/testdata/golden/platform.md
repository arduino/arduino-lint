Arduino Lint provides 2 rules for the [`platform`](https://arduino.github.io/arduino-cli/latest/platform-specification/) project type:

---

<a id="PF001"></a>

## boards.txt missing (`PF001`)

The `boards.txt` configuration file was not found in the platform folder

More information: [**here**](https://arduino.github.io/arduino-cli/latest/platform-specification/#boardstxt)<br />
Enabled for superproject type: all<br />
Category: configuration files<br />
Subcategory: boards.txt

##### Rule levels

| `compliance`  | Level |
|:--------------|:------|
| permissive    | ERROR |
| specification | ERROR |
| strict        | ERROR |

---

<a id="PF009"></a>

## use of compiler.&lt;pattern type&gt;.extra\_flags &amp; foo &#39;bar&#39; &#34;baz&#34; (`PF009`)

A board definition in the platform's `boards.txt` configuration file is using one of the `compiler.<pattern type>.extra_flags` properties (e.g., `compiler.cpp.extra_flags`). These are intended to be left for use by the user as a standardized interface for customizing the compilation commands. The platform author can define as many arbitrary properties as they like, so there is no need for them to take the user's properties.


Enabled for superproject type: all<br />
Category: configuration files<br />
Subcategory: boards.txt

##### Rule levels

| `compliance`  | Level   |
|:--------------|:--------|
| permissive    | WARNING |
| specification | WARNING |
| strict        | ERROR   |
