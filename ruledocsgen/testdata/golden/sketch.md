Arduino Lint provides 1 rules for the [`sketch`](https://arduino.github.io/arduino-cli/latest/sketch-specification/) project type:

---

<a id="SS001"></a>

## name mismatch (`SS001`)

There is no `.ino` sketch file with name matching the sketch folder. The primary sketch file name must match the folder for the sketch to be valid.

More information: [**here**](https://arduino.github.io/arduino-cli/latest/sketch-specification/#primary-sketch-file)<br />
Enabled for superproject type: all<br />
Category: structure<br />
Subcategory: root folder

##### Rule levels

| `compliance`  | Level   |
|:--------------|:--------|
| permissive    | WARNING |
| specification | ERROR   |
| strict        | ERROR   |
